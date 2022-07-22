package handler

import (
	"allincecup-server/pkg/service"
	"github.com/gin-gonic/gin"
)

type Handler struct {
	services *service.Service
}

func NewHandler(services *service.Service) *Handler {
	return &Handler{services: services}
}

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

			admin.GET("/all-orders")
		}

		orders := api.Group("/orders")
		{
			orders.POST("/new-order")
		}

		client := api.Group("/client", h.userAuthorized)
		{
			client.GET("/users-orders")
			client.GET("/users-cart")
		}
	}

	return router
}
