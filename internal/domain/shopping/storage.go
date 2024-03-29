package shopping

import (
	"fmt"
	sq "github.com/Masterminds/squirrel"
	"github.com/jmoiron/sqlx"
	"github.com/zh0vtyj/alliancecup-server/internal/domain/models"
	"github.com/zh0vtyj/alliancecup-server/pkg/client/postgres"
)

type Storage interface {
	AddToCart(userId int, info CartProduct) error
	PriceValidation(productId, quantity int) (float64, error)
	GetProductsInCart(userId int) ([]CartProduct, error)
	DeleteFromCart(productId, userId int) error
	GetFavourites(userId int) (products []models.Product, err error)
	AddToFavourites(userId, productId int) error
	DeleteFromFavourites(userId, productId int) error
}

type storage struct {
	db *sqlx.DB
	qb sq.StatementBuilderType
}

func NewShoppingPostgres(db *sqlx.DB, psql sq.StatementBuilderType) *storage {
	return &storage{
		db: db,
		qb: psql,
	}
}

func (s *storage) AddToCart(userId int, info CartProduct) error {
	queryAddToCart := fmt.Sprintf(
		`
		INSERT INTO %s 
			(cart_id, product_id, quantity) 
		values 
			((SELECT id FROM %s WHERE user_id=$1 LIMIT 1), $2, $3)
		`,
		postgres.CartsProductsTable,
		postgres.CartsTable,
	)
	_, err := s.db.Exec(queryAddToCart, userId, info.Id, info.Quantity)
	if err != nil {
		return err
	}

	return err
}

func (s *storage) PriceValidation(productId, quantity int) (float64, error) {
	var price float64
	query := fmt.Sprintf("SELECT price FROM %s WHERE id=$1", postgres.ProductsTable)
	if err := s.db.Get(&price, query, productId); err != nil {
		return 0, err
	}

	return price * float64(quantity), nil
}

func (s *storage) GetProductsInCart(userId int) ([]CartProduct, error) {
	var cartId int
	queryGetCartId := fmt.Sprintf("SELECT id FROM %s WHERE user_id=$1", postgres.CartsTable)
	if err := s.db.Get(&cartId, queryGetCartId, userId); err != nil {
		return nil, err
	}

	var productsInCart []CartProduct
	queryCartProducts, args, err := s.qb.Select(
		"carts_products.product_id",
		"products.article",
		"products.product_title",
		"products.img_url",
		"products.img_uuid",
		"products.amount_in_stock",
		"products.price",
		"products.packaging",
		"products.characteristics",
		"products.created_at",
		"carts_products.quantity",
		"products.price * carts_products.quantity as price_for_quantity",
	).
		From(postgres.CartsProductsTable).
		LeftJoin(postgres.ProductsTable + " ON carts_products.product_id = products.id").
		Where(sq.Eq{"carts_products.cart_id": cartId}).ToSql()
	if err != nil {
		return nil, fmt.Errorf("failed to build sql query to get products from cart due to: %v", err)
	}

	if err = s.db.Select(&productsInCart, queryCartProducts, args...); err != nil {
		return nil, fmt.Errorf("failed to get products from user cart due to: %v", err)
	}

	return productsInCart, nil
}

func (s *storage) DeleteFromCart(productId, userId int) error {
	queryDelete := fmt.Sprintf(
		`DELETE FROM %s WHERE cart_id = (SELECT id FROM %s WHERE user_id = $1) AND product_id = $2`,
		postgres.CartsProductsTable,
		postgres.CartsTable,
	)
	_, err := s.db.Exec(queryDelete, userId, productId)
	return err
}

func (s *storage) AddToFavourites(userId, productId int) error {
	queryAddToFavourites := fmt.Sprintf(
		"INSERT INTO %s (user_id, product_id) values ($1, $2)",
		postgres.FavouritesTable,
	)

	_, err := s.db.Exec(queryAddToFavourites, userId, productId)

	return err
}

func (s *storage) GetFavourites(userId int) (products []models.Product, err error) {
	queryGetFavourites := fmt.Sprintf(
		`
		SELECT 
			products.id, 
			products.article, 
			products.product_title, 
			products.img_url, 
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

func (s *storage) DeleteFromFavourites(userId, productId int) error {
	queryDeleteProduct := fmt.Sprintf(
		"DELETE FROM %s WHERE user_id = $1 AND product_id = $2",
		postgres.FavouritesTable,
	)

	_, err := s.db.Exec(queryDeleteProduct, userId, productId)
	if err != nil {
		return err
	}

	return nil
}
