package repository

import (
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	server "github.com/zh0vtyj/allincecup-server"
)

type Authorization interface {
	CreateUser(user server.User, role string) (int, int, error)
	GetUser(email string, password string) (server.User, error)
	NewSession(session server.Session) (*server.Session, error)
	GetSessionByRefresh(refresh string) (*server.Session, error)
	DeleteSessionByRefresh(refresh string) error
	DeleteSessionByUserId(id int) error
}

type Category interface {
	GetAll() ([]server.Category, error)
	Update(category server.Category) (int, error)
	Create(category server.Category) (int, error)
	Delete(id int, title string) error
}

type Products interface {
	Search(searchInput string) ([]server.Product, error)
	GetWithParams(params server.SearchParams, createdAt, search string) ([]server.Product, error)
	GetProductById(id int) (server.ProductInfoDescription, error)
	AddProduct(product server.Product, info []server.ProductInfo) (int, error)
	Update(product server.ProductInfoDescription) (int, error)
	Delete(productId int) error
}

type Shopping interface {
	AddToCart(userId int, info server.CartProduct) (float64, error)
	PriceValidation(productId, quantity int) (float64, error)
	GetProductsInCart(userId int) ([]server.CartProduct, error)
	DeleteFromCart(productId int) error
	AddToFavourites(userId, productId int) error
	GetFavourites(userId int) ([]server.Product, error)
	DeleteFromFavourites(userId, productId int) error
}

type Orders interface {
	New(order server.OrderFullInfo) (uuid.UUID, error)
	GetUserOrders(userId int, createdAt string) ([]server.OrderInfo, error)
	GetOrderById(orderId uuid.UUID) (server.OrderInfo, error)
	GetAdminOrders(status string, lastOrderCreatedAt string) ([]server.Order, error)
}

type Repository struct {
	Authorization
	Category
	Products
	Shopping
	Orders
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{
		Authorization: NewAuthPostgres(db),
		Products:      NewProductsPostgres(db),
		Category:      NewCategoryPostgres(db),
		Shopping:      NewShoppingPostgres(db),
		Orders:        NewOrdersPostgres(db),
	}
}
