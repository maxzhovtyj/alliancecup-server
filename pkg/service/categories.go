package service

import (
	server "github.com/zh0vtyj/allincecup-server/pkg/models"
	"github.com/zh0vtyj/allincecup-server/pkg/repository"
)

type CategoryService struct {
	repo repository.Category
}

func NewCategoryService(repo repository.Category) *CategoryService {
	return &CategoryService{repo: repo}
}

func (c *CategoryService) GetAll() ([]server.Category, error) {
	return c.repo.GetAll()
}

func (c *CategoryService) Update(category server.Category) (int, error) {
	return c.repo.Update(category)
}

func (c *CategoryService) Create(category server.Category) (int, error) {
	id, err := c.repo.Create(category)
	if err != nil {
		return 0, err
	}
	return id, err
}

func (c *CategoryService) Delete(id int) error {
	return c.repo.Delete(id)
}

func (c *CategoryService) AddFiltration(filtration server.CategoryFiltration) (int, error) {
	return c.repo.AddFiltration(filtration)
}

func (c *CategoryService) GetFiltration(fkName string, id int) ([]server.CategoryFiltration, error) {
	return c.repo.GetFiltration(fkName, id)
}
