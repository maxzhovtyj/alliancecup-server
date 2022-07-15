package repository

import "github.com/jmoiron/sqlx"

type OrdersPostgres struct {
	db *sqlx.DB
}

func NewOrdersPostgres(db *sqlx.DB) *OrdersPostgres {
	return &OrdersPostgres{db: db}
}

func (o *OrdersPostgres) GetUserOrders() {

}

func (o *OrdersPostgres) NewUserOrder() {

}

func (o *OrdersPostgres) GetAllOrders() {

}
