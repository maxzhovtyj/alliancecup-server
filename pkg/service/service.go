package service

import (
	"github.com/zh0vtyj/allincecup-server/internal/category"
	"github.com/zh0vtyj/allincecup-server/internal/order"
	"github.com/zh0vtyj/allincecup-server/internal/product"
	"github.com/zh0vtyj/allincecup-server/internal/shopping"
	"github.com/zh0vtyj/allincecup-server/internal/supply"
	"github.com/zh0vtyj/allincecup-server/internal/user"
	"github.com/zh0vtyj/allincecup-server/pkg/repository"
)

type Service struct {
	Authorization user.AuthorizationService
	Category      category.Service
	Product       product.Service
	Order         order.Service
	Shopping      shopping.Service
	Supply        supply.Service
}

func NewService(repos *repository.Repository) *Service {
	return &Service{
		Authorization: user.NewAuthService(repos.Authorization),
		Product:       product.NewProductsService(repos.Product),
		Category:      category.NewCategoryService(repos.Category),
		Order:         order.NewOrdersService(repos.Order),
		Shopping:      shopping.NewShoppingService(repos.Shopping),
		Supply:        supply.NewSupplyService(repos.Supply),
	}
}
