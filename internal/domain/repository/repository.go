package repository

import (
	sq "github.com/Masterminds/squirrel"
	"github.com/jmoiron/sqlx"
	"github.com/zh0vtyj/alliancecup-server/internal/domain/category"
	"github.com/zh0vtyj/alliancecup-server/internal/domain/inventory"
	"github.com/zh0vtyj/alliancecup-server/internal/domain/order"
	"github.com/zh0vtyj/alliancecup-server/internal/domain/product"
	"github.com/zh0vtyj/alliancecup-server/internal/domain/review"
	"github.com/zh0vtyj/alliancecup-server/internal/domain/shopping"
	"github.com/zh0vtyj/alliancecup-server/internal/domain/supply"
	"github.com/zh0vtyj/alliancecup-server/internal/domain/user"
	"github.com/zh0vtyj/alliancecup-server/pkg/logging"
)

type Repository struct {
	Authorization user.Storage
	Category      category.Storage
	Product       product.Storage
	Order         order.Storage
	Shopping      shopping.Storage
	Supply        supply.Storage
	Review        review.Storage
	Inventory     inventory.Storage
	logger        *logging.Logger
}

func NewRepository(db *sqlx.DB, logger *logging.Logger) *Repository {
	psql := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)

	return &Repository{
		Authorization: user.NewAuthPostgres(db, psql),
		Product:       product.NewProductsPostgres(db, psql),
		Category:      category.NewCategoryPostgres(db),
		Shopping:      shopping.NewShoppingPostgres(db, psql),
		Order:         order.NewOrdersPostgres(db, psql),
		Supply:        supply.NewSupplyPostgres(db, psql),
		Review:        review.NewReviewStorage(db, psql),
		Inventory:     inventory.NewInventoryStorage(db, psql, logger),
		logger:        logger,
	}
}
