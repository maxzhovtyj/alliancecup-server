package handler

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/swaggo/files"       // swagger embed files
	"github.com/swaggo/gin-swagger" // gin-swagger middleware
	_ "github.com/zh0vtyj/allincecup-server/docs"
	"github.com/zh0vtyj/allincecup-server/internal/domain/service"
	"github.com/zh0vtyj/allincecup-server/pkg/logging"
	"net/http"
)

const (
	StatusInProgress   = "IN_PROGRESS"
	StatusProcessed    = "PROCESSED"
	StatusCompleted    = "COMPLETED"
	refreshTokenCookie = "refresh_token"
	domain             = "localhost"
)

const (
	authUrl              = "/auth"
	apiUrl               = "/api"
	adminUrl             = "/admin"
	clientUrl            = "/client"
	signInUrl            = "/sign-in"
	signUpUrl            = "/sign-up"
	logoutUrl            = "/logout"
	changePasswordUrl    = "/change-password"
	refreshUrl           = "/refresh"
	categoriesUrl        = "/categories"
	categoryUrl          = "/category"
	productsUrl          = "/products"
	productUrl           = "/product"
	reviewsUrl           = "/reviews"
	reviewUrl            = "/review"
	cartUrl              = "/cart"
	favouritesUrl        = "/favourites"
	filtrationUrl        = "/filtration"
	ordersUrl            = "/orders"
	orderUrl             = "/order"
	userOrdersUrl        = "/user-orders"
	orderInfoTypesUrl    = "/order-info-types"
	moderatorUrl         = "/moderator"
	superAdminUrl        = "/super"
	supplyUrl            = "/supply"
	supplyProductsUrl    = "/supply-products"
	inventoryUrl         = "/inventory"
	inventoriesUrl       = "/inventories"
	inventoryProductsUrl = "/inventory-products"
	invoiceUrl           = "/invoice"
)

type Handler struct {
	services *service.Service
	logger   *logging.Logger
}

func NewHandler(services *service.Service, logger *logging.Logger) *Handler {
	return &Handler{
		services: services,
		logger:   logger,
	}
}

func (h *Handler) InitRoutes() *gin.Engine {
	router := gin.New()

	c := cors.New(cors.Config{
		AllowOrigins: []string{
			"http://localhost:3000",
		},
		AllowMethods: []string{
			http.MethodGet, http.MethodDelete, http.MethodPost, http.MethodPut,
		},
		AllowHeaders: []string{
			"Authorization",
			"Content-Type",
			"User-Agent",
		},
		AllowCredentials: true,
		ExposeHeaders:    []string{},
	})

	router.Use(c)

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	auth := router.Group(authUrl)
	{
		auth.POST(signUpUrl, h.signUp)
		auth.POST(signInUrl, h.signIn)
		auth.POST(refreshUrl, h.refresh)
	}

	api := router.Group(apiUrl, h.userIdentity)
	{
		api.GET(categoriesUrl, h.getCategories)
		api.GET(filtrationUrl, h.getFiltration)
		api.GET(productsUrl, h.getProducts)
		api.GET(productUrl, h.getProductById)
		api.POST(orderUrl, h.newOrder)
		api.GET(orderInfoTypesUrl, h.deliveryPaymentTypes)
		api.POST(reviewUrl, h.addReview)
		api.GET(reviewsUrl, h.getReviews)

		api.POST("forgot-password", h.forgotPassword)

		api.GET(invoiceUrl, h.getOrderInvoice)

		admin := api.Group(adminUrl, h.userHasPermission)
		{
			admin.POST(moderatorUrl, h.createModerator)

			admin.POST(productUrl, h.addProduct)
			admin.PUT(productUrl, h.updateProduct)
			admin.DELETE(productUrl, h.deleteProduct)
			//admin.PUT("/product-amount") // TODO

			admin.POST(categoryUrl, h.addCategory)
			admin.PUT(categoryUrl, h.updateCategory)
			admin.DELETE(categoryUrl, h.deleteCategory)

			admin.GET(ordersUrl, h.adminGetOrders)
			//admin.PUT("/processed-order", h.processedOrder) // TODO amount_in_stock handling

			admin.POST(filtrationUrl, h.addFiltrationItem)

			admin.POST(supplyUrl, h.newSupply)
			admin.GET(supplyUrl, h.getAllSupply)
			admin.GET(supplyProductsUrl, h.getSupplyProducts)
			admin.DELETE(supplyUrl, h.deleteSupply)

			superAdmin := admin.Group(superAdminUrl, h.superAdmin)
			{
				superAdmin.GET(inventoryUrl, h.getProductsToInventory)
				superAdmin.POST(inventoryUrl, h.doInventory)
				superAdmin.GET(inventoriesUrl, h.getInventories)
				superAdmin.GET(inventoryProductsUrl, h.getInventoryProducts)
			}

			//admin.POST("write-off") // TODO product write off
			admin.DELETE(reviewUrl, h.deleteReview)
		}

		client := api.Group(clientUrl, h.userAuthorized)
		{
			client.GET("personal-info", h.personalInfo)
			client.PUT(changePasswordUrl, h.changePassword)
			client.DELETE(logoutUrl, h.logout)

			client.GET(userOrdersUrl, h.userOrders)

			client.POST(cartUrl, h.addToCart)
			client.GET(cartUrl, h.getFromCartById)
			client.DELETE(cartUrl, h.deleteFromCart)

			client.POST(favouritesUrl, h.addToFavourites)
			client.GET(favouritesUrl, h.getFavourites)
			client.DELETE(favouritesUrl, h.deleteFromFavourites)

			client.GET(orderUrl, h.getOrderById)
		}
	}

	return router
}
