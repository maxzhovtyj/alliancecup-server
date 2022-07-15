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

func (p *ProductsPostgres) AddProduct(product server.Product, info []server.ProductInfo) (int, error) {
	tx, err := p.db.Begin()

	var categoryId int
	queryGetCategoryId := fmt.Sprintf("SELECT id FROM %s WHERE category_title=$1", categoriesTable)
	categoryRow := tx.QueryRow(queryGetCategoryId, product.CategoryTitle)
	if err = categoryRow.Scan(&categoryId); err != nil {
		_ = tx.Rollback()
		return 0, err
	}

	var typeId int
	queryGetTypeId := fmt.Sprintf("SELECT id FROM %s WHERE type_title=$1", productTypesTable)
	typeRow := tx.QueryRow(queryGetTypeId, product.TypeTitle)
	if err = typeRow.Scan(&typeId); err != nil {
		_ = tx.Rollback()
		return 0, err
	}

	var productId int
	queryInsertProduct := fmt.Sprintf("INSERT INTO %s (article, category_id, product_title, img_url, type_id, amount_in_stock, price, units_in_package, packages_in_box) values ($1, $2, $3, $4, $5, $6, $7, $8, $9) RETURNING id", productsTable)
	insertRow := tx.QueryRow(queryInsertProduct,
		product.Article,
		categoryId,
		product.ProductTitle,
		product.ImgUrl,
		typeId,
		product.AmountInStock,
		product.Price,
		product.UnitsInPackage,
		product.PackagesInBox,
	)

	if err = insertRow.Scan(&productId); err != nil {
		_ = tx.Rollback()
		return 0, err
	}

	for i := range info {
		queryInsertInfo := fmt.Sprintf("INSERT INTO %s (product_id, info_title, description) values ($1, $2, $3)", productsInfoTable)
		_, err = tx.Exec(queryInsertInfo, productId, info[i].InfoTitle, info[i].Description)
		if err != nil {
			_ = tx.Rollback()
			return 0, err
		}
	}

	return productId, tx.Commit()
}
