package repository

import (
	server "allincecup-server"
	"github.com/jmoiron/sqlx"
)

type Authorization interface {
	CreateUser(user server.User) (int, error)
	GetUser(email string, password string) (server.User, error)
}

type ShopItemCup interface {
}

type ShopList interface {
}

type Categories interface {
}

type Products interface {
	AddProduct(product server.Product) (int, error)
}

type Repository struct {
	Authorization
	ShopItemCup
	ShopList
	Categories
	Products
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{
		Authorization: NewAuthPostgres(db),
		Products:      NewProductsPostgres(db),
		Categories:    NewCategoriesPostgres(db),
	}
}
