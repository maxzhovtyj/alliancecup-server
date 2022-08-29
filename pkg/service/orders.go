package service

import (
	"fmt"
	"github.com/google/uuid"
	server "github.com/zh0vtyj/allincecup-server/pkg/models"
	"github.com/zh0vtyj/allincecup-server/pkg/repository"
)

type OrdersService struct {
	repo         repository.Orders
	productsRepo repository.Products
}

func NewOrdersService(repo repository.Orders, productsRepo repository.Products) *OrdersService {
	return &OrdersService{repo: repo, productsRepo: productsRepo}
}

func (o *OrdersService) OrderSumCount(products []server.OrderProducts) (float64, error) {
	var sum float64
	for _, product := range products {
		p, err := o.productsRepo.GetProductById(product.ProductId)
		if err != nil {
			return 0, err
		}
		inputPrice := product.PriceForQuantity / float64(product.Quantity)
		if p.Info.Price != inputPrice {
			return 0, fmt.Errorf("invalid product price")
		}
		if p.Info.Price*float64(product.Quantity) != product.PriceForQuantity {
			return 0, fmt.Errorf("price for quantity mismatch")
		}
		sum += p.Info.Price * float64(product.Quantity)
	}
	return sum, nil
}

func (o *OrdersService) New(order server.OrderFullInfo) (uuid.UUID, error) {
	sum, err := o.OrderSumCount(order.Products)
	if err != nil {
		return [16]byte{}, err
	}
	if sum != order.Info.OrderSumPrice {
		return [16]byte{}, fmt.Errorf("sum price mismatch, %f (computed) !== %f (given)", sum, order.Info.OrderSumPrice)
	}

	id, err := o.repo.New(order)
	if err != nil {
		return [16]byte{}, err
	}

	return id, nil
}

func (o *OrdersService) GetUserOrders(userId int, createdAt string) ([]server.OrderInfo, error) {
	return o.repo.GetUserOrders(userId, createdAt)
}

func (o *OrdersService) GetOrderById(orderId uuid.UUID) (server.OrderInfo, error) {
	return o.repo.GetOrderById(orderId)
}

func (o *OrdersService) GetAdminOrders(status string, lastOrderCreatedAt string) ([]server.Order, error) {
	return o.repo.GetAdminOrders(status, lastOrderCreatedAt)
}

func (o *OrdersService) DeliveryPaymentTypes() (server.DeliveryPaymentTypes, error) {
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

func (o *OrdersService) ProcessedOrder(orderId uuid.UUID, toStatus string) error {
	err := o.repo.ChangeOrderStatus(orderId, toStatus)
	if err != nil {
		return err
	}
	return nil
}
