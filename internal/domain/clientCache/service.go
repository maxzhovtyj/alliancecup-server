package clientCache

import (
	"github.com/go-redis/redis/v9"
	"github.com/zh0vtyj/allincecup-server/pkg/logging"
)

type Service interface {
	NewUser() error
}

type service struct {
	cache  *redis.Client
	logger *logging.Logger
}

func NewService(cache *redis.Client, logger *logging.Logger) Service {
	return &service{
		cache:  cache,
		logger: logger,
	}
}

func (s *service) NewUser() error {
	return nil
}
