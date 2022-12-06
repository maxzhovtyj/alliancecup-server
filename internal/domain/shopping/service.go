package shopping

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/go-redis/redis/v9"
	"github.com/google/uuid"
	"github.com/zh0vtyj/allincecup-server/internal/domain/models"
	"time"
)

const (
	userCartCacheTTL       = 72 * time.Hour
	userFavouritesCacheTTL = 72 * time.Hour
)

type Service interface {
	NewCart(userId int) (uuid.UUID, error)
	GetCart(cartId string) ([]CartProduct, float64, error)
	AddToCart(info CartProduct, cartId string, userId int) error
	DeleteFromCart(productId, userId int, cartId string) error
	NewFavourites(userId int) (favouritesUUID uuid.UUID, err error)
	AddToFavourites(product models.Product, favouritesId string, userId int) error
	GetFavourites(userFavouritesId string) (favourites []models.Product, err error)
	DeleteFromFavourites(userFavouritesId string, userId, productId int) error
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

func validateNoDuplicateInCart(cart []CartProduct, product CartProduct) error {
	for _, p := range cart {
		if p.Id == product.Id {
			return fmt.Errorf("duplicate, product already in cart")
		}
	}

	return nil
}

func validateNoDuplicateInFavourites(products []models.Product, productId int) error {
	for _, p := range products {
		if p.Id == productId {
			return fmt.Errorf("duplicate, product already in favourites")
		}
	}

	return nil
}

func remove(slice []CartProduct, s int) []CartProduct {
	return append(slice[:s], slice[s+1:]...)
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

func (s *service) AddToCart(info CartProduct, cartId string, userId int) error {
	price, err := s.repo.PriceValidation(info.Id, info.Quantity)
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

func (s *service) GetCart(cartId string) (cart []CartProduct, sum float64, err error) {
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

func (s *service) DeleteFromCart(productId, userId int, cartId string) error {
	getCart := s.cache.Get(context.Background(), cartId)
	if getCart.Err() != nil {
		if getCart.Err() == redis.Nil {
			return fmt.Errorf("failed to find cart with id = %s", cartId)
		}
		return getCart.Err()
	}

	cartBytes, err := getCart.Bytes()
	if err != nil {
		return err
	}

	var cartProducts []CartProduct
	if len(cartBytes) != 0 {
		err = json.Unmarshal(cartBytes, &cartProducts)
		if err != nil {
			return fmt.Errorf("failed to unmarshal byte to cart product model")
		}
	}

	for i, p := range cartProducts {
		if p.Id == productId {
			cartProducts = remove(cartProducts, i)
		}
	}

	cartProductsMarshal, err := json.Marshal(cartProducts)
	if err != nil {
		return err
	}

	setCart := s.cache.Set(context.Background(), cartId, cartProductsMarshal, userCartCacheTTL)
	if setCart.Err() != nil {
		return err
	}

	if userId != 0 {
		err = s.repo.DeleteFromCart(productId, userId)
		if err != nil {
			return fmt.Errorf("failed to delete product in cart, %v", err)
		}
	}

	return err
}

func (s *service) NewFavourites(userId int) (favouritesUUID uuid.UUID, err error) {
	favouritesUUID = uuid.New()

	var productsBytes []byte
	if userId != 0 {
		products, err := s.repo.GetFavourites(userId)
		if err != nil {
			return [16]byte{}, err
		}

		productsBytes, err = json.Marshal(products)
		if err != nil {
			return [16]byte{}, err
		}
	}

	set := s.cache.Set(context.Background(), favouritesUUID.String(), productsBytes, userFavouritesCacheTTL)
	if set.Err() != nil {
		return [16]byte{}, fmt.Errorf("failed to create favourites in cache")
	}

	return favouritesUUID, nil
}

func (s *service) AddToFavourites(product models.Product, favouritesId string, userId int) error {
	getFavourites := s.cache.Get(context.Background(), favouritesId)
	if getFavourites.Err() != nil {
		if getFavourites.Err() == redis.Nil {
			return fmt.Errorf("failed to find favourites in cache id, %v", getFavourites.Err())
		}

		return fmt.Errorf("failed to get favourites from cache, %v", getFavourites.Err())
	}

	favouritesBytes, err := getFavourites.Bytes()
	if err != nil {
		return err
	}

	var favouriteProducts []models.Product

	if len(favouritesBytes) != 0 {
		err = json.Unmarshal(favouritesBytes, &favouriteProducts)
		if err != nil {
			return err
		}
	}

	err = validateNoDuplicateInFavourites(favouriteProducts, product.Id)
	if err != nil {
		return err
	}

	favouriteProducts = append(favouriteProducts, product)
	favouriteProductsMarshal, err := json.Marshal(favouriteProducts)
	if err != nil {
		return err
	}

	setFavourites := s.cache.Set(context.Background(), favouritesId, favouriteProductsMarshal, userFavouritesCacheTTL)
	if setFavourites.Err() != nil {
		return fmt.Errorf("failed to update favourite products, %v", setFavourites.Err())
	}

	if userId != 0 {
		err = s.repo.AddToFavourites(userId, product.Id)
		if err != nil {
			return err
		}
	}

	return err
}

func (s *service) GetFavourites(userFavouritesId string) (favourites []models.Product, err error) {
	var favouritesUUID uuid.UUID

	if userFavouritesId != "" {
		favouritesUUID, err = uuid.Parse(userFavouritesId)
		if err != nil {
			return nil, fmt.Errorf("invalid favourites id provided, %v", err)
		}

		getFavourites := s.cache.Get(context.Background(), favouritesUUID.String())
		if err = getFavourites.Err(); err != nil {
			return nil, fmt.Errorf("failed to find favourites id in cache, %v", err)
		}

		var favouritesBytes []byte
		favouritesBytes, err = getFavourites.Bytes()
		if err != nil {
			return nil, err
		}

		if len(favouritesBytes) != 0 {
			err = json.Unmarshal(favouritesBytes, &favourites)
			if err != nil {
				return nil, err
			}
		}
	}

	return favourites, err
}

func (s *service) DeleteFromFavourites(userFavouritesId string, userId, productId int) error {
	getFavourites := s.cache.Get(context.Background(), userFavouritesId)
	if getFavourites.Err() != nil {
		return fmt.Errorf("failed to find user favourites id in cache, %v", getFavourites.Err())
	}

	favouritesBytes, err := getFavourites.Bytes()
	if err != nil {
		return fmt.Errorf("failed to get favourite products bytes, %v", err)
	}

	var favourites []models.Product

	if len(favouritesBytes) != 0 {
		err = json.Unmarshal(favouritesBytes, &favourites)
		if err != nil {
			return fmt.Errorf("failed to unmarshal favourite products, %v", err)
		}

		//for i, p := range favourites {
		//	if p.Id == productId {
		//		favourites = remove(favourites, i)
		//	}
		//}
	}

	favouritesMarshal, err := json.Marshal(favourites)
	if err != nil {
		return fmt.Errorf("failed to marshal updated favourite products, %v", err)
	}

	setFavourites := s.cache.Set(context.Background(), userFavouritesId, favouritesMarshal, userFavouritesCacheTTL)
	if setFavourites.Err() != nil {
		return fmt.Errorf("failed to update favourite products, %v", setFavourites.Err())
	}

	if userId != 0 {
		err = s.repo.DeleteFromFavourites(userId, productId)
		if err != nil {
			return err
		}
	}

	return err
}
