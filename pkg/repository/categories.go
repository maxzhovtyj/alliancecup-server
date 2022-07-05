package repository

import (
	"github.com/jmoiron/sqlx"
)

type CategoriesPostgres struct {
	db *sqlx.DB
}

func NewCategoriesPostgres(db *sqlx.DB) *CategoriesPostgres {
	return &CategoriesPostgres{db: db}
}
