package handler

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/swaggo/files"       // swagger embed files
	"github.com/swaggo/gin-swagger" // gin-swagger middleware
	_ "github.com/zh0vtyj/allincecup-server/docs"
	"github.com/zh0vtyj/allincecup-server/internal/config"
	"github.com/zh0vtyj/allincecup-server/internal/domain/service"
	"github.com/zh0vtyj/allincecup-server/pkg/logging"
	"net/http"
)

const (
	StatusInProgress   = "IN_PROGRESS"
	StatusProcessed    = "PROCESSED"
	StatusCompleted    = "COMPLETED"
	refreshTokenCookie = "refresh_token"
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
	processedOrder       = "/processed-order"
	forgotPassword       = "/forgot-password"
	moderatorUrl         = "/moderator"
	superAdminUrl        = "/super"
	supplyUrl            = "/supply"
	supplyProductsUrl    = "/supply-products"
	inventoryUrl         = "/inventory"
	inventoriesUrl       = "/inventories"
	inventoryProductsUrl = "/inventory-products"
	saveInventory        = "/save-inventory"
	invoiceUrl           = "/invoice"
	personalInfoUrl      = "personal-info"
	shoppingUrl          = "/shopping"
)

type Handler struct {
	services *service.Service
	logger   *logging.Logger
	cfg      *config.Config
}

func NewHandler(services *service.Service, logger *logging.Logger, cfg *config.Config) *Handler {
	return &Handler{
		services: services,
		logger:   logger,
		cfg:      cfg,
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
			"UserCart",
			"UserFavourites",
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
		api.POST(reviewUrl, h.addReview)
		api.GET(reviewsUrl, h.getReviews)

		api.POST(forgotPassword, h.forgotPassword)

		api.GET(invoiceUrl, h.getOrderInvoice)

		admin := api.Group(adminUrl, h.moderatorPermission)
		{
			admin.POST(productUrl, h.addProduct)
			admin.PUT(productUrl, h.updateProduct)
			admin.DELETE(productUrl, h.deleteProduct)

			admin.POST(categoryUrl, h.addCategory)
			admin.PUT(categoryUrl, h.updateCategory)
			admin.DELETE(categoryUrl, h.deleteCategory)

			admin.GET(ordersUrl, h.adminGetOrders)
			admin.PUT(processedOrder, h.processedOrder)

			admin.POST(filtrationUrl, h.addFiltrationItem)

			admin.POST(supplyUrl, h.newSupply)
			admin.GET(supplyUrl, h.getAllSupply)
			admin.GET(supplyProductsUrl, h.getSupplyProducts)
			admin.DELETE(supplyUrl, h.deleteSupply)

			admin.GET(orderUrl, h.getOrderById)
			admin.POST(orderUrl, h.adminNewOrder)

			superAdmin := admin.Group(superAdminUrl, h.superAdmin)
			{
				superAdmin.GET(moderatorUrl, h.getModerators)
				superAdmin.POST(moderatorUrl, h.createModerator)
				superAdmin.DELETE(moderatorUrl, h.deleteModerator)

				superAdmin.GET(inventoryUrl, h.getProductsToInventory)
				superAdmin.POST(saveInventory, h.saveInventory)
				superAdmin.POST(inventoryUrl, h.doInventory)
				superAdmin.GET(inventoriesUrl, h.getInventories)
				superAdmin.GET(inventoryProductsUrl, h.getInventoryProducts)
			}

			//admin.POST("write-off") // TODO product write off
			admin.DELETE(reviewUrl, h.deleteReview)
		}

		client := api.Group(clientUrl, h.userAuthorized)
		{
			client.GET(personalInfoUrl, h.personalInfo)
			client.PUT(personalInfoUrl, h.updatePersonalInfo)

			client.PUT(changePasswordUrl, h.changePassword)
			client.DELETE(logoutUrl, h.logout)

			client.GET(userOrdersUrl, h.userOrders)
		}

		shopping := api.Group(shoppingUrl, h.getShoppingInfo)
		{
			shopping.GET(orderInfoTypesUrl, h.deliveryPaymentTypes)
			shopping.POST(orderUrl, h.newOrder)

			shopping.GET(cartUrl, h.getFromCartById)
			shopping.POST(cartUrl, h.addToCart)
			shopping.DELETE(cartUrl, h.deleteFromCart)

			shopping.GET(favouritesUrl, h.getFavourites)
			shopping.POST(favouritesUrl, h.addToFavourites)
			shopping.DELETE(favouritesUrl, h.deleteFromFavourites)
		}
	}

	return router
}
