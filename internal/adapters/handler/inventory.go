package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/zh0vtyj/allincecup-server/internal/domain/inventory"
	"net/http"
	"strconv"
)

func (h *Handler) getProductsToInventory(ctx *gin.Context) {
	products, err := h.services.Inventory.Products()
	if err != nil {
		newErrorResponse(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	ctx.JSON(http.StatusOK, products)
}

func (h *Handler) doInventory(ctx *gin.Context) {
	var input []inventory.InsertProductDTO

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

func (h *Handler) getInventories(ctx *gin.Context) {
	createdAt := ctx.Query("createdAt")

	inventories, err := h.services.Inventory.GetAll(createdAt)
	if err != nil {
		newErrorResponse(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	ctx.JSON(http.StatusOK, inventories)
}

func (h *Handler) getInventoryProducts(ctx *gin.Context) {
	inventoryId, err := strconv.Atoi(ctx.Query("id"))
	if err != nil {
		newErrorResponse(ctx, http.StatusBadRequest, "invalid inventory id")
		return
	}

	products, err := h.services.Inventory.GetInventoryProducts(inventoryId)
	if err != nil {
		newErrorResponse(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	ctx.JSON(http.StatusOK, products)
}
