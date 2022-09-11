package inventory

import "github.com/zh0vtyj/allincecup-server/pkg/logging"

type Service interface {
	Products()
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

func (s *service) Products() {
	s.repo.GetProducts()
}
