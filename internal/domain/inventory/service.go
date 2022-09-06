package inventory

import "github.com/zh0vtyj/allincecup-server/pkg/logging"

type Service interface {
}

type service struct {
	repo   Storage
	logger *logging.Logger
}

func NewInventoryService(repo Storage, logger *logging.Logger) Service {
	return service{
		repo:   repo,
		logger: logger,
	}
}
