package service

import (
	server "allincecup-server"
	"allincecup-server/internal/domain"
	"allincecup-server/pkg/repository"
)

type Authorization interface {
	CreateUser(user server.User) (int, int, error)
	CreateModerator(user server.User) (int, int, error)
	GenerateTokens(email string, password string) (string, string, error)
	ParseToken(token string) (int, int, error)
	ParseRefreshToken(refreshToken string) error
	RefreshAccessToken(refreshToken string) (string, error)
	CreateNewSession(session *domain.Session) (*domain.Session, error)
}

type Category interface {
	GetAll() ([]server.Category, error)
	Update(category server.Category) (int, error)
	Create(category server.Category) (int, error)
	Delete(id int, title string) error
}

type Products interface {
	GetProductById(id int) (server.ProductInfoDescription, error)
	AddProduct(product server.Product, info []server.ProductInfo) (int, error)
	Update(product server.ProductInfoDescription) (int, error)
	Delete(productId int) error
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
