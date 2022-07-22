package service

import (
	server "allincecup-server"
	"allincecup-server/pkg/repository"
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

func (c *CategoryService) Delete(id int, title string) error {
	return c.repo.Delete(id, title)
}
