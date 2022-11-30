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
	GetUserOrders(userId int, createdAt string) ([]SelectDTO, error)
	GetOrderById(orderId int) (SelectDTO, error)
	AdminGetOrders(status, lastOrderCreatedAt, search string) ([]Order, error)
	DeliveryPaymentTypes() (shopping.DeliveryPaymentTypes, error)
	ProcessedOrder(orderId int) error
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

func (s *service) OrderSumCount(products []Product) (float64, error) {
	var sum float64
	for _, item := range products {
		p, err := s.productRepo.GetProductById(item.Id)
		if err != nil {
			return 0, err
		}
		inputPrice := item.PriceForQuantity / float64(item.Quantity)
		if p.Price != inputPrice {
			return 0, fmt.Errorf("invalid product price")
		}
		if p.Price*float64(item.Quantity) != item.PriceForQuantity {
			return 0, fmt.Errorf("price for quantity mismatch")
		}
		sum += p.Price * float64(item.Quantity)
	}
	return sum, nil
}

func (s *service) New(order CreateDTO, cartUUID string) (int, error) {
	sum, err := s.OrderSumCount(order.Products)
	if err != nil {
		return 0, err
	}
	if sum != order.Order.SumPrice {
		return 0, fmt.Errorf("sum price mismatch, %f (computed) !== %f (given)", sum, order.Order.SumPrice)
	}

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

func (s *service) ProcessedOrder(orderId int) error {
	err := s.repo.ProcessedOrder(orderId)
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
