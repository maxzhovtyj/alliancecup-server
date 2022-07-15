package repository

import (
	server "allincecup-server"
	"allincecup-server/internal/domain"
	"github.com/jmoiron/sqlx"
)

type Authorization interface {
	CreateUser(user server.User, role string) (int, int, error)
	GetUser(email string, password string) (server.User, error)
	NewSession(session domain.Session) (*domain.Session, error)
	GetSessionByRefresh(refresh string) (*domain.Session, error)
	DeleteSessionByRefresh(refresh string) error
}

type Category interface {
	GetAll() ([]server.Category, error)
	Create(category server.Category) (int, error)
}

type Products interface {
	AddProduct(product server.Product, info []server.ProductInfo) (int, error)
}

type Repository struct {
	Authorization
	Category
	Products
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{
		Authorization: NewAuthPostgres(db),
		Products:      NewProductsPostgres(db),
		Category:      NewCategoryPostgres(db),
	}
}
