package service

import (
	server "allincecup-server"
	"allincecup-server/pkg/repository"
)

type ProductsService struct {
	repo repository.Products
}

func NewProductsService(repo repository.Products) *ProductsService {
	return &ProductsService{repo: repo}
}

func (s *ProductsService) AddProduct(product server.Product) (int, error) {
	id, err := s.repo.AddProduct(product)

	if err != nil {
		return 0, err
	}
	return id, nil
}
