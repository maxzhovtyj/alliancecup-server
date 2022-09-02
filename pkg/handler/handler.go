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

	auth := router.Group("/auth")
	{
		auth.POST("/sign-up", h.signUp)
		auth.POST("/sign-in", h.signIn)
		auth.POST("/refresh", h.refresh)
	}

	api := router.Group("/api", h.userIdentity)
	{
		api.GET("/categories", h.getCategories)
		api.GET("/filtration-list", h.getFiltration)
		api.GET("/products", h.getProducts)
		api.GET("/product", h.getProductById)
		api.POST("/new-order", h.newOrder)
		api.GET("/order-info-types", h.deliveryPaymentTypes)
		api.POST("/review")
		api.GET("/reviews")

		admin := api.Group("/admin", h.userHasPermission)
		{
			admin.POST("/moderator", h.createModerator)

			admin.POST("/product", h.addProduct)
			admin.PUT("/product", h.updateProduct)
			admin.DELETE("/product", h.deleteProduct)
			//admin.PUT("/update-product-amount") // TODO

			admin.POST("/category", h.addCategory)
			admin.PUT("/category", h.updateCategory)
			admin.DELETE("/category", h.deleteCategory)

			admin.GET("/orders", h.adminGetOrders)
			//admin.PUT("/processed-order", h.processedOrder) // TODO amount_in_stock handling

			admin.POST("/filtration", h.addFiltrationItem)

			admin.POST("/supply", h.newSupply)
			admin.GET("/supply", h.getAllSupply)
			admin.DELETE("/supply", h.deleteSupply)
		}

		client := api.Group("/client", h.userAuthorized)
		{
			client.PUT("/change-password", h.changePassword)
			client.DELETE("/logout", h.logout)

			client.GET("/user-order", h.userOrders)

			client.POST("/cart", h.addToCart)
			client.GET("/cart", h.getFromCartById)
			client.DELETE("/cart", h.deleteFromCart)

			client.POST("/favourites", h.addToFavourites)
			client.GET("/favourites", h.getFavourites)
			client.DELETE("favourites", h.deleteFromFavourites)

			client.GET("/get-order", h.getOrderById)
		}
	}

	return router
}
