package handler

import (
	"errors"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/swaggo/files"       // swagger embed files
	"github.com/swaggo/gin-swagger" // gin-swagger middleware
	_ "github.com/zh0vtyj/alliancecup-server/docs"
	"github.com/zh0vtyj/alliancecup-server/internal/config"
	"github.com/zh0vtyj/alliancecup-server/internal/domain/service"
	"github.com/zh0vtyj/alliancecup-server/pkg/logging"
	"net/http"
)

const (
	StatusInProgress   = "IN_PROGRESS"
	StatusProcessed    = "PROCESSED"
	StatusCompleted    = "COMPLETED"
	refreshTokenCookie = "refresh_token"
	fileMaxSize        = 5 << 20
)

const (
	apiUrl               = "/api"
	adminUrl             = "/admin"
	clientUrl            = "/client"
	logoutUrl            = "/logout"
	changePasswordUrl    = "/change-password"
	categoriesUrl        = "/categories"
	categoryUrl          = "/category"
	categoryImageUrl     = "/category-image"
	productsUrl          = "/products"
	productUrl           = "/product"
	productImageUrl      = "/product-image"
	productVisibilityUrl = "/product-visibility"
	reviewsUrl           = "/reviews"
	reviewUrl            = "/review"
	cartUrl              = "/cart"
	favouritesUrl        = "/favourites"
	filtrationUrl        = "/filtration"
	filtrationListUrl    = "/filtration-list"
	filtrationItemUrl    = "/filtration-item"
	filtrationImageUrl   = "/filtration-image"
	ordersUrl            = "/orders"
	orderUrl             = "/order"
	userOrdersUrl        = "/user-orders"
	orderInfoTypesUrl    = "/order-info-types"
	processedOrder       = "/processed-order"
	completeOrder        = "/complete-order"

	superAdminUrl        = "/super"
	supplyUrl            = "/supply"
	supplyProductsUrl    = "/supply-products"
	inventoryUrl         = "/inventory"
	inventoriesUrl       = "/inventories"
	inventoryProductsUrl = "/inventory-products"
	saveInventory        = "/save-inventory"
	invoiceUrl           = "/invoice"
	personalInfoUrl      = "/personal-info"
	shoppingUrl          = "/shopping"
	forgotPassword       = "/forgot-password"
	restorePasswordUrl   = "/restore-password"
)

var ErrEmptyFile = errors.New("file is empty")

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

func (h *Handler) InitRoutes(cfg *config.Config) *gin.Engine {
	router := gin.New()

	corsConfig := cors.Config{
		AllowOrigins: cfg.Cors.AllowedOrigins,
		AllowMethods: []string{
			http.MethodGet, http.MethodDelete, http.MethodPost, http.MethodPut, http.MethodPatch,
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
	}

	c := cors.New(corsConfig)
	router.Use(c, gin.Recovery())

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	h.initAuthRoutes(router)
	h.initApi(router)

	return router
}

func (h *Handler) initApi(router *gin.Engine) {
	api := router.Group(apiUrl, h.userIdentity)
	{
		api.GET(categoryUrl, h.getCategory)
		api.GET(categoriesUrl, h.getCategories)
		api.GET(filtrationUrl, h.getFiltration)

		api.GET(productsUrl, h.getProducts)
		api.GET(productUrl, h.getProductById)

		h.initReviewsRoutes(api)

		api.POST(forgotPassword, h.forgotPassword)

		api.GET(invoiceUrl, h.getOrderInvoice)

		admin := api.Group(adminUrl, h.moderatorPermission)
		{
			h.initAdminProductsRoutes(admin)
			h.initAdminCategoriesRoutes(admin)
			h.initAdminFiltrationRoutes(admin)
			h.initAdminOrderRoutes(admin)

			admin.POST(supplyUrl, h.newSupply)
			admin.GET(supplyUrl, h.getAllSupply)
			admin.GET(supplyProductsUrl, h.getSupplyProducts)

			superAdmin := admin.Group(superAdminUrl, h.superAdmin)
			{
				superAdmin.DELETE(supplyUrl, h.deleteSupply)

				h.initAdminModeratorsRoutes(superAdmin)
				h.initAdminInventoryRoutes(superAdmin)
			}

			admin.DELETE(reviewUrl, h.deleteReview)
		}

		client := api.Group(clientUrl, h.userAuthorized)
		{
			h.initClientRoutes(client)

			client.GET(userOrdersUrl, h.userOrders)
		}

		shopping := api.Group(shoppingUrl, h.getShoppingInfo)
		{
			shopping.GET(orderInfoTypesUrl, h.deliveryPaymentTypes)
			shopping.POST(orderUrl, h.newOrder)

			h.initShoppingShoppingRoutes(shopping)
		}
	}
}
