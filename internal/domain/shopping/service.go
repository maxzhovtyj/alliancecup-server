package shopping

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/go-redis/redis/v9"
	"github.com/google/uuid"
	"time"
)

type Service interface {
	NewCart(userId int) (uuid.UUID, error)
	AddToCart(userId int, info CartProduct) (float64, error)
	GetCart(cartId string, userId int) ([]CartProductFullInfo, float64, error)
	DeleteFromCart(productId int) error
	AddToFavourites(userId, productId int) error
	DeleteFromFavourites(userId, productId int) error
}

type service struct {
	repo  Storage
	cache *redis.Client
}

func NewShoppingService(repo Storage, cache *redis.Client) Service {
	return &service{
		repo:  repo,
		cache: cache,
	}
}

func (s *service) NewCart(userId int) (cartUUID uuid.UUID, err error) {
	cartUUID = uuid.New()

	var productsBytes []byte
	if userId != 0 {
		products, err := s.repo.GetProductsInCart(userId)
		if err != nil {
			return [16]byte{}, err
		}

		productsBytes, err = json.Marshal(products)
		if err != nil {
			return [16]byte{}, err
		}
	}

	set := s.cache.Set(context.Background(), cartUUID.String(), productsBytes, 72*time.Hour)
	if set.Err() != nil {
		return [16]byte{}, fmt.Errorf("failed to create card in cache")
	}

	return cartUUID, nil
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

func (s *service) GetCart(cartId string, userId int) (cart []CartProductFullInfo, sum float64, err error) {
	// todo
	var cartUUID uuid.UUID
	if cartId != "" {
		cartUUID, err = uuid.Parse(cartId)
		if err != nil {
			return nil, 0, fmt.Errorf("invalid cart id provided, %v", err)
		}

		getCart := s.cache.Get(context.Background(), cartUUID.String())
		if getCart.Err() == redis.Nil {
			return nil, 0, err
		}

		var cartBytes []byte
		cartBytes, err = getCart.Bytes()
		if err != nil {
			return nil, 0, err
		}

		if len(cartBytes) != 0 {
			err = json.Unmarshal(cartBytes, &cart)
			if err != nil {
				return nil, 0, err
			}
		}
	}

	for _, e := range cart {
		sum += e.PriceForQuantity
	}

	return cart, sum, err
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
