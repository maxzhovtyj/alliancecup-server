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
		api.GET("/all-categories", h.getCategories)
		api.GET("/get-products", h.getProducts)
		api.GET("/product", h.getProductById)
		api.POST("/new-order", h.newOrder)
		api.GET("/order-info-types", h.deliveryPaymentTypes)

		admin := api.Group("/admin", h.userHasPermission)
		{
			admin.POST("/new-moderator", h.createModerator)

			admin.POST("/add-product", h.addProduct)
			admin.PUT("/update-product", h.updateProduct)
			//admin.PUT("/update-product-amount") // TODO
			admin.DELETE("/delete-product", h.deleteProduct)

			admin.POST("/add-category", h.addCategory)
			admin.PUT("/update-category", h.updateCategory)
			admin.DELETE("/delete-category", h.deleteCategory)

			admin.GET("/all-orders", h.adminGetOrders)
			admin.PUT("/processed-order", h.processedOrder)
		}

		client := api.Group("/client", h.userAuthorized)
		{
			client.DELETE("/logout", h.logout)

			client.GET("/user-orders", h.userOrders)

			client.POST("/add-to-cart", h.addToCart)
			client.GET("/user-cart", h.getFromCartById)
			client.DELETE("/delete-from-cart", h.deleteFromCart)

			client.POST("/add-to-favourites", h.addToFavourites)
			client.GET("/get-favourites", h.getFavourites)
			client.DELETE("/delete-from-favourites", h.deleteFromFavourites)

			client.GET("/get-order", h.getOrderById)
		}
	}

	return router
}
