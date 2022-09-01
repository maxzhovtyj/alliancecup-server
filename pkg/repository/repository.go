package repository

import (
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/zh0vtyj/allincecup-server/pkg/models"
)

type Authorization interface {
	CreateUser(user models.User, role string) (int, int, error)
	GetUser(email string, password string) (models.User, error)
	NewSession(session models.Session) (*models.Session, error)
	GetSessionByRefresh(refresh string) (*models.Session, error)
	DeleteSessionByRefresh(refresh string) error
	DeleteSessionByUserId(id int) error
	UpdateRefreshToken(userId int, newRefreshToken string) error
	GetUserPasswordHash(userId int) (string, error)
	UpdatePassword(userId int, newPassword string) error
}

type Category interface {
	GetAll() ([]models.Category, error)
	GetFiltration(fkName string, id int) ([]models.CategoryFiltration, error)
	Update(category models.Category) (int, error)
	Create(category models.Category) (int, error)
	Delete(id int) error
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
	PriceValidation(productId, quantity int) (float64, error)
	GetProductsInCart(userId int) ([]models.CartProductFullInfo, error)
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
	GetDeliveryTypes() ([]models.DeliveryType, error)
	GetPaymentTypes() ([]models.PaymentType, error)
	ChangeOrderStatus(orderId uuid.UUID, toStatus string) error
}

type Supply interface {
	New(supply models.SupplyDTO) error
	GetAll(createdAt string) ([]models.SupplyInfoDTO, error)
	UpdateProductsAmount(products []models.ProductSupplyDTO, operation string) error
	DeleteAndGetProducts(id int) ([]models.ProductSupplyDTO, error)
}

type Repository struct {
	Authorization
	Category
	Products
	Shopping
	Orders
	Supply
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{
		Authorization: NewAuthPostgres(db),
		Products:      NewProductsPostgres(db),
		Category:      NewCategoryPostgres(db),
		Shopping:      NewShoppingPostgres(db),
		Orders:        NewOrdersPostgres(db),
		Supply:        NewSupplyPostgres(db),
	}
}
