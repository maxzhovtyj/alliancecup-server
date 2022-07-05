package service

import "allincecup-server/pkg/repository"

type CategoryService struct {
	repo repository.Categories
}

func NewCategoriesService(repo repository.Categories) *CategoryService {
	return &CategoryService{repo: repo}
}
