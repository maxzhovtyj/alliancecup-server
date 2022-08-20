package service

import (
	"errors"
	server "github.com/zh0vtyj/allincecup-server"
	"github.com/zh0vtyj/allincecup-server/pkg/repository"
)

type ShoppingService struct {
	repo repository.Shopping
}

func NewShoppingService(repo repository.Shopping) *ShoppingService {
	return &ShoppingService{repo: repo}
}

func (s *ShoppingService) AddToCart(userId int, info server.CartProduct) (float64, error) {
	price, err := s.repo.PriceValidation(info.ProductId, info.Quantity)
	if err != nil {
		return 0, err
	}

	if price != info.PriceForQuantity {
		return 0, errors.New("price mismatch")
	}

	return s.repo.AddToCart(userId, info)
}

func (s *ShoppingService) GetProductsInCart(userId int) ([]server.CartProductFullInfo, float64, error) {
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

func (s *ShoppingService) DeleteFromCart(productId int) error {
	return s.repo.DeleteFromCart(productId)
}

func (s *ShoppingService) AddToFavourites(userId, productId int) error {
	return s.repo.AddToFavourites(userId, productId)
}

func (s *ShoppingService) GetFavourites(userId int) ([]server.Product, error) {
	return s.repo.GetFavourites(userId)
}

func (s *ShoppingService) DeleteFromFavourites(userId, productId int) error {
	return s.repo.DeleteFromFavourites(userId, productId)
}
