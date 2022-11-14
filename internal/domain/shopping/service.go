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

const userCartCacheTTL = 72 * time.Hour

type Service interface {
	NewCart(userId int) (uuid.UUID, error)
	AddToCart(info CartProduct, cartId string, userId int) error
	GetCart(cartId string) ([]CartProductFullInfo, float64, error)
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

	set := s.cache.Set(context.Background(), cartUUID.String(), productsBytes, userCartCacheTTL)
	if set.Err() != nil {
		return [16]byte{}, fmt.Errorf("failed to create card in cache")
	}

	return cartUUID, nil
}

func validateNoDuplicateInCart(cart []CartProduct, product CartProduct) error {
	for _, p := range cart {
		if p.ProductId == product.ProductId {
			return fmt.Errorf("duplicate, product already in cart")
		}
	}

	return nil
}

func (s *service) AddToCart(info CartProduct, cartId string, userId int) error {
	price, err := s.repo.PriceValidation(info.ProductId, info.Quantity)
	if err != nil {
		return err
	}

	if price != info.PriceForQuantity {
		return errors.New("price mismatch")
	}

	var cartProducts []CartProduct

	cacheCart := s.cache.Get(context.Background(), cartId)
	cacheCartBytes, err := cacheCart.Bytes()
	if err != nil {
		return err
	}

	if len(cacheCartBytes) != 0 {
		err = json.Unmarshal(cacheCartBytes, &cartProducts)
		if err != nil {
			return err
		}
	}

	err = validateNoDuplicateInCart(cartProducts, info)
	if err != nil {
		return err
	}

	cartProducts = append(cartProducts, info)
	cartProductsBytes, err := json.Marshal(cartProducts)
	if err != nil {
		return err
	}

	updateUserCart := s.cache.Set(context.Background(), cartId, cartProductsBytes, userCartCacheTTL)
	if err = updateUserCart.Err(); err != nil {
		return fmt.Errorf("failed to update user cart in cache, %v", updateUserCart.Err())
	}

	if userId != 0 {
		err = s.repo.AddToCart(userId, info)
		if err != nil {
			return err
		}
	}

	return err
}

func (s *service) GetCart(cartId string) (cart []CartProductFullInfo, sum float64, err error) {
	var cartUUID uuid.UUID

	if cartId != "" {
		cartUUID, err = uuid.Parse(cartId)
		if err != nil {
			return nil, 0, fmt.Errorf("invalid cart id provided, %v", err)
		}

		getCart := s.cache.Get(context.Background(), cartUUID.String())
		if err = getCart.Err(); err != nil {
			return nil, 0, fmt.Errorf("failed to find cart id in cache, %v", err)
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
