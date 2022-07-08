package repository

import (
	"fmt"
	"github.com/jmoiron/sqlx"
)

type CategoryPostgres struct {
	db *sqlx.DB
}

func NewCategoryPostgres(db *sqlx.DB) *CategoryPostgres {
	return &CategoryPostgres{db: db}
}

func (c *CategoryPostgres) Create(title string) (int, error) {
	var id int
	query := fmt.Sprintf("INSERT INTO %s (title) values ($1) RETURNING id", categoriesTable)
	row := c.db.QueryRow(query, title)

	if err := row.Scan(&id); err != nil {
		return 0, err
	}

	return id, nil
}

func (c *CategoryPostgres) GetAll() ([]string, error) {
	var categories []string
	query := fmt.Sprintf("SELECT * FROM %s", categoriesTable)
	err := c.db.Get(categories, query)
	if err != nil {
		return nil, err
	}

	return categories, nil
}
