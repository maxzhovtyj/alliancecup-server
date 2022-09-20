package handler

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func (h *Handler) doStockInventory(ctx *gin.Context) {
	// TODO
	panic("implement me")
}

func (h *Handler) getInventory(ctx *gin.Context) {
	products, err := h.services.Inventory.Products()
	if err != nil {
		newErrorResponse(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	ctx.JSON(http.StatusOK, products)
}
