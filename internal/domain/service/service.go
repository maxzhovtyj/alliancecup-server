package service

import (
	"github.com/go-redis/redis/v9"
	"github.com/minio/minio-go/v7"
	"github.com/zh0vtyj/allincecup-server/internal/domain/category"
	"github.com/zh0vtyj/allincecup-server/internal/domain/inventory"
	"github.com/zh0vtyj/allincecup-server/internal/domain/order"
	"github.com/zh0vtyj/allincecup-server/internal/domain/product"
	"github.com/zh0vtyj/allincecup-server/internal/domain/repository"
	"github.com/zh0vtyj/allincecup-server/internal/domain/review"
	"github.com/zh0vtyj/allincecup-server/internal/domain/shopping"
	"github.com/zh0vtyj/allincecup-server/internal/domain/supply"
	"github.com/zh0vtyj/allincecup-server/internal/domain/user"
	"github.com/zh0vtyj/allincecup-server/pkg/logging"
)

type Service struct {
	Authorization user.AuthorizationService
	Category      category.Service
	Product       product.Service
	Order         order.Service
	Shopping      shopping.Service
	Supply        supply.Service
	Review        review.Service
	Inventory     inventory.Service
	logger        *logging.Logger
}

func NewService(
	repos *repository.Repository,
	logger *logging.Logger,
	cache *redis.Client,
	fileStorage *minio.Client,
) *Service {
	return &Service{
		Authorization: user.NewAuthService(repos.Authorization),
		Product:       product.NewProductsService(repos.Product, fileStorage),
		Category:      category.NewCategoryService(repos.Category, fileStorage),
		Order:         order.NewOrdersService(repos.Order, repos.Product),
		Shopping:      shopping.NewShoppingService(repos.Shopping),
		Supply:        supply.NewSupplyService(repos.Supply),
		Review:        review.NewReviewService(repos.Review),
		Inventory:     inventory.NewInventoryService(repos.Inventory, logger),
		logger:        logger,
	}
}
