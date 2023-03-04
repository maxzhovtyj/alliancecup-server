package inventory

import (
	"github.com/go-redis/redis/v9"
	"github.com/zh0vtyj/alliancecup-server/pkg/logging"
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
	return s.repo.GetProducts()
}

func (s *service) Save(products []CurrentProductDTO) error {
	return s.repo.Save(products)
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
