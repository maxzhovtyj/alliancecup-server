package inventory

import (
	sq "github.com/Masterminds/squirrel"
	"github.com/jmoiron/sqlx"
	"github.com/zh0vtyj/allincecup-server/pkg/client/postgres"
	"github.com/zh0vtyj/allincecup-server/pkg/logging"
)

var inventoryProductsColumn = []string{
	"inventory.id",
	"inventory_products.inventory_id",
	"inventory_products.product_id",
	"inventory_products.last_inventory",
	"inventory_products.initial_amount",
	"inventory_products.supply",
	"inventory_products.spends",
	"inventory_products.write_off",
	"inventory_products.write_off_price",
	"inventory_products.planned_amount",
	"inventory_products.difference",
	"inventory_products.difference_price",
	"products.product_title",
	"products.last_inventory",
}

type Storage interface {
	GetProducts()
}

type storage struct {
	db     *sqlx.DB
	logger *logging.Logger
}

func NewInventoryStorage(db *sqlx.DB, logger *logging.Logger) Storage {
	return &storage{
		db:     db,
		logger: logger,
	}
}

func (s *storage) GetProducts() {
	psql := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)

	querySelectProducts, args, err := psql.
		Select(inventoryProductsColumn...).
		Join(postgres.ProductsTable + " ON inventory_products.product_id=products.id").
		Join(postgres.InventoryTable + " ON inventory_products.inventory_id=inventory.id").
		From(postgres.InventoryProductsTable).
		ToSql()

	//s.db.Select()

	s.logger.Println(querySelectProducts, args, err)
}
