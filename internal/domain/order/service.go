package order

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v9"
	"github.com/jung-kurt/gofpdf"
	goinvoice "github.com/maxzhovtyj/go-invoice"
	"github.com/zh0vtyj/allincecup-server/internal/domain/product"
	"github.com/zh0vtyj/allincecup-server/internal/domain/shopping"
	"os"
)

type Service interface {
	New(order CreateDTO, cartUUID string) (int, error)
	AdminNew(order CreateDTO) (int, error)
	GetUserOrders(userId int, createdAt string) ([]SelectDTO, error)
	GetOrderById(orderId int) (SelectDTO, error)
	AdminGetOrders(status, lastOrderCreatedAt, search string) ([]Order, error)
	DeliveryPaymentTypes() (shopping.DeliveryPaymentTypes, error)
	ProcessedOrder(orderId int, status string) error
	GetInvoice(orderId int) (gofpdf.Fpdf, error)
}

type service struct {
	repo        Storage
	productRepo product.Storage
	cache       *redis.Client
}

func NewOrdersService(repo Storage, productRepo product.Storage, cache *redis.Client) Service {
	return &service{
		repo:        repo,
		productRepo: productRepo,
		cache:       cache,
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

func (s *service) ProcessedOrder(orderId int, status string) error {
	err := s.repo.ProcessedOrder(orderId, status)
	if err != nil {
		return err
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

	var orderDeliveryInfo string
	for i, d := range order.Delivery {
		orderDeliveryInfo += d.DeliveryDescription
		if len(order.Delivery)-1 != i {
			orderDeliveryInfo += ", "
		}
	}

	doc.SetCustomer(&goinvoice.Customer{
		LastName:     order.Info.UserLastName,
		FirstName:    order.Info.UserFirstName,
		MiddleName:   order.Info.UserMiddleName,
		PhoneNumber:  order.Info.UserPhoneNumber,
		Email:        order.Info.UserEmail,
		DeliveryType: order.Info.DeliveryTypeTitle,
		DeliveryInfo: orderDeliveryInfo,
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
