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

func (s *ProductsService) AddProduct(product server.Product, info []server.ProductInfo) (int, error) {
	return s.repo.AddProduct(product, info)
}
