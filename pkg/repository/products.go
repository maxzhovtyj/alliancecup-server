package repository

import (
	server "allincecup-server"
	"fmt"
	"github.com/jmoiron/sqlx"
)

type ProductsPostgres struct {
	db *sqlx.DB
}

func NewProductsPostgres(db *sqlx.DB) *ProductsPostgres {
	return &ProductsPostgres{db: db}
}

func (p *ProductsPostgres) AddProduct(product server.Product) (int, error) {
	var id int
	query :=
		fmt.Sprintf(
			"INSERT INTO %s (category_id, title, price, size, characteristic, description, amount, in_stock) values ($1, $2, $3, $4, $5, $6, $7, $8) RETURNING id",
			productsTable)
	row :=
		p.db.QueryRow(query, product.CategoryId, product.Title, product.Price, product.Size, product.Characteristic, product.Description, product.Amount, product.InStock)

	if err := row.Scan(&id); err != nil {
		return 0, err
	}

	return id, nil
}
