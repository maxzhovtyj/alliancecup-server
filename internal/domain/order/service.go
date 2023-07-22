package order

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v9"
	"github.com/jung-kurt/gofpdf"
	goinvoice "github.com/maxzhovtyj/go-invoice"
	"github.com/sirupsen/logrus"
	"github.com/zh0vtyj/alliancecup-server/internal/domain/product"
	"github.com/zh0vtyj/alliancecup-server/internal/domain/shopping"
	"github.com/zh0vtyj/alliancecup-server/pkg/telegram"
	"os"
	"strings"
)

type Service interface {
	New(order CreateDTO, cartUUID string) (int, error)
	AdminNew(order CreateDTO) (int, error)
	GetUserOrders(userId int, createdAt string) ([]SelectDTO, error)
	GetOrderById(orderId int) (SelectDTO, error)
	AdminGetOrders(status, lastOrderCreatedAt, search string) ([]Order, error)
	DeliveryPaymentTypes() (shopping.DeliveryPaymentTypes, error)
	HandleOrderStatus(orderId int, status string) error
	GetInvoice(orderId int) (gofpdf.Fpdf, error)
}

type service struct {
	repo         Storage
	productRepo  product.Storage
	cache        *redis.Client
	tgBotManager telegram.Manager
}

func NewOrdersService(
	repo Storage,
	productRepo product.Storage,
	cache *redis.Client,
	tgBotManager telegram.Manager,
) Service {
	return &service{
		repo:         repo,
		productRepo:  productRepo,
		cache:        cache,
		tgBotManager: tgBotManager,
	}
}

func (s *service) New(order CreateDTO, cartUUID string) (int, error) {
	id, err := s.repo.New(order)
	if err != nil {
		return 0, err
	}

	delCartCache := s.cache.Del(context.Background(), cartUUID)
	if delCartCache.Err() != nil {
		return 0, fmt.Errorf("failed to delete cart from cache")
	}

	s.sentTelegramNewOrder(id)

	return id, nil
}

func (s *service) AdminNew(order CreateDTO) (int, error) {
	id, err := s.repo.New(order)
	if err != nil {
		return 0, err
	}

	return id, nil
}

func (s *service) GetUserOrders(userId int, createdAt string) ([]SelectDTO, error) {
	return s.repo.GetUserOrders(userId, createdAt)
}

func (s *service) GetOrderById(orderId int) (SelectDTO, error) {
	return s.repo.GetOrderById(orderId)
}

func (s *service) AdminGetOrders(status, lastOrderCreatedAt, search string) ([]Order, error) {
	return s.repo.AdminGetOrders(status, lastOrderCreatedAt, search)
}

func (s *service) DeliveryPaymentTypes() (shopping.DeliveryPaymentTypes, error) {
	deliveryTypes, err := s.repo.GetDeliveryTypes()
	if err != nil {
		return shopping.DeliveryPaymentTypes{}, fmt.Errorf("failed to load delivery types due to: %v", err)
	}

	paymentTypes, err := s.repo.GetPaymentTypes()
	if err != nil {
		return shopping.DeliveryPaymentTypes{}, fmt.Errorf("failed to load payment types due to: %v", err)
	}

	return shopping.DeliveryPaymentTypes{
		DeliveryTypes: deliveryTypes,
		PaymentTypes:  paymentTypes,
	}, nil
}

func (s *service) HandleOrderStatus(orderId int, status string) error {
	order, err := s.repo.GetOrderById(orderId)
	if err != nil {
		return err
	}

	if order.Info.Status == status {
		return fmt.Errorf("order status is already %s", status)
	}

	err = s.repo.ChangeOrderStatus(orderId, status)
	if err != nil {
		return err
	}

	if status == "PROCESSED" {
		err = s.repo.ProcessedOrder(orderId)
		if err != nil {
			return err
		}
	}

	return nil
}

