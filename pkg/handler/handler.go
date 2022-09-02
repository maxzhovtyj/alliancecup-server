package handler

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/swaggo/files"       // swagger embed files
	"github.com/swaggo/gin-swagger" // gin-swagger middleware
	_ "github.com/zh0vtyj/allincecup-server/docs"
	"github.com/zh0vtyj/allincecup-server/pkg/service"
	"net/http"
)

type Handler struct {
	services *service.Service
}

func NewHandler(services *service.Service) *Handler {
	return &Handler{services: services}
}

const (
	StatusInProgress   = "IN_PROGRESS"
	StatusProcessed    = "PROCESSED"
	StatusCompleted    = "COMPLETED"
	refreshTokenCookie = "refresh_token"
	domain             = "localhost"
)

const (
	authUrl           = "/auth"
	apiUrl            = "/api"
	adminUrl          = "/admin"
	clientUrl         = "/client"
	signInUrl         = "/sign-in"
	signUpUrl         = "/sign-up"
	logoutUrl         = "/logout"
	changePasswordUrl = "/change-password"
	refreshUrl        = "/refresh"
	categoriesUrl     = "/categories"
	categoryUrl       = "/category"
	productsUrl       = "/products"
	productUrl        = "/product"
	supplyUrl         = "/supply"
	reviewsUrl        = "/reviewUrl"
	reviewUrl         = "/review"
	cartUrl           = "/cart"
	favouritesUrl     = "/favourites"
	filtrationUrl     = "/filtration"
	ordersUrl         = "/orders"
	orderUrl          = "/order"
	userOrdersUrl     = "/user-orders"
	orderInfoTypesUrl = "/order-info-types"
	moderatorUrl      = "/moderator"
)

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
		api.POST(reviewUrl)
		api.GET(reviewsUrl)

		admin := api.Group(adminUrl, h.userHasPermission)
		{
			admin.POST(moderatorUrl, h.createModerator)

			admin.POST(productUrl, h.addProduct)
			admin.PUT(productUrl, h.updateProduct)
			admin.DELETE(productUrl, h.deleteProduct)
			//admin.PUT("/update-product-amount") // TODO

			admin.POST(categoryUrl, h.addCategory)
			admin.PUT(categoryUrl, h.updateCategory)
			admin.DELETE(categoryUrl, h.deleteCategory)

			admin.GET(ordersUrl, h.adminGetOrders)
			//admin.PUT("/processed-order", h.processedOrder) // TODO amount_in_stock handling

			admin.POST(filtrationUrl, h.addFiltrationItem)

			admin.POST(supplyUrl, h.newSupply)
			admin.GET(supplyUrl, h.getAllSupply)
			admin.DELETE(supplyUrl, h.deleteSupply)

			admin.DELETE(reviewUrl)
		}

		client := api.Group(clientUrl, h.userAuthorized)
		{
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
