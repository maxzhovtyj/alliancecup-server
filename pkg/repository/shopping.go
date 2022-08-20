package repository

import (
	"fmt"
	sq "github.com/Masterminds/squirrel"
	"github.com/jmoiron/sqlx"
	server "github.com/zh0vtyj/allincecup-server"
)

type ShoppingPostgres struct {
	db *sqlx.DB
}

func NewShoppingPostgres(db *sqlx.DB) *ShoppingPostgres {
	return &ShoppingPostgres{db: db}
}

func (s *ShoppingPostgres) AddToCart(userId int, info server.CartProduct) (float64, error) {
	var userCartId int
	queryGetCartId := fmt.Sprintf("SELECT id FROM %s WHERE user_id=$1 LIMIT 1", cartsTable)
	err := s.db.Get(&userCartId, queryGetCartId, userId)
	if err != nil {
		return 0, err
	}

	queryAddToCart := fmt.Sprintf("INSERT INTO %s (cart_id, product_id, quantity, price_for_quantity) values ($1, $2, $3, $4)", cartsProductsTable)
	_, err = s.db.Exec(queryAddToCart, userCartId, info.ProductId, info.Quantity, info.PriceForQuantity)
	if err != nil {
		return 0, err
	}

	return info.PriceForQuantity, err
}

func (s *ShoppingPostgres) PriceValidation(productId, quantity int) (float64, error) {
	var price float64
	query := fmt.Sprintf("SELECT price FROM %s WHERE id=$1", productsTable)
	if err := s.db.Get(&price, query, productId); err != nil {
		return 0, err
	}

	return price * float64(quantity), nil
}

func (s *ShoppingPostgres) GetProductsInCart(userId int) ([]server.CartProductFullInfo, error) {
	var cartId int
	queryGetCartId := fmt.Sprintf("SELECT id FROM %s WHERE user_id=$1", cartsTable)
	if err := s.db.Get(&cartId, queryGetCartId, userId); err != nil {
		return nil, err
	}

	psql := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)

	var productsInCart []server.CartProductFullInfo
	queryCartProducts, args, err := psql.Select(
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
		From(cartsProductsTable).
		LeftJoin(productsTable + " ON carts_products.product_id=products.id").
		Where(sq.Eq{"carts_products.cart_id": cartId}).ToSql()
	if err != nil {
		return nil, fmt.Errorf("failed to build sql query to get products from cart due to: %v", err)
	}

	if err = s.db.Select(&productsInCart, queryCartProducts, args...); err != nil {
		return nil, fmt.Errorf("failed to get products from user cart due to: %v", err)
	}

	return productsInCart, nil
}

func (s *ShoppingPostgres) DeleteFromCart(productId int) error {
	queryDelete := fmt.Sprintf("DELETE FROM %s WHERE product_id=$1", cartsProductsTable)
	_, err := s.db.Exec(queryDelete, productId)
	return err
}

func (s *ShoppingPostgres) AddToFavourites(userId, productId int) error {
	queryAddToFavourites := fmt.Sprintf("INSERT INTO %s (user_id, product_id) values ($1, $2)", favouritesTable)
	_, err := s.db.Exec(queryAddToFavourites, userId, productId)
	return err
}

func (s *ShoppingPostgres) GetFavourites(userId int) ([]server.Product, error) {
	var products []server.Product

	queryGetFavourites := fmt.Sprintf(
		"SELECT products.id, products.article, products.product_title, products.img_url, products.price, products.units_in_package, products.packages_in_box, products.created_at FROM %s, %s WHERE user_id=$1 AND %s.product_id=%s.id",
		favouritesTable,
		productsTable,
		favouritesTable,
		productsTable,
	)

	if err := s.db.Select(&products, queryGetFavourites, userId); err != nil {
		return nil, err
	}

	return products, nil
}

func (s *ShoppingPostgres) DeleteFromFavourites(userId, productId int) error {
	queryDeleteProduct := fmt.Sprintf("DELETE FROM %s WHERE user_id=$1 AND product_id=$2", favouritesTable)
	_, err := s.db.Exec(queryDeleteProduct, userId, productId)
	if err != nil {
		return err
	}
	return nil
}
