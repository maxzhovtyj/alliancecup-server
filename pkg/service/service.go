package service

import (
	server "allincecup-server"
	"allincecup-server/pkg/repository"
)

type Authorization interface {
	CreateUser(user server.User) (int, error)
	GenerateTokens(email string, password string) (string, string, error)
	ParseToken(token string) (int, int, error)
}

type ShopItemCup interface {
}

type ShopList interface {
}

type Categories interface {
}

type Products interface {
	AddProduct(product server.Product) (int, error)
}

type Service struct {
	Authorization
	ShopItemCup
	ShopList
	Categories
	Products
}

func NewService(repos *repository.Repository) *Service {
	return &Service{
		Authorization: NewAuthService(repos.Authorization),
		Products:      NewProductsService(repos.Products),
		Categories:    NewCategoriesService(repos.Categories),
	}
}
