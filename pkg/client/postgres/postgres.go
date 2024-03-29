package postgres

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/zh0vtyj/alliancecup-server/internal/config"
)

const (
	UsersTable                = "users"
	RolesTable                = "roles"
	CartsTable                = "carts"
	CartsProductsTable        = "carts_products"
	SessionsTable             = "sessions"
	CategoriesTable           = "categories"
	CategoriesFiltrationTable = "categories_filtration"
	ProductsTable             = "products"
	FavouritesTable           = "favourites"
	DeliveryTypesTable        = "delivery_types"
	PaymentTypesTable         = "payment_types"
	OrdersTable               = "orders"
	OrdersProductsTable       = "orders_products"
	SupplyTable               = "supply"
	SupplyPaymentTable        = "supply_payment"
	SupplyProductsTable       = "supply_products"
	ProductsReviewTable       = "products_review"
	InventoryTable            = "inventory"
	InventoryProductsTable    = "inventory_products"
)

func NewPostgresDB(cfg config.Storage) (*sqlx.DB, error) {
	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s dbname=%s password=%s sslmode=%s",
		cfg.Host, cfg.Port, cfg.Username, cfg.DBName, cfg.Password, cfg.SSLMode,
	)

	db, err := sqlx.Open("postgres", dsn)
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	return db, nil
}
