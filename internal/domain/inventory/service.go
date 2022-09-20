package inventory

import "github.com/zh0vtyj/allincecup-server/pkg/logging"

type Service interface {
	Products() ([]CurrentProductDTO, error)
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
