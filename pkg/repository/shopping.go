package repository

import (
	server "allincecup-server"
	"fmt"
	"github.com/jmoiron/sqlx"
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

func (s *ShoppingPostgres) GetProductsInCart(userId int) ([]server.CartProduct, error) {
	var cartId int
	queryGetCartId := fmt.Sprintf("SELECT id FROM %s WHERE user_id=$1", cartsTable)
	if err := s.db.Get(&cartId, queryGetCartId, userId); err != nil {
		return nil, err
	}

	var productsInCart []server.CartProduct
	queryGetProductsInCart := fmt.Sprintf("SELECT * FROM %s WHERE cart_id=$1", cartsProductsTable)
	if err := s.db.Select(&productsInCart, queryGetProductsInCart, cartId); err != nil {
		return nil, err
	}

	return productsInCart, nil
}

func (s *ShoppingPostgres) AddToFavourites(userId, productId int) error {
	queryAddToFavourites := fmt.Sprintf("INSERT INTO %s (user_id, product_id) values ($1, $2)", favouritesTable)
	_, err := s.db.Exec(queryAddToFavourites, userId, productId)
	return err
}

func (s *ShoppingPostgres) GetFavourites(userId int) ([]server.Product, error) {
	var products []server.Product

	queryGetFavourites := fmt.Sprintf(
		"SELECT products.id, products.article, products.product_title, products.img_url, products.price, products.units_in_package, products.packages_in_box FROM %s, %s WHERE user_id=$1 AND %s.product_id=%s.id",
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
