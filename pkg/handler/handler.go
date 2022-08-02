package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/zh0vtyj/allincecup-server/pkg/service"
)

type Handler struct {
	services *service.Service
}

func NewHandler(services *service.Service) *Handler {
	return &Handler{services: services}
}

const (
	statusInProgress = "IN_PROGRESS"
	statusProcessed  = "PROCESSED"
	statusCompleted  = "COMPLETED"
)

func (h *Handler) InitRoutes() *gin.Engine {
	router := gin.New()

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

		admin := api.Group("/admin", h.userHasPermission)
		{
			admin.POST("/new-moderator", h.createModerator)

			admin.POST("/add-product", h.addProduct)
			admin.PUT("/update-product", h.updateProduct)
			admin.DELETE("/delete-product", h.deleteProduct)

			admin.POST("/add-category", h.addCategory)
			admin.PUT("/update-category", h.updateCategory)
			admin.DELETE("/delete-category", h.deleteCategory)

			admin.GET("/all-orders", h.adminGetOrders)
		}

		api.POST("/new-order", h.newOrder)

		client := api.Group("/client", h.userAuthorized)
		{
			client.GET("/orders", h.userOrders)
			client.DELETE("/logout", h.logout)

			client.POST("/add-to-cart", h.addToCart)
			client.GET("/user-cart", h.getFromCartById)
			client.DELETE("/delete-from-cart", h.deleteFromCart)

			client.POST("/add-to-favourites", h.addToFavourites)
			client.GET("/get-favourites", h.getFavourites)

			client.GET("/get-order", h.getOrderById)
		}
	}

	return router
}
