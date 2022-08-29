package service

import (
	"github.com/google/uuid"
	"github.com/zh0vtyj/allincecup-server/pkg/models"
	"github.com/zh0vtyj/allincecup-server/pkg/repository"
)

type Authorization interface {
	CreateUser(user models.User) (int, int, error)
	CreateModerator(user models.User) (int, int, error)
	GenerateTokens(email string, password string) (string, string, error)
	ParseToken(token string) (int, int, error)
	ParseRefreshToken(refreshToken string) error
	RefreshTokens(refreshToken, clientIp, userAgent string) (string, string, int, int, error)
	CreateNewSession(session *models.Session) (*models.Session, error)
	Logout(id int) error
	ChangePassword(userId int, oldPassword, newPassword string) error
}

type Category interface {
	GetAll() ([]models.Category, error)
	GetFiltration(fkName string, id int) ([]models.CategoryFiltration, error)
	Update(category models.Category) (int, error)
	Create(category models.Category) (int, error)
	Delete(id int, title string) error
	AddFiltration(filtration models.CategoryFiltration) (int, error)
}

type Products interface {
	Search(searchInput string) ([]models.Product, error)
	GetWithParams(params models.SearchParams) ([]models.Product, error)
	GetProductById(id int) (models.ProductInfoDescription, error)
	AddProduct(product models.Product, info []models.ProductInfo) (int, error)
	Update(product models.ProductInfoDescription) (int, error)
	Delete(productId int) error
}

type Shopping interface {
	AddToCart(userId int, info models.CartProduct) (float64, error)
	GetProductsInCart(userId int) ([]models.CartProductFullInfo, float64, error)
	DeleteFromCart(productId int) error
	AddToFavourites(userId, productId int) error
	GetFavourites(userId int) ([]models.Product, error)
	DeleteFromFavourites(userId, productId int) error
}

type Orders interface {
	New(order models.OrderFullInfo) (uuid.UUID, error)
	GetUserOrders(userId int, createdAt string) ([]models.OrderInfo, error)
	GetOrderById(orderId uuid.UUID) (models.OrderInfo, error)
	GetAdminOrders(status string, lastOrderCreatedAt string) ([]models.Order, error)
	DeliveryPaymentTypes() (models.DeliveryPaymentTypes, error)
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
