package repository

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	server "github.com/zh0vtyj/allincecup-server/pkg/models"
)

type CategoryPostgres struct {
	db *sqlx.DB
}

func NewCategoryPostgres(db *sqlx.DB) *CategoryPostgres {
	return &CategoryPostgres{db: db}
}

func (c *CategoryPostgres) GetAll() ([]server.Category, error) {
	var categories []server.Category

	queryGetCategories := fmt.Sprintf("SELECT * FROM %s", categoriesTable)
	err := c.db.Select(&categories, queryGetCategories)

	return categories, err
}

func (c *CategoryPostgres) GetFiltration(fkName string, id int) ([]server.CategoryFiltration, error) {
	var filtration []server.CategoryFiltration

	//fkName can be either category_id or filtration_list_id
	queryGetFiltration := fmt.Sprintf("SELECT * FROM %s WHERE %s=$1", categoriesFiltrationTable, fkName)
	err := c.db.Select(&filtration, queryGetFiltration, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get category filtration from db due to: %v", err)
	}

	return filtration, nil
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

func (c *CategoryPostgres) AddFiltration(filtration server.CategoryFiltration) (int, error) {
	var filtrationId int

	queryAddFiltration := fmt.Sprintf(
		"INSERT INTO %s (info_description, filtration_title, filtration_description, img_url, category_id, filtration_list_id) VALUES ($1, $2, $3, $4, $5, $6) RETURNING id",
		categoriesFiltrationTable,
	)
	row := c.db.QueryRow(
		queryAddFiltration,
		filtration.InfoDescription,
		filtration.FiltrationTitle,
		filtration.FiltrationDescription,
		filtration.ImgUrl,
		filtration.CategoryId,
		filtration.FiltrationListId,
	)

	if err := row.Scan(&filtrationId); err != nil {
		return 0, fmt.Errorf("failed to execute query to a db due to: %v", err)
	}

	return filtrationId, nil
}
