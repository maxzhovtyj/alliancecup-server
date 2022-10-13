package shopping

import (
	"fmt"
	sq "github.com/Masterminds/squirrel"
	"github.com/jmoiron/sqlx"
	"github.com/zh0vtyj/allincecup-server/pkg/client/postgres"
)

type Storage interface {
	AddToCart(userId int, info CartProduct) (float64, error)
	PriceValidation(productId, quantity int) (float64, error)
	GetProductsInCart(userId int) ([]CartProductFullInfo, error)
	DeleteFromCart(productId int) error
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

func (s *storage) AddToCart(userId int, info CartProduct) (float64, error) {
	var userCartId int
	queryGetCartId := fmt.Sprintf("SELECT id FROM %s WHERE user_id=$1 LIMIT 1", postgres.CartsTable)
	err := s.db.Get(&userCartId, queryGetCartId, userId)
	if err != nil {
		return 0, err
	}

	queryAddToCart := fmt.Sprintf("INSERT INTO %s (cart_id, product_id, quantity, price_for_quantity) values ($1, $2, $3, $4)", postgres.CartsProductsTable)
	_, err = s.db.Exec(queryAddToCart, userCartId, info.ProductId, info.Quantity, info.PriceForQuantity)
	if err != nil {
		return 0, err
	}

	return info.PriceForQuantity, err
}

func (s *storage) PriceValidation(productId, quantity int) (float64, error) {
	var price float64
	query := fmt.Sprintf("SELECT price FROM %s WHERE id=$1", postgres.ProductsTable)
	if err := s.db.Get(&price, query, productId); err != nil {
		return 0, err
	}

	return price * float64(quantity), nil
}

func (s *storage) GetProductsInCart(userId int) ([]CartProductFullInfo, error) {
	var cartId int
	queryGetCartId := fmt.Sprintf("SELECT id FROM %s WHERE user_id=$1", postgres.CartsTable)
	if err := s.db.Get(&cartId, queryGetCartId, userId); err != nil {
		return nil, err
	}

	var productsInCart []CartProductFullInfo
	queryCartProducts, args, err := s.qb.Select(
		"carts_products.product_id",
		"products.article",
		"products.product_title",
		"products.img_url",
		"products.amount_in_stock",
		"products.price",
		"products.units_in_package",
		"products.packages_in_box",
		"products.created_at",
		"carts_products.quantity",
		"carts_products.price_for_quantity",
	).
		From(postgres.CartsProductsTable).
		LeftJoin(postgres.ProductsTable + " ON carts_products.product_id=products.id").
		Where(sq.Eq{"carts_products.cart_id": cartId}).ToSql()
	if err != nil {
		return nil, fmt.Errorf("failed to build sql query to get products from cart due to: %v", err)
	}

	if err = s.db.Select(&productsInCart, queryCartProducts, args...); err != nil {
		return nil, fmt.Errorf("failed to get products from user cart due to: %v", err)
	}

	return productsInCart, nil
}

func (s *storage) DeleteFromCart(productId int) error {
	queryDelete := fmt.Sprintf("DELETE FROM %s WHERE product_id=$1", postgres.CartsProductsTable)
	_, err := s.db.Exec(queryDelete, productId)
	return err
}

func (s *storage) AddToFavourites(userId, productId int) error {
	queryAddToFavourites := fmt.Sprintf("INSERT INTO %s (user_id, product_id) values ($1, $2)", postgres.FavouritesTable)
	_, err := s.db.Exec(queryAddToFavourites, userId, productId)
	return err
}

func (s *storage) DeleteFromFavourites(userId, productId int) error {
	queryDeleteProduct := fmt.Sprintf("DELETE FROM %s WHERE user_id=$1 AND product_id=$2", postgres.FavouritesTable)
	_, err := s.db.Exec(queryDeleteProduct, userId, productId)
	if err != nil {
		return err
	}
	return nil
}
