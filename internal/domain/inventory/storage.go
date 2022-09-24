package inventory

import (
	"fmt"
	sq "github.com/Masterminds/squirrel"
	"github.com/jmoiron/sqlx"
	"github.com/zh0vtyj/allincecup-server/pkg/client/postgres"
	"github.com/zh0vtyj/allincecup-server/pkg/logging"
	"time"
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
	"inventory_products.inventory_id as last_inventory_id",
	"inventory.created_at as last_inventory",
}

type Storage interface {
	GetProducts() ([]CurrentProductDTO, error)
	DoInventory(products []ProductDTO) error
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
	// TODO select if there are not inventories yet
	var products []CurrentProductDTO
	querySelectProducts, args, err := psql.
		Select(inventoryProductsColumn...).
		Join(postgres.InventoryProductsTable + " ON products.id = inventory_products.product_id").
		Join(postgres.InventoryTable + " ON inventory_products.inventory_id = inventory.id").
		From(postgres.ProductsTable).
		ToSql()

	fmt.Println(querySelectProducts)

	err = s.db.Select(&products, querySelectProducts, args...)
	if err != nil {
		return nil, err
	}

	return products, err
}

func (s *storage) DoInventory(products []ProductDTO) error {
	psql := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)

	tx, _ := s.db.Begin()

	// create new inventory and get its id
	var inventoryId int
	var createdAt time.Time
	queryNewInventory := fmt.Sprintf("INSERT INTO %s (created_at) values (now() at time zone 'utc-3') RETURNING id, created_at", postgres.InventoryTable)
	row := tx.QueryRow(queryNewInventory)
	if err := row.Scan(&inventoryId, &createdAt); err != nil {
		_ = tx.Rollback()
		return err
	}

	// inserting essential info
	for _, p := range products {
		sql, args, _ := psql.Insert(postgres.InventoryProductsTable).Values(
			inventoryId,
			p.ProductId,
			p.LastInventoryId,
			p.InitialAmount,
			p.Supply,
			p.Spend,
			p.WriteOff,
			p.PlannedAmount,
			p.RealAmount,
		).ToSql()

		_, err := tx.Exec(sql, args...)
		if err != nil {
			_ = tx.Rollback()
			return err
		}

		// update product fields
		// fmt.Sprintf("UPDATE %s SET ", postgres.ProductsTable)
	}

	return tx.Commit()
}
