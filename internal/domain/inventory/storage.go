package inventory

import (
	"fmt"
	sq "github.com/Masterminds/squirrel"
	"github.com/jmoiron/sqlx"
	"github.com/zh0vtyj/allincecup-server/pkg/client/postgres"
	"github.com/zh0vtyj/allincecup-server/pkg/logging"
)

//var inventoryProductsColumn = []string{
//	"inventory.id",
//	"inventory_products.inventory_id",
//	"inventory_products.product_id",
//	"inventory_products.last_inventory",
//	"inventory_products.initial_amount",
//	"inventory_products.supply",
//	"inventory_products.spends",
//	"inventory_products.write_off",
//	"inventory_products.write_off_price",
//	"inventory_products.planned_amount",
//	"inventory_products.difference",
//	"inventory_products.difference_price",
//	"products.product_title",
//	"products.last_inventory",
//}

var inventoryProductsColumn = []string{
	"products.id",
	"products.product_title",
	"products.current_write_off",
	"products.current_spend",
	"products.current_supply",
	"products.amount_in_stock",
	"products.last_inventory",
}

type Storage interface {
	GetProducts() ([]CurrentProductDTO, error)
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

func (s *storage) GetProducts() ([]CurrentProductDTO, error) {
	psql := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)

	// TODO select amount_in_stock from the last inventory

	var products []CurrentProductDTO
	querySelectProducts, args, err := psql.
		Select(inventoryProductsColumn...).
		Join(postgres.InventoryProductsTable + " ON products.id = inventory_products.product_id").
		From(postgres.ProductsTable).
		ToSql()

	err = s.db.Select(&products, querySelectProducts, args...)
	if err != nil {
		return nil, err
	}

	fmt.Println(products)

	s.logger.Println(querySelectProducts, args, err)

	return products, err
}
