package repository

import (
	"github.com/jmoiron/sqlx"
	"github.com/zh0vtyj/allincecup-server/internal/domain/category"
	"github.com/zh0vtyj/allincecup-server/internal/domain/order"
	"github.com/zh0vtyj/allincecup-server/internal/domain/product"
	"github.com/zh0vtyj/allincecup-server/internal/domain/review"
	"github.com/zh0vtyj/allincecup-server/internal/domain/shopping"
	"github.com/zh0vtyj/allincecup-server/internal/domain/supply"
	"github.com/zh0vtyj/allincecup-server/internal/domain/user"
)

type Repository struct {
	Authorization user.AuthorizationStorage
	Category      category.Storage
	Product       product.Storage
	Order         order.Storage
	Shopping      shopping.Storage
	Supply        supply.Storage
	Review        review.Storage
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{
		Authorization: user.NewAuthPostgres(db),
		Product:       product.NewProductsPostgres(db),
		Category:      category.NewCategoryPostgres(db),
		Shopping:      shopping.NewShoppingPostgres(db),
		Order:         order.NewOrdersPostgres(db),
		Supply:        supply.NewSupplyPostgres(db),
		Review:        review.NewReviewStorage(db),
	}
}
