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

func (s *ShoppingPostgres) AddToCart(userId int, info server.ProductOrder) (float64, error) {
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

func (s *ShoppingPostgres) GetProductsInCart(userId int) ([]server.ProductOrder, error) {
	var cartId int
	queryGetCartId := fmt.Sprintf("SELECT id FROM %s WHERE user_id=$1", cartsTable)
	if err := s.db.Get(&cartId, queryGetCartId, userId); err != nil {
		return nil, err
	}

	var productsInCart []server.ProductOrder
	queryGetProductsInCart := fmt.Sprintf("SELECT * FROM %s WHERE cart_id=$1", cartsProductsTable)
	if err := s.db.Select(&productsInCart, queryGetProductsInCart, cartId); err != nil {
		return nil, err
	}

	return productsInCart, nil
}
