package product

import (
	"fmt"
	sq "github.com/Masterminds/squirrel"
	"github.com/jmoiron/sqlx"
	server "github.com/zh0vtyj/allincecup-server/internal/shopping"
	"github.com/zh0vtyj/allincecup-server/pkg/client/postgres"
	"strconv"
	"strings"
)

type Storage interface {
	Search(searchInput string) ([]Product, error)
	GetWithParams(params server.SearchParams) ([]Product, error)
	GetProductById(id int) (Description, error)
	AddProduct(product Product, info []Info) (int, error)
	GetFavourites(userId int) ([]Product, error)
	Update(product Description) (int, error)
	Delete(productId int) error
}

type storage struct {
	db *sqlx.DB
}

func NewProductsPostgres(db *sqlx.DB) *storage {
	return &storage{db: db}
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

func (s *storage) GetWithParams(params server.SearchParams) ([]Product, error) {
	psql := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)

	query := psql.Select(productsColumnsSelect...).
		From(postgres.ProductsTable).
		LeftJoin(postgres.CategoriesTable + " ON products.category_id=categories.id").
		LeftJoin(postgres.ProductTypesTable + " ON products_types.id=products.type_id")

	if params.Characteristic != "" {
		query = query.Join(postgres.ProductsInfoTable + " ON products.id=products_info.product_id").
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
		searchToLower := strings.ToLower(params.Search)
		query = query.Where(sq.Like{"LOWER(products.product_title)": "%" + searchToLower + "%"})
	}

	if params.CreatedAt != "" {
		query = query.Where(sq.Lt{"products.created_at": params.CreatedAt})
	}

	if params.CategoryId != 0 {
		query = query.Where("products.category_id=?", params.CategoryId)
	}

	ordered := query.OrderBy("products.created_at DESC").Limit(9)

	querySql, args, err := ordered.ToSql()
	if err != nil {
		return nil, err
	}

	var products []Product
	err = s.db.Select(&products, querySql, args...)

	return products, err
}

func (s *storage) Search(searchInput string) ([]Product, error) {
	var products []Product
	psql := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)
	querySearch, args, err := psql.Select(productsColumnsSelect...).
		From(postgres.ProductsTable).
		LeftJoin(postgres.CategoriesTable + " ON categories.id=products.category_id").
		LeftJoin(postgres.ProductTypesTable + " ON products_types.id=products.type_id").
		Where(sq.Like{"products.product_title": "%" + searchInput + "%"}).
		ToSql()

	err = s.db.Select(&products, querySearch, args...)
	return products, err
}

func (s *storage) AddProduct(product Product, info []Info) (int, error) {
	tx, err := s.db.Begin()

	var categoryId int
	queryGetCategoryId := fmt.Sprintf("SELECT id FROM %s WHERE category_title=$1", postgres.CategoriesTable)
	categoryRow := tx.QueryRow(queryGetCategoryId, product.CategoryTitle)
	if err = categoryRow.Scan(&categoryId); err != nil {
		_ = tx.Rollback()
		return 0, err
	}

	var typeId int
	queryGetTypeId := fmt.Sprintf("SELECT id FROM %s WHERE type_title=$1", postgres.ProductTypesTable)
	typeRow := tx.QueryRow(queryGetTypeId, product.TypeTitle)
	if err = typeRow.Scan(&typeId); err != nil {
		_ = tx.Rollback()
		return 0, err
	}

	var productId int
	queryInsertProduct := fmt.Sprintf(
		"INSERT INTO %s (article, category_id, product_title, img_url, type_id, amount_in_stock, price, units_in_package, packages_in_box) values ($1, $2, $3, $4, $5, $6, $7, $8, $9) RETURNING id",
		postgres.ProductsTable,
	)
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
		queryInsertInfo := fmt.Sprintf(
			"INSERT INTO %s (product_id, info_title, description) values ($1, $2, $3)",
			postgres.ProductsInfoTable,
		)
		_, err = tx.Exec(queryInsertInfo, productId, info[i].InfoTitle, info[i].Description)
		if err != nil {
			_ = tx.Rollback()
			return 0, err
		}
	}

	return productId, tx.Commit()
}

