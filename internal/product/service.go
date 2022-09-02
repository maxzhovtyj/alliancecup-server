package product

import (
	server "github.com/zh0vtyj/allincecup-server/internal/shopping"
)

type Service interface {
	Search(searchInput string) ([]Product, error)
	GetWithParams(params server.SearchParams) ([]Product, error)
	GetProductById(id int) (Description, error)
	AddProduct(product Product, info []Info) (int, error)
	GetFavourites(userId int) ([]Product, error)
	Update(product Description) (int, error)
	Delete(productId int) error
}

type service struct {
	repo Storage
}

func NewProductsService(repo Storage) Service {
	return &service{repo: repo}
}

func (s *service) Search(searchInput string) ([]Product, error) {
	searchInput = "%" + searchInput + "%"
	products, err := s.repo.Search(searchInput)
	if err != nil {
		return nil, err
	}
	return products, err
}

func (s *service) GetWithParams(params server.SearchParams) ([]Product, error) {
	return s.repo.GetWithParams(params)
}

func (s *service) AddProduct(product Product, info []Info) (int, error) {
	return s.repo.AddProduct(product, info)
}

func (s *service) GetFavourites(userId int) ([]Product, error) {
	return s.repo.GetFavourites(userId)
}

func (s *service) Update(product Description) (int, error) {
	return s.repo.Update(product)
}

func (s *service) Delete(productId int) error {
	return s.repo.Delete(productId)
}

func (s *service) GetProductById(id int) (Description, error) {
	return s.repo.GetProductById(id)
}
