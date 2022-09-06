package postgres

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/zh0vtyj/allincecup-server/internal/config"
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
	ProductTypesTable         = "products_types"
	ProductsInfoTable         = "products_info"
	FavouritesTable           = "favourites"
	DeliveryTypesTable        = "delivery_types"
	PaymentTypesTable         = "payment_types"
	OrdersTable               = "orders"
	OrdersProductsTable       = "orders_products"
	OrdersDeliveryTable       = "orders_delivery"
	SupplyTable               = "supply"
	SupplyPaymentTable        = "supply_payment"
	SupplyProductsTable       = "supply_products"
	ProductsReviewTable       = "products_review"
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
