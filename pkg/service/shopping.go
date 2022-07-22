package service

import (
	server "allincecup-server"
	"allincecup-server/pkg/repository"
	"errors"
)

type ShoppingService struct {
	repo repository.Shopping
}

func NewShoppingService(repo repository.Shopping) *ShoppingService {
	return &ShoppingService{repo: repo}
}

func (s *ShoppingService) AddToCart(userId int, info server.ProductOrder) (float64, error) {
	price, err := s.repo.PriceValidation(info.ProductId, info.Quantity)
	if err != nil {
		return 0, err
	}

	if price != info.PriceForQuantity {
		return 0, errors.New("price mismatch")
	}

	return s.repo.AddToCart(userId, info)
}

func (s *ShoppingService) GetProductsInCart(userId int) ([]server.ProductOrder, float64, error) {
	products, err := s.repo.GetProductsInCart(userId)
	if err != nil {
		return nil, 0, err
	}
	
	var sum float64
	for _, e := range products {
		sum += e.PriceForQuantity
	}

	return products, sum, err
}
