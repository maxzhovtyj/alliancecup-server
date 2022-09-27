package inventory

import "github.com/zh0vtyj/allincecup-server/pkg/logging"

type Service interface {
	Products() ([]CurrentProductDTO, error)
	New(dto []InsertProductDTO) error
	GetAll(createdAt string) ([]DTO, error)
	GetInventoryProducts(id int) ([]SelectProductDTO, error)
}

type service struct {
	repo   Storage
	logger *logging.Logger
}

func NewInventoryService(repo Storage, logger *logging.Logger) Service {
	return &service{
		repo:   repo,
		logger: logger,
	}
}

func (s *service) Products() ([]CurrentProductDTO, error) {
	products, err := s.repo.GetProducts()
	if err != nil {
		return nil, err
	}

	return products, err
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
