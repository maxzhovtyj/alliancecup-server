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

func (c *CategoryPostgres) Update(category server.Category) (int, error) {
	var id int

	queryUpdate := fmt.Sprintf("UPDATE %s SET category_title=$1, img_url=$2 WHERE id=$3 RETURNING id", categoriesTable)
	row := c.db.QueryRow(queryUpdate, category.CategoryTitle, category.ImgUrl, category.Id)
	if err := row.Scan(&id); err != nil {
		return 0, err
	}

	return id, nil
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

func (c *CategoryPostgres) Delete(id int, title string) error {
	queryDeleteCategory := fmt.Sprintf("DELETE FROM %s WHERE id=$1 OR category_title=$2", categoriesTable)
	_, err := c.db.Exec(queryDeleteCategory, id, title)
	return err
}
