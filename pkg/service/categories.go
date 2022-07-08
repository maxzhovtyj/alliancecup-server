package service

import "allincecup-server/pkg/repository"

type CategoryService struct {
	repo repository.Category
}

func NewCategoryService(repo repository.Category) *CategoryService {
	return &CategoryService{repo: repo}
}

func (c *CategoryService) Create(title string) (int, error) {
	id, err := c.repo.Create(title)
	if err != nil {
		return 0, err
	}
	return id, err
}
