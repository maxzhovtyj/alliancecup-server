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
		auth.POST("/refresh-tokens")
	}

	api := router.Group("/api", h.userIdentity)
	{
		admin := api.Group("/admin", h.userHasPermission)
		{
			admin.POST("/add-product", h.addProduct)
			admin.POST("/add-category")
			admin.DELETE("/delete-product")
			admin.DELETE("/delete-category")
			admin.GET("/all-orders")
		}
		client := api.Group("/client")
		{
			client.GET("/all-categories")
			client.GET("/products/:category_id")
			orders := client.Group("/orders")
			{
				orders.POST("/new-order")
				orders.DELETE("/cancel-order")

			}
		}
	}

	return router
}
