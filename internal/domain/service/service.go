package service

import (
	"github.com/zh0vtyj/allincecup-server/internal/domain/category"
	"github.com/zh0vtyj/allincecup-server/internal/domain/order"
	"github.com/zh0vtyj/allincecup-server/internal/domain/product"
	"github.com/zh0vtyj/allincecup-server/internal/domain/repository"
	"github.com/zh0vtyj/allincecup-server/internal/domain/review"
	"github.com/zh0vtyj/allincecup-server/internal/domain/shopping"
	"github.com/zh0vtyj/allincecup-server/internal/domain/supply"
	"github.com/zh0vtyj/allincecup-server/internal/domain/user"
)

type Service struct {
	Authorization user.AuthorizationService
	Category      category.Service
	Product       product.Service
	Order         order.Service
	Shopping      shopping.Service
	Supply        supply.Service
	Review        review.Service
}

func NewService(repos *repository.Repository) *Service {
	return &Service{
		Authorization: user.NewAuthService(repos.Authorization),
		Product:       product.NewProductsService(repos.Product),
		Category:      category.NewCategoryService(repos.Category),
		Order:         order.NewOrdersService(repos.Order),
		Shopping:      shopping.NewShoppingService(repos.Shopping),
		Supply:        supply.NewSupplyService(repos.Supply),
		Review:        review.NewReviewService(repos.Review),
	}
}
