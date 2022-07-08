package service

import (
	server "allincecup-server"
	"allincecup-server/internal/domain"
	"allincecup-server/pkg/repository"
)

type Authorization interface {
	CreateUser(user server.User) (int, error)
	GenerateTokens(email string, password string) (string, string, error)
	ParseToken(token string) (int, int, error)
	ParseRefreshToken(refreshToken string) error
	RefreshAccessToken(refreshToken string) (string, error)
	CreateNewSession(session *domain.Session) (*domain.Session, error)
}

type Category interface {
	Create(title string) (int, error)
}

type Products interface {
	AddProduct(product server.Product) (int, error)
}

type Service struct {
	Authorization
	Category
	Products
}

func NewService(repos *repository.Repository) *Service {
	return &Service{
		Authorization: NewAuthService(repos.Authorization),
		Products:      NewProductsService(repos.Products),
		Category:      NewCategoryService(repos.Category),
	}
}
