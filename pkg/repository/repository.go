package repository

import (
	"github.com/jmoiron/sqlx"
	"github.com/zh0vtyj/allincecup-server/internal/category"
	"github.com/zh0vtyj/allincecup-server/internal/order"
	"github.com/zh0vtyj/allincecup-server/internal/product"
	"github.com/zh0vtyj/allincecup-server/internal/shopping"
	"github.com/zh0vtyj/allincecup-server/internal/supply"
	"github.com/zh0vtyj/allincecup-server/internal/user"
)

type Repository struct {
	Authorization user.AuthorizationStorage
	Category      category.Storage
	Product       product.Storage
	Order         order.Storage
	Shopping      shopping.Storage
	Supply        supply.Storage
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{
		Authorization: user.NewAuthPostgres(db),
		Product:       product.NewProductsPostgres(db),
		Category:      category.NewCategoryPostgres(db),
		Shopping:      shopping.NewShoppingPostgres(db),
		Order:         order.NewOrdersPostgres(db),
		Supply:        supply.NewSupplyPostgres(db),
	}
}
