package order

import (
	"fmt"
	"github.com/google/uuid"
	server "github.com/zh0vtyj/allincecup-server/internal/shopping"
)

type Service interface {
	New(order Info) (uuid.UUID, error)
	GetUserOrders(userId int, createdAt string) ([]FullInfo, error)
	GetOrderById(orderId uuid.UUID) (FullInfo, error)
	GetAdminOrders(status string, lastOrderCreatedAt string) ([]Order, error)
	DeliveryPaymentTypes() (server.DeliveryPaymentTypes, error)
	ProcessedOrder(orderId uuid.UUID, toStatus string) error
}

type service struct {
	repo Storage
}

func NewOrdersService(repo Storage) Service {
	return &service{repo: repo}
}

//func (o *service) OrderSumCount(products []Product) (float64, error) {
//	var sum float64
//	for _, item := range products {
//		p, err := o.productsRepo.GetProductById(item.ProductId)
//		if err != nil {
//			return 0, err
//		}
//		inputPrice := item.PriceForQuantity / float64(item.Quantity)
//		if p.Info.Price != inputPrice {
//			return 0, fmt.Errorf("invalid product price")
//		}
//		if p.Info.Price*float64(item.Quantity) != item.PriceForQuantity {
//			return 0, fmt.Errorf("price for quantity mismatch")
//		}
//		sum += p.Info.Price * float64(item.Quantity)
//	}
//	return sum, nil
//}

func (o *service) New(order Info) (uuid.UUID, error) {
	// TODO
	//sum, err := o.OrderSumCount(order.Products)
	//if err != nil {
	//	return [16]byte{}, err
	//}
	//if sum != order.Order.OrderSumPrice {
	//	return [16]byte{}, fmt.Errorf("sum price mismatch, %f (computed) !== %f (given)", sum, order.Order.OrderSumPrice)
	//}

	id, err := o.repo.New(order)
	if err != nil {
		return [16]byte{}, err
	}

	return id, nil
}

func (o *service) GetUserOrders(userId int, createdAt string) ([]FullInfo, error) {
	return o.repo.GetUserOrders(userId, createdAt)
}

func (o *service) GetOrderById(orderId uuid.UUID) (FullInfo, error) {
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

func (o *service) ProcessedOrder(orderId uuid.UUID, toStatus string) error {
	err := o.repo.ChangeOrderStatus(orderId, toStatus)
	if err != nil {
		return err
	}
	return nil
}
