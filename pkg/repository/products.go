package repository

import (
	"fmt"
	sq "github.com/Masterminds/squirrel"
	"github.com/jmoiron/sqlx"
	server "github.com/zh0vtyj/allincecup-server/pkg/models"
	"strconv"
	"strings"
)

type ProductsPostgres struct {
	db *sqlx.DB
}

func NewProductsPostgres(db *sqlx.DB) *ProductsPostgres {
	return &ProductsPostgres{db: db}
}

var productsColumnsSelect = []string{
	"products.id",
	"products.article",
	"categories.category_title",
	"products.product_title",
	"products.img_url",
	"products_types.type_title",
	"products.amount_in_stock",
	"products.price",
	"products.units_in_package",
	"products.packages_in_box",
	"products.created_at",
}

func (p *ProductsPostgres) GetWithParams(params server.SearchParams) ([]server.Product, error) {
	psql := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)

	query := psql.Select(productsColumnsSelect...).
		From(productsTable).
		LeftJoin(categoriesTable + " ON products.category_id=categories.id").
		LeftJoin(productTypesTable + " ON products_types.id=products.type_id")

	if params.Characteristic != "" {
		query = query.Join(productsInfoTable + " ON products.id=products_info.product_id").
			Where(sq.Eq{"products_info.description": params.Characteristic})
		// TODO select all rows with exact match of characteristics slice
	}

	if params.PriceRange != "" {
		price := strings.Split(params.PriceRange, ":")
		gtPrice, err := strconv.ParseFloat(price[0], 64)
		ltPrice, err := strconv.ParseFloat(price[1], 64)
		if err != nil {
			return nil, err
		}
		query = query.Where("products.price BETWEEN ? AND ?", gtPrice, ltPrice)
	}

	if params.Search != "" {
		query = query.Where(sq.Like{"LOWER(products.product_title)": "%" + params.Search + "%"})
	}

	if params.CreatedAt != "" {
		query = query.Where(sq.Lt{"products.created_at": params.CreatedAt})
	}

	query = query.Where("products.category_id=?", params.CategoryId)

	ordered := query.OrderBy("products.created_at DESC").Limit(9)

	querySql, args, err := ordered.ToSql()
	if err != nil {
		return nil, err
	}

	var products []server.Product
	err = p.db.Select(&products, querySql, args...)

	return products, err
}

func (p *ProductsPostgres) Search(searchInput string) ([]server.Product, error) {
	var products []server.Product
	psql := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)
	querySearch, args, err := psql.Select(productsColumnsSelect...).
		From(productsTable).
		LeftJoin(categoriesTable + " ON categories.id=products.category_id").
		LeftJoin(productTypesTable + " ON products_types.id=products.type_id").
		Where(sq.Like{"products.product_title": "%" + searchInput + "%"}).
		ToSql()

	err = p.db.Select(&products, querySearch, args...)
	return products, err
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

func (p *ProductsPostgres) Update(product server.ProductInfoDescription) (int, error) {
	tx, _ := p.db.Begin()

	var newCategoryId int
	queryGetCategoryId := fmt.Sprintf("SELECT id FROM %s WHERE category_title=$1 LIMIT 1", categoriesTable)
	if err := p.db.Get(&newCategoryId, queryGetCategoryId, product.Info.CategoryTitle); err != nil {
		_ = tx.Rollback()
		return 0, err
	}

	var newTypeId int
	queryGetTypeId := fmt.Sprintf("SELECT id FROM %s WHERE type_title=$1 LIMIT 1", productTypesTable)
	if err := p.db.Get(&newTypeId, queryGetTypeId, product.Info.TypeTitle); err != nil {
		_ = tx.Rollback()
		return 0, err
	}

	queryUpdateProduct := fmt.Sprintf("UPDATE %s SET article=$1, category_id=$2, product_title=$3, img_url=$4, type_id=$5, amount_in_stock=$6, price=$7, units_in_package=$8, packages_in_box=$9 WHERE id=$10", productsTable)
	_, err := tx.Exec(
		queryUpdateProduct,
		product.Info.Article,
		newCategoryId,
		product.Info.ProductTitle,
		product.Info.ImgUrl,
		newTypeId,
		product.Info.AmountInStock,
		product.Info.Price,
		product.Info.UnitsInPackage,
		product.Info.PackagesInBox,
		product.Info.Id,
	)
	if err != nil {
		_ = tx.Rollback()
		return 0, err
	}

	queryDeleteOldDescription := fmt.Sprintf("DELETE FROM %s WHERE product_id=$1", productsInfoTable)
	_, err = tx.Exec(queryDeleteOldDescription, product.Info.Id)
	if err != nil {
		_ = tx.Rollback()
		return 0, err
	}

	for i := range product.Description {
		queryInsertInfo := fmt.Sprintf("INSERT INTO %s (product_id, info_title, description) values ($1, $2, $3)", productsInfoTable)
		_, err = tx.Exec(queryInsertInfo, product.Info.Id, product.Description[i].InfoTitle, product.Description[i].Description)
		if err != nil {
			_ = tx.Rollback()
			return 0, err
		}
	}

	return product.Info.Id, tx.Commit()
}

func (p *ProductsPostgres) Delete(productId int) error {
	tx, _ := p.db.Begin()

	queryDeleteProduct := fmt.Sprintf("DELETE FROM %s WHERE id=$1", productsTable)
	_, err := tx.Exec(queryDeleteProduct, productId)
	if err != nil {
		_ = tx.Rollback()
		return err
	}

	queryDeleteProductDescription := fmt.Sprintf("DELETE FROM %s WHERE product_id=$1", productsInfoTable)
	_, err = tx.Exec(queryDeleteProductDescription, productId)
	if err != nil {
		_ = tx.Rollback()
		return err
	}

	return tx.Commit()
}

func (p *ProductsPostgres) GetProductById(id int) (server.ProductInfoDescription, error) {
	var product server.ProductInfoDescription

	psql := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)

	queryGetProduct, args, err := psql.Select(productsColumnsSelect...).
		From(productsTable).
		LeftJoin(categoriesTable + " ON categories.id=products.category_id").
		LeftJoin(productTypesTable + " ON products_types.id=products.type_id").
		Where(sq.Eq{"products.id": id}).ToSql()

	err = p.db.Get(&product.Info, queryGetProduct, args...)
	if err != nil {
		return server.ProductInfoDescription{}, err
	}

	queryGetProductInfo := fmt.Sprintf("SELECT * FROM %s WHERE product_id=$1", productsInfoTable)
	err = p.db.Select(&product.Description, queryGetProductInfo, id)
	if err != nil {
		return server.ProductInfoDescription{}, err
	}

	return product, err
}
