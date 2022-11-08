package inventory

import (
	"context"
	"encoding/json"
	"github.com/go-redis/redis/v9"
	"github.com/zh0vtyj/allincecup-server/pkg/client/redisdb"
	"github.com/zh0vtyj/allincecup-server/pkg/logging"
)

type Service interface {
	Products() ([]CurrentProductDTO, error)
	New(dto []InsertProductDTO) error
	GetAll(createdAt string) ([]DTO, error)
	GetInventoryProducts(id int) ([]SelectProductDTO, error)
	Save(products []CurrentProductDTO) error
}

type service struct {
	repo   Storage
	cache  *redis.Client
	logger *logging.Logger
}

func NewInventoryService(repo Storage, cache *redis.Client, logger *logging.Logger) Service {
	return &service{
		repo:   repo,
		cache:  cache,
		logger: logger,
	}
}

func (s *service) Products() (products []CurrentProductDTO, err error) {
	get := s.cache.Get(context.Background(), redisdb.InventoryProducts)
	if get.Err() == redis.Nil {
		s.logger.Printf("Cache %s key does not exist", redisdb.InventoryProducts)
		products, err = s.repo.GetProducts()

		s.cache.Set(context.Background(), redisdb.InventoryProducts, products, 0)
		if err != nil {
			return nil, err
		}
	} else {
		jsonBytes, err := get.Bytes()
		err = json.Unmarshal(jsonBytes, &products)
		if err != nil {
			return nil, err
		}
	}

	return products, err
}

func (s *service) Save(products []CurrentProductDTO) error {
	productsBytes, err := json.Marshal(products)
	if err != nil {
		return err
	}

	set := s.cache.Set(context.Background(), redisdb.InventoryProducts, productsBytes, 0)
	if err = set.Err(); err != nil {
		return err
	}

	return nil
}

func (s *service) New(dto []InsertProductDTO) error {
	return s.repo.DoInventory(dto)
}

func (s *service) GetAll(createdAt string) ([]DTO, error) {
	return s.repo.GetInventories(createdAt)
}

func (s *service) GetInventoryProducts(id int) ([]SelectProductDTO, error) {
	return s.repo.getInventoryProductsById(id)
}
