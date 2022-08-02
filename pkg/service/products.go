package service

import (
	server "github.com/zh0vtyj/allincecup-server"
	"github.com/zh0vtyj/allincecup-server/pkg/repository"
)

type ProductsService struct {
	repo repository.Products
}

func NewProductsService(repo repository.Products) *ProductsService {
	return &ProductsService{repo: repo}
}

func (s *ProductsService) Search(searchInput string) ([]server.Product, error) {
	searchInput = "%" + searchInput + "%"
	products, err := s.repo.Search(searchInput)
	if err != nil {
		return nil, err
	}
	return products, err
}

func (s *ProductsService) GetWithParams(params server.SearchParams, lastProductCreatedAt, search string) ([]server.Product, error) {
	return s.repo.GetWithParams(params, lastProductCreatedAt, search)
}

func (s *ProductsService) AddProduct(product server.Product, info []server.ProductInfo) (int, error) {
	return s.repo.AddProduct(product, info)
}

func (s *ProductsService) Update(product server.ProductInfoDescription) (int, error) {
	return s.repo.Update(product)
}

func (s *ProductsService) Delete(productId int) error {
	return s.repo.Delete(productId)
}

func (s *ProductsService) GetProductById(id int) (server.ProductInfoDescription, error) {
	return s.repo.GetProductById(id)
}
