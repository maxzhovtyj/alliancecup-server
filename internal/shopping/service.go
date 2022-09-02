package shopping

import (
	"errors"
)

type Service interface {
	AddToCart(userId int, info CartProduct) (float64, error)
	GetProductsInCart(userId int) ([]CartProductFullInfo, float64, error)
	DeleteFromCart(productId int) error
	AddToFavourites(userId, productId int) error
	DeleteFromFavourites(userId, productId int) error
}

type service struct {
	repo Storage
}

func NewShoppingService(repo Storage) Service {
	return &service{repo: repo}
}

func (s *service) AddToCart(userId int, info CartProduct) (float64, error) {
	price, err := s.repo.PriceValidation(info.ProductId, info.Quantity)
	if err != nil {
		return 0, err
	}

	if price != info.PriceForQuantity {
		return 0, errors.New("price mismatch")
	}

	return s.repo.AddToCart(userId, info)
}

func (s *service) GetProductsInCart(userId int) ([]CartProductFullInfo, float64, error) {
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

func (s *service) DeleteFromCart(productId int) error {
	return s.repo.DeleteFromCart(productId)
}

func (s *service) AddToFavourites(userId, productId int) error {
	return s.repo.AddToFavourites(userId, productId)
}

func (s *service) DeleteFromFavourites(userId, productId int) error {
	return s.repo.DeleteFromFavourites(userId, productId)
}