func (s *service) GetInvoice(orderId int) (gofpdf.Fpdf, error) {
	order, err := s.repo.GetOrderById(orderId)
	if err != nil {
		return gofpdf.Fpdf{}, err
	}

	pwd, err := os.Getwd()
	if err != nil {
		return gofpdf.Fpdf{}, err
	}

	doc, err := goinvoice.NewWithCyrillic(pwd)
	if err != nil {
		return gofpdf.Fpdf{}, err
	}

	doc.Language = "UK"

	doc.SetPwd(pwd)
	doc.OrderId = fmt.Sprintf("%d", orderId)
	doc.SetInvoice(&goinvoice.Invoice{
		Logo:  nil,
		Title: "Alliance Cup",
	})
	doc.SetCompany(&goinvoice.Company{
		Name:        "AllianceCup",
		Address:     "Шухевича, 22",
		PhoneNumber: "+380(96) 512-15-16",
		City:        "Рівне",
		Country:     "Україна",
	})

	deliveryInfo := order.Info.Delivery.String()
	if deliveryInfo == "{}" {
		deliveryInfo = ""
	}

	doc.SetCustomer(&goinvoice.Customer{
		LastName:     order.Info.UserLastName,
		FirstName:    order.Info.UserFirstName,
		MiddleName:   order.Info.UserMiddleName,
		PhoneNumber:  order.Info.UserPhoneNumber,
		Email:        order.Info.UserEmail,
		DeliveryType: order.Info.DeliveryTypeTitle,
		DeliveryInfo: deliveryInfo,
	})

	for _, p := range order.Products {
		doc.AppendProductItem(&goinvoice.Product{
			Title:     p.ProductTitle,
			Price:     p.Price,
			Quantity:  float64(p.Quantity),
			Total:     p.PriceForQuantity,
			Packaging: "уп",
		})
	}

	pdf, err := doc.BuildPdf()
	if err != nil {
		return gofpdf.Fpdf{}, err
	}

	return pdf, err
}

func (s *service) sentTelegramNewOrder(oid int) {
	newOrder, err := s.GetOrderById(oid)
	if err != nil {
		errMsg := fmt.Sprintf(`Створено нове замовлення №%d. Не вдалося отримати інформацію`, oid)
		err = s.tgBotManager.Send(errMsg)
		if err != nil {
			logrus.Error("failed to send telegram message")
		}
	}

	msg := s.getNewOrderMessage(newOrder)

	err = s.tgBotManager.Send(msg)
	if err != nil {
		logrus.Error("failed to send telegram message")
	}
}

func (s *service) getNewOrderMessage(order SelectDTO) string {
	var orderProductsMsg string
	for i, p := range order.Products {
		orderProductsMsg += fmt.Sprintf(
			`

Товар %d, id #%d
Назва: %s
Арктикул: %s
Пакування: %s
Кількість на складі: %f
Ціна: %f,
Кількість: %d,
Сума: %f,

`,
			i,
			p.Id,
			p.ProductTitle,
			p.Article,
			p.Packaging.String(),
			p.AmountInStock,
			p.Price,
			p.Quantity,
			p.PriceForQuantity,
		)
	}

	newOrderMessage := fmt.Sprintf(`
Нове Замовлення №%d:
ПІБ: №%d %s
Номер телефону: %s
Email: %s
Доставка: %s %s
Оплата: %s
Сума: %f
Виконано адміністратором: %d,
Коментар: %s

Товари:
%s
`,
		order.Info.Id,
		order.Info.UserId,
		strings.Join([]string{order.Info.UserLastName, order.Info.UserFirstName, order.Info.UserMiddleName}, " "),
		order.Info.UserEmail,
		order.Info.UserPhoneNumber,
		order.Info.DeliveryTypeTitle,
		order.Info.Delivery.String(),
		order.Info.PaymentTypeTitle,
		order.Info.SumPrice,
		order.Info.ExecutedBy,
		*order.Info.Comment,
		orderProductsMsg,
	)

	return newOrderMessage
}
