package repository

import (
	server "allincecup-server"
	"github.com/jmoiron/sqlx"
)

type ProductsPostgres struct {
	db *sqlx.DB
}

func NewProductsPostgres(db *sqlx.DB) *ProductsPostgres {
	return &ProductsPostgres{db: db}
}

func (p *ProductsPostgres) AddProduct(product server.Product) (int, error) {
	return 0, nil
}
