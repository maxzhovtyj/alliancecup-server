package service

import (
	"fmt"
	"github.com/google/uuid"
	server "github.com/zh0vtyj/allincecup-server"
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
		return [16]byte{}, fmt.Errorf("sum price mismatch")
	}

	id, err := o.repo.New(order)
	if err != nil {
		return [16]byte{}, err
	}

	return id, nil
}

func (o *OrdersService) GetUserOrders(userId int, createdAt string) ([]server.Order, error) {
	return o.repo.GetUserOrders(userId, createdAt)
}
