package order

import (
	"fmt"
	"github.com/zh0vtyj/allincecup-server/internal/domain/product"
	server "github.com/zh0vtyj/allincecup-server/internal/domain/shopping"
)

type Service interface {
	New(order Info) (int, error)
	GetUserOrders(userId int, createdAt string) ([]FullInfo, error)
	GetOrderById(orderId int) (FullInfo, error)
	GetAdminOrders(status string, lastOrderCreatedAt string) ([]Order, error)
	DeliveryPaymentTypes() (server.DeliveryPaymentTypes, error)
	ProcessedOrder(orderId int) error
}

type service struct {
	repo        Storage
	productRepo product.Storage
}

func NewOrdersService(repo Storage, productRepo product.Storage) Service {
	return &service{
		repo:        repo,
		productRepo: productRepo,
	}
}

func (o *service) OrderSumCount(products []Product) (float64, error) {
	var sum float64
	for _, item := range products {
		p, err := o.productRepo.GetProductById(item.ProductId)
		if err != nil {
			return 0, err
		}
		inputPrice := item.PriceForQuantity / float64(item.Quantity)
		if p.Product.Price != inputPrice {
			return 0, fmt.Errorf("invalid product price")
		}
		if p.Product.Price*float64(item.Quantity) != item.PriceForQuantity {
			return 0, fmt.Errorf("price for quantity mismatch")
		}
		sum += p.Product.Price * float64(item.Quantity)
	}
	return sum, nil
}

func (o *service) New(order Info) (int, error) {
	// TODO
	sum, err := o.OrderSumCount(order.Products)
	if err != nil {
		return 0, err
	}
	if sum != order.Order.OrderSumPrice {
		return 0, fmt.Errorf("sum price mismatch, %f (computed) !== %f (given)", sum, order.Order.OrderSumPrice)
	}

	id, err := o.repo.New(order)
	if err != nil {
		return 0, err
	}

	return id, nil
}

func (o *service) GetUserOrders(userId int, createdAt string) ([]FullInfo, error) {
	return o.repo.GetUserOrders(userId, createdAt)
}

func (o *service) GetOrderById(orderId int) (FullInfo, error) {
	return o.repo.GetOrderById(orderId)
}

func (o *service) GetAdminOrders(status string, lastOrderCreatedAt string) ([]Order, error) {
	return o.repo.GetAdminOrders(status, lastOrderCreatedAt)
}

func (o *service) DeliveryPaymentTypes() (server.DeliveryPaymentTypes, error) {
	deliveryTypes, err := o.repo.GetDeliveryTypes()
	if err != nil {
		return server.DeliveryPaymentTypes{}, fmt.Errorf("failed to load delivery types due to: %v", err)
	}

	paymentTypes, err := o.repo.GetPaymentTypes()
	if err != nil {
		return server.DeliveryPaymentTypes{}, fmt.Errorf("failed to load payment types due to: %v", err)
	}

	return server.DeliveryPaymentTypes{
		DeliveryTypes: deliveryTypes,
		PaymentTypes:  paymentTypes,
	}, nil
}

func (o *service) ProcessedOrder(orderId int) error {
	err := o.repo.ProcessedOrder(orderId)
	if err != nil {
		return err
	}

	return nil
}
