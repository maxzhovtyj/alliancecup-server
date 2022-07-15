package repository

import (
	server "allincecup-server"
	"fmt"
	"github.com/jmoiron/sqlx"
)

type CategoryPostgres struct {
	db *sqlx.DB
}

func NewCategoryPostgres(db *sqlx.DB) *CategoryPostgres {
	return &CategoryPostgres{db: db}
}

func (c *CategoryPostgres) GetAll() ([]server.Category, error) {
	var categories []server.Category

	query := fmt.Sprintf("SELECT * FROM %s", categoriesTable)
	err := c.db.Select(&categories, query)

	return categories, err
}

func (c *CategoryPostgres) Create(category server.Category) (int, error) {
	var id int
	query := fmt.Sprintf("INSERT INTO %s (category_title, img_url) values ($1, $2) RETURNING id", categoriesTable)
	row := c.db.QueryRow(query, category.CategoryTitle, category.ImgUrl)

	if err := row.Scan(&id); err != nil {
		return 0, err
	}

	return id, nil
}
