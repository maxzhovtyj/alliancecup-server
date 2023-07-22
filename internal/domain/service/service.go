package service

import (
	"github.com/go-redis/redis/v9"
	"github.com/minio/minio-go/v7"
	"github.com/zh0vtyj/alliancecup-server/internal/config"
	"github.com/zh0vtyj/alliancecup-server/internal/domain/category"
	"github.com/zh0vtyj/alliancecup-server/internal/domain/inventory"
	"github.com/zh0vtyj/alliancecup-server/internal/domain/order"
	"github.com/zh0vtyj/alliancecup-server/internal/domain/product"
	"github.com/zh0vtyj/alliancecup-server/internal/domain/repository"
	"github.com/zh0vtyj/alliancecup-server/internal/domain/review"
	"github.com/zh0vtyj/alliancecup-server/internal/domain/shopping"
	"github.com/zh0vtyj/alliancecup-server/internal/domain/supply"
	"github.com/zh0vtyj/alliancecup-server/internal/domain/user"
	"github.com/zh0vtyj/alliancecup-server/pkg/logging"
	"github.com/zh0vtyj/alliancecup-server/pkg/telegram"
)

type Service struct {
	Authorization user.Service
	Category      category.Service
	Product       product.Service
	Order         order.Service
	Shopping      shopping.Service
	Supply        supply.Service
	Review        review.Service
	Inventory     inventory.Service
}

func New(
	repos *repository.Repository,
	auth config.Auth,
	logger *logging.Logger,
	cache *redis.Client,
	fileStorage *minio.Client,
	tgBotManager telegram.Manager,
) *Service {
	return &Service{
		Authorization: user.NewAuthService(repos.Authorization, auth),
		Product:       product.NewProductsService(repos.Product, fileStorage),
		Category:      category.NewCategoryService(repos.Category, fileStorage),
		Order:         order.NewOrdersService(repos.Order, repos.Product, cache, tgBotManager),
		Shopping:      shopping.NewShoppingService(repos.Shopping, cache),
		Supply:        supply.NewSupplyService(repos.Supply),
		Review:        review.NewReviewService(repos.Review),
		Inventory:     inventory.NewInventoryService(repos.Inventory, cache, logger),
	}
}
