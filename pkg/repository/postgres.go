package repository

import (
	"fmt"
	"github.com/jmoiron/sqlx"
)

const (
	usersTable                = "users"
	rolesTable                = "roles"
	categoriesTable           = "categories"
	categoriesFiltrationTable = "categories_filtration"
	productsTable             = "products"
	ordersTable               = "orders"
	sessionsTable             = "sessions"
	cartsTable                = "carts"
	productTypesTable         = "products_types"
	productsInfoTable         = "products_info"
	cartsProductsTable        = "carts_products"
	favouritesTable           = "favourites"
	ordersProductsTable       = "orders_products"
	ordersDeliveryTable       = "orders_delivery"
	deliveryTypesTable        = "delivery_types"
	paymentTypesTable         = "payment_types"
	supplyTable               = "supply"
	supplyPaymentTable        = "supply_payment"
	productsSupplyTable       = "supply_products"
)

type Config struct {
	Host     string
	Port     string
	Username string
	Password string
	DBName   string
	SSLMode  string
}

func NewPostgresDB(cfg Config) (*sqlx.DB, error) {
	db, err := sqlx.Open("postgres", fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s sslmode=%s",
		cfg.Host, cfg.Port, cfg.Username, cfg.DBName, cfg.Password, cfg.SSLMode))
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	return db, nil
}
