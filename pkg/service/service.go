package service

import (
	"github.com/google/uuid"
	server "github.com/zh0vtyj/allincecup-server"
	"github.com/zh0vtyj/allincecup-server/pkg/repository"
)

type Authorization interface {
	CreateUser(user server.User) (int, int, error)
	CreateModerator(user server.User) (int, int, error)
	GenerateTokens(email string, password string) (string, string, error)
	ParseToken(token string) (int, int, error)
	ParseRefreshToken(refreshToken string) error
	RefreshTokens(refreshToken, clientIp, userAgent string) (string, string, int, int, error)
	CreateNewSession(session *server.Session) (*server.Session, error)
	Logout(id int) error
	ChangePassword(userId int, oldPassword, newPassword string) error
}

type Category interface {
	GetAll() ([]server.Category, error)
	GetFiltration(fkName string, id int) ([]server.CategoryFiltration, error)
	Update(category server.Category) (int, error)
	Create(category server.Category) (int, error)
	Delete(id int, title string) error
	AddFiltration(filtration server.CategoryFiltration) (int, error)
}

type Products interface {
	Search(searchInput string) ([]server.Product, error)
	GetWithParams(params server.SearchParams, lastProductCreatedAt, search string) ([]server.Product, error)
	GetProductById(id int) (server.ProductInfoDescription, error)
	AddProduct(product server.Product, info []server.ProductInfo) (int, error)
	Update(product server.ProductInfoDescription) (int, error)
	Delete(productId int) error
}

type Shopping interface {
	AddToCart(userId int, info server.CartProduct) (float64, error)
	GetProductsInCart(userId int) ([]server.CartProductFullInfo, float64, error)
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
	DeliveryPaymentTypes() (server.DeliveryPaymentTypes, error)
	ProcessedOrder(orderId uuid.UUID, toStatus string) error
}

type Service struct {
	Authorization
	Category
	Products
	Shopping
	Orders
}

func NewService(repos *repository.Repository) *Service {
	return &Service{
		Authorization: NewAuthService(repos.Authorization),
		Products:      NewProductsService(repos.Products),
		Category:      NewCategoryService(repos.Category),
		Shopping:      NewShoppingService(repos.Shopping),
		Orders:        NewOrdersService(repos.Orders, repos.Products),
	}
}
