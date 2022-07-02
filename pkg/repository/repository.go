package repository

import "github.com/jmoiron/sqlx"

type Authorization interface {
}

type ShopItemCup interface {
}

type ShopList interface {
}

type Category interface {
}

type Repository struct {
	Authorization
	ShopItemCup
	ShopList
	Category
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{}
}
