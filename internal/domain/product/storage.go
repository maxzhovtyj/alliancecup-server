package product

import (
	"fmt"
	sq "github.com/Masterminds/squirrel"
	"github.com/jmoiron/sqlx"
	server "github.com/zh0vtyj/allincecup-server/internal/domain/shopping"
	"github.com/zh0vtyj/allincecup-server/pkg/client/postgres"
	"strconv"
	"strings"
)

type Storage interface {
	Search(searchInput string) ([]Product, error)
	GetWithParams(params server.SearchParams) ([]Product, error)
	GetProductById(id int) (Product, error)
	Create(product Product) (int, error)
	GetFavourites(userId int) ([]Product, error)
	Update(product Product) (int, error)
	UpdateImage(product Product) (int, error)
	Delete(productId int) error
}

type storage struct {
	qb sq.StatementBuilderType
	db *sqlx.DB
}

func NewProductsPostgres(db *sqlx.DB, psql sq.StatementBuilderType) *storage {
	return &storage{
		db: db,
		qb: psql,
	}
}

var productsColumnsSelect = []string{
	"products.id",
	"products.article",
	"categories.category_title",
	"products.product_title",
	"products.img_url",
	"products.img_uuid",
	"products.amount_in_stock",
	"products.price",
	"products.characteristics",
	"products.packaging",
	"products.is_active",
	"products.created_at",
}

func (s *storage) GetWithParams(params server.SearchParams) ([]Product, error) {
	query := s.qb.Select(productsColumnsSelect...).
		From(postgres.ProductsTable).
		LeftJoin(postgres.CategoriesTable + " ON products.category_id = categories.id")

	if len(params.Characteristic) != 0 {
		for _, chr := range params.Characteristic {
			query = query.Where("characteristics ->> ? = ?", chr.Name, chr.Value)
		}
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
		query = query.Where("products.category_id = ?", params.CategoryId)
	}

	if params.IsActive != nil {
		query = query.Where("products.is_active = ?", params.IsActive)
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
	querySearch, args, err := s.qb.Select(productsColumnsSelect...).
		From(postgres.ProductsTable).
		LeftJoin(postgres.CategoriesTable + " ON categories.id = products.category_id").
		Where(sq.Like{"products.product_title": "%" + searchInput + "%"}).
		ToSql()

	err = s.db.Select(&products, querySearch, args...)
	return products, err
}

func (s *storage) Create(product Product) (int, error) {
	tx, err := s.db.Begin()

	var productId int
	queryInsertProduct := fmt.Sprintf(
		`
		INSERT INTO %s 
			(article, category_id, product_title, img_url, img_uuid, amount_in_stock, price, characteristics, packaging) 
		VALUES (
			$1, 
			(SELECT id FROM %s WHERE category_title = $2),
			$3, $4,
			$5, $6, $7, $8, $9
		) 
		RETURNING id
		`,
		postgres.ProductsTable,
		postgres.CategoriesTable,
	)
	insertRow := tx.QueryRow(
		queryInsertProduct,
		product.Article,
		product.CategoryTitle,
		product.ProductTitle,
		product.ImgUrl,
		product.ImgUUID,
		product.AmountInStock,
		product.Price,
		product.Characteristics,
		product.Packaging,
	)

	if err = insertRow.Scan(&productId); err != nil {
		_ = tx.Rollback()
		return 0, err
	}

	return productId, tx.Commit()
}

func (s *storage) Update(product Product) (int, error) {
	tx, _ := s.db.Begin()

	queryUpdateProduct := fmt.Sprintf(
		`
		UPDATE %s 
		SET article = $1,
			category_id = (SELECT id FROM %s WHERE category_title = $2),
			product_title = $3,
			img_url = $4,
			price = $5,
			characteristics = $6,
			packaging = $7
		WHERE id = $8
		`,
		postgres.ProductsTable,
		postgres.CategoriesTable,
	)

	_, err := tx.Exec(
		queryUpdateProduct,
		product.Article,
		product.CategoryTitle,
		product.ProductTitle,
		product.ImgUrl,
		product.Price,
		product.Characteristics,
		product.Packaging,
		product.Id,
	)
	if err != nil {
		_ = tx.Rollback()
		return 0, err
	}

	return product.Id, tx.Commit()
}

func (s *storage) UpdateImage(product Product) (int, error) {
	tx, _ := s.db.Begin()

	queryUpdateProduct := fmt.Sprintf(
		`UPDATE %s SET img_uuid = $1 WHERE id = $2`,
		postgres.ProductsTable,
	)

	_, err := tx.Exec(
		queryUpdateProduct,
		product.ImgUUID,
		product.Id,
	)
	if err != nil {
		_ = tx.Rollback()
		return 0, err
	}

	return product.Id, tx.Commit()
}

func (s *storage) Delete(productId int) error {
	tx, _ := s.db.Begin()

	queryDeleteProduct := fmt.Sprintf("DELETE FROM %s WHERE id=$1", postgres.ProductsTable)
	_, err := tx.Exec(queryDeleteProduct, productId)
	if err != nil {
		_ = tx.Rollback()
		return err
	}

	return tx.Commit()
}

func (s *storage) GetProductById(id int) (Product, error) {
	var product Product

	queryGetProduct, args, err := s.qb.Select(productsColumnsSelect...).
		From(postgres.ProductsTable).
		LeftJoin(postgres.CategoriesTable + " ON categories.id = products.category_id").
		Where(sq.Eq{"products.id": id}).ToSql()

	err = s.db.Get(&product, queryGetProduct, args...)
	if err != nil {
		return Product{}, err
	}

	return product, err
}

func (s *storage) GetFavourites(userId int) ([]Product, error) {
	var products []Product

	queryGetFavourites := fmt.Sprintf(
		`
		SELECT 
			products.id, 
			products.article, 
			products.product_title, 
			products.img_url, 
			products.img_uuid, 
			products.price, 
			products.packaging, 
			products.amount_in_stock, 
			products.created_at 
		FROM %s, %s 
		WHERE user_id = $1 AND %s.product_id = %s.id
		`,
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
