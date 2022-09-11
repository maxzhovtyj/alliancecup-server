package handler

import "github.com/gin-gonic/gin"

func (h *Handler) doStockInventory(ctx *gin.Context) {
	// TODO
	panic("implement me")
}

func (h *Handler) getInventory(ctx *gin.Context) {
	h.services.Inventory.Products()
}
