package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/zh0vtyj/allincecup-server/internal/domain/inventory"
	"net/http"
)

func (h *Handler) getInventory(ctx *gin.Context) {
	products, err := h.services.Inventory.Products()
	if err != nil {
		newErrorResponse(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	ctx.JSON(http.StatusOK, products)
}

func (h *Handler) doInventory(ctx *gin.Context) {
	var input []inventory.ProductDTO

	if err := ctx.BindJSON(&input); err != nil {
		newErrorResponse(ctx, http.StatusBadRequest, err.Error())
		return
	}

	err := h.services.Inventory.New(input)
	if err != nil {
		newErrorResponse(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	ctx.JSON(http.StatusCreated, map[string]any{
		"message": "inventory created",
	})
}
