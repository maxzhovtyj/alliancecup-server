package inventory

import (
	"fmt"
	sq "github.com/Masterminds/squirrel"
	"github.com/jmoiron/sqlx"
	"github.com/zh0vtyj/alliancecup-server/pkg/client/postgres"
	"github.com/zh0vtyj/alliancecup-server/pkg/logging"
	"time"
)

type Storage interface {
	GetProducts() ([]CurrentProductDTO, error)
	DoInventory(products []InsertProductDTO) error
	GetInventories(createdAt string) ([]DTO, error)
	getInventoryProductsById(inventoryId int) ([]SelectProductDTO, error)
	Save(products []CurrentProductDTO) error
}

type storage struct {
	db     *sqlx.DB
	logger *logging.Logger
	qb     sq.StatementBuilderType
}

func NewInventoryStorage(db *sqlx.DB, psql sq.StatementBuilderType, logger *logging.Logger) Storage {
	return &storage{
		db:     db,
		logger: logger,
		qb:     psql,
	}
}

var ProductsToInventory = []string{
	"products.id",
	"products.product_title",
	"products.price as product_price",
	"products.current_write_off",
	"products.current_write_off * products.price as write_off_price",
	"products.current_spend",
	"products.current_supply",
	"products.current_real_amount",
	"products.current_real_amount * products.price as real_amount_price",
	"products.current_real_amount - products.amount_in_stock as difference",
	"(products.current_real_amount - products.amount_in_stock) * products.price as difference_price",
	"products.amount_in_stock",
	"products.last_inventory_id",                       // last inventory id
	"inventory.created_at as last_inventory",           // last inventory time
	"inventory_products.real_amount as initial_amount", // amount from the last inventory
}

func (s *storage) GetProducts() ([]CurrentProductDTO, error) {
	var products []CurrentProductDTO
	querySelectProducts, args, err := s.qb.
		Select(ProductsToInventory...).
		LeftJoin(postgres.InventoryTable + " ON products.last_inventory_id = inventory.id").
		LeftJoin(postgres.InventoryProductsTable + " ON products.last_inventory_id = inventory_products.inventory_id").
		From(postgres.ProductsTable).
		ToSql()

	err = s.db.Select(&products, querySelectProducts, args...)
	if err != nil {
		return nil, err
	}

	return products, err
}

func (s *storage) DoInventory(products []InsertProductDTO) error {
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

	queryUpdateProduct := fmt.Sprintf(
		`
		UPDATE %s
		SET
			current_supply = 0,
			current_spend = 0,
			current_write_off = 0,
			current_real_amount = 0,
			last_inventory_id = $1
		WHERE id = $2
		`,
		postgres.ProductsTable,
	)

	// inserting essential info
	for _, p := range products {
		sql, args, _ := s.qb.Insert(postgres.InventoryProductsTable).Values(
			inventoryId,
			p.ProductId,
			p.ProductPrice,
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
			s.logger.Println(sql, args)
			return fmt.Errorf("failed to insert inventory product, %v", err)
		}

		_, err = tx.Exec(queryUpdateProduct, inventoryId, p.ProductId)
		if err != nil {
			_ = tx.Rollback()
			return err
		}
	}

	return tx.Commit()
}

func (s *storage) GetInventories(createdAt string) ([]DTO, error) {
	var inventories []DTO

	selectInventoryQuery := s.qb.Select("id, created_at").From(postgres.InventoryTable)

	if createdAt != "" {
		selectInventoryQuery = selectInventoryQuery.Where(sq.Lt{"created_at": createdAt})
	}

	selectInventoryQuerySQL, args, err := selectInventoryQuery.
		OrderBy("created_at DESC").
		Limit(12).
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("failed to build query to select inventory due to: %v", err)
	}

	err = s.db.Select(&inventories, selectInventoryQuerySQL, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to select inventories from db due to: %v", err)
	}

	return inventories, err
}

var inventoryProducts = []string{
	"inventory_products.inventory_id",
	"inventory_products.product_id",
	"products.product_title",
	"inventory_products.product_price",
	"inventory_products.last_inventory_id",
	"inventory_products.initial_amount",
	"inventory_products.supply",
	"inventory_products.spend",
	"inventory_products.write_off",
	"inventory_products.write_off * inventory_products.product_price as write_off_price",
	"inventory_products.planned_amount",
	"inventory_products.real_amount",
	"inventory_products.real_amount * inventory_products.product_price as real_amount_price",
	"inventory_products.real_amount - inventory_products.planned_amount as difference",
	"(inventory_products.real_amount - inventory_products.planned_amount) * product_price as difference_price",
}

func (s *storage) getInventoryProductsById(inventoryId int) ([]SelectProductDTO, error) {
	var products []SelectProductDTO

	query, args, err := s.qb.
		Select(inventoryProducts...).
		Join(postgres.ProductsTable + " ON products.id = inventory_products.product_id").
		From(postgres.InventoryProductsTable).
		Where(sq.Eq{"inventory_id": inventoryId}).
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("failed to build query to select inventory products, %v", err)
	}

	err = s.db.Select(&products, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to select inventory products due to: %v", err)
	}

	return products, err
}

func (s *storage) Save(products []CurrentProductDTO) error {
	queryUpdateCurrentAmount := fmt.Sprintf(
		"UPDATE %s SET current_real_amount = $1 WHERE id = $2",
		postgres.ProductsTable,
	)

	for _, p := range products {
		_, err := s.db.Exec(queryUpdateCurrentAmount, p.RealAmount, p.ProductId)
		if err != nil {
			return fmt.Errorf("failed to update current product real amount, %v", err)
		}
	}

	return nil
}
