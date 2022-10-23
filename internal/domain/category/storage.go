package category

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/zh0vtyj/allincecup-server/pkg/client/postgres"
)

type Storage interface {
	GetAll() ([]Category, error)
	GetFiltration(fkName string, id int) ([]Filtration, error)
	Update(category Category) (int, error)
	Create(category Category) (int, error)
	Delete(id int) error
	AddFiltration(filtration Filtration) (int, error)
}

type storage struct {
	db *sqlx.DB
}

func NewCategoryPostgres(db *sqlx.DB) *storage {
	return &storage{db: db}
}

func (c *storage) GetAll() ([]Category, error) {
	var categories []Category

	queryGetCategories := fmt.Sprintf(
		`
		SELECT id,
			   category_title,
			   img_url,
			   category_description 
		FROM %s
		`,
		postgres.CategoriesTable,
	)
	err := c.db.Select(&categories, queryGetCategories)

	return categories, err
}

func (c *storage) GetFiltration(fkName string, id int) ([]Filtration, error) {
	var filtration []Filtration

	//fkName can be either category_id or filtration_list_id
	queryGetFiltration := fmt.Sprintf("SELECT * FROM %s WHERE %s=$1", postgres.CategoriesFiltrationTable, fkName)
	err := c.db.Select(&filtration, queryGetFiltration, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get category filtration from db due to: %v", err)
	}

	return filtration, nil
}

func (c *storage) Update(category Category) (int, error) {
	var id int

	queryUpdate := fmt.Sprintf("UPDATE %s SET category_title=$1, img_url=$2 WHERE id=$3 RETURNING id", postgres.CategoriesTable)
	row := c.db.QueryRow(queryUpdate, category.CategoryTitle, category.ImgUrl, category.Id)
	if err := row.Scan(&id); err != nil {
		return 0, err
	}

	return id, nil
}

func (c *storage) Create(category Category) (int, error) {
	var id int
	query := fmt.Sprintf(
		`
		INSERT INTO %s 
			(category_title, img_url, img_uuid) 
		values 
			($1, $2, $3) 
		RETURNING id
		`,
		postgres.CategoriesTable,
	)
	row := c.db.QueryRow(query, category.CategoryTitle, category.ImgUrl, category.ImgUUID)

	if err := row.Scan(&id); err != nil {
		return 0, err
	}

	return id, nil
}

func (c *storage) Delete(id int) error {
	queryDeleteCategory := fmt.Sprintf("DELETE FROM %s WHERE id=$1", postgres.CategoriesTable)
	_, err := c.db.Exec(queryDeleteCategory, id)
	return err
}

func (c *storage) AddFiltration(filtration Filtration) (int, error) {
	var filtrationId int

	queryAddFiltration := fmt.Sprintf(
		"INSERT INTO %s (info_description, filtration_title, filtration_description, img_url, category_id, filtration_list_id) VALUES ($1, $2, $3, $4, $5, $6) RETURNING id",
		postgres.CategoriesFiltrationTable,
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