func (s *storage) Update(product Description) (int, error) {
	tx, _ := s.db.Begin()

	var newCategoryId int
	queryGetCategoryId := fmt.Sprintf("SELECT id FROM %s WHERE category_title=$1 LIMIT 1", postgres.CategoriesTable)
	if err := s.db.Get(&newCategoryId, queryGetCategoryId, product.Info.CategoryTitle); err != nil {
		_ = tx.Rollback()
		return 0, err
	}

	var newTypeId int
	queryGetTypeId := fmt.Sprintf("SELECT id FROM %s WHERE type_title=$1 LIMIT 1", postgres.ProductTypesTable)
	if err := s.db.Get(&newTypeId, queryGetTypeId, product.Info.TypeTitle); err != nil {
		_ = tx.Rollback()
		return 0, err
	}

	queryUpdateProduct := fmt.Sprintf(
		"UPDATE %s SET article=$1, category_id=$2, product_title=$3, img_url=$4, type_id=$5, amount_in_stock=$6, price=$7, units_in_package=$8, packages_in_box=$9 WHERE id=$10",
		postgres.ProductsTable)
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

	queryDeleteOldDescription := fmt.Sprintf("DELETE FROM %s WHERE product_id=$1", postgres.ProductsInfoTable)
	_, err = tx.Exec(queryDeleteOldDescription, product.Info.Id)
	if err != nil {
		_ = tx.Rollback()
		return 0, err
	}

	for i := range product.Description {
		queryInsertInfo := fmt.Sprintf("INSERT INTO %s (product_id, info_title, description) values ($1, $2, $3)", postgres.ProductsInfoTable)
		_, err = tx.Exec(queryInsertInfo, product.Info.Id, product.Description[i].InfoTitle, product.Description[i].Description)
		if err != nil {
			_ = tx.Rollback()
			return 0, err
		}
	}

	return product.Info.Id, tx.Commit()
}

func (s *storage) Delete(productId int) error {
	tx, _ := s.db.Begin()

	queryDeleteProduct := fmt.Sprintf("DELETE FROM %s WHERE id=$1", postgres.ProductsTable)
	_, err := tx.Exec(queryDeleteProduct, productId)
	if err != nil {
		_ = tx.Rollback()
		return err
	}

	queryDeleteProductDescription := fmt.Sprintf("DELETE FROM %s WHERE product_id=$1", postgres.ProductsInfoTable)
	_, err = tx.Exec(queryDeleteProductDescription, productId)
	if err != nil {
		_ = tx.Rollback()
		return err
	}

	return tx.Commit()
}

func (s *storage) GetProductById(id int) (Description, error) {
	var product Description

	psql := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)

	queryGetProduct, args, err := psql.Select(productsColumnsSelect...).
		From(postgres.ProductsTable).
		LeftJoin(postgres.CategoriesTable + " ON categories.id=products.category_id").
		LeftJoin(postgres.ProductTypesTable + " ON products_types.id=products.type_id").
		Where(sq.Eq{"products.id": id}).ToSql()

	err = s.db.Get(&product.Info, queryGetProduct, args...)
	if err != nil {
		return Description{}, err
	}

	queryGetProductInfo := fmt.Sprintf("SELECT * FROM %s WHERE product_id=$1", postgres.ProductsInfoTable)
	err = s.db.Select(&product.Description, queryGetProductInfo, id)
	if err != nil {
		return Description{}, err
	}

	return product, err
}

func (s *storage) GetFavourites(userId int) ([]Product, error) {
	var products []Product

	queryGetFavourites := fmt.Sprintf(
		"SELECT products.id, products.article, products.product_title, products.img_url, products.price, products.units_in_package, products.packages_in_box, products.amount_in_stock, products.created_at FROM %s, %s WHERE user_id=$1 AND %s.product_id=%s.id",
		postgres.FavouritesTable,
		postgres.ProductsTable,
		postgres.FavouritesTable,
		postgres.ProductsTable,
	)

	if err := s.db.Select(&products, queryGetFavourites, userId); err != nil {
		return nil, err
	}

	return products, nil
}
