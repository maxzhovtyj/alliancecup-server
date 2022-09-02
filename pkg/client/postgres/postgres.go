package postgres

import (
	"fmt"
	"github.com/jmoiron/sqlx"
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
	OrdersTable               = "order"
	OrdersProductsTable       = "orders_products"
	OrdersDeliveryTable       = "orders_delivery"
	SupplyTable               = "supply"
	SupplyPaymentTable        = "supply_payment"
	SupplyProductsTable       = "supply_products"
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
