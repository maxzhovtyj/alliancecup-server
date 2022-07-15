package handler

import (
	server "allincecup-server"
	"github.com/gin-gonic/gin"
	"net/http"
)

type addProductInput struct {
	Info        server.Product       `json:"info"`
	Description []server.ProductInfo `json:"description"`
}

func (h *Handler) addProduct(ctx *gin.Context) {
	var input addProductInput

	if err := ctx.BindJSON(&input); err != nil {
		newErrorResponse(ctx, http.StatusBadRequest, err.Error())
		return
	}

	id, err := h.services.Products.AddProduct(input.Info, input.Description)
	if err != nil {
		newErrorResponse(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	ctx.JSON(http.StatusOK, map[string]interface{}{
		"id":      id,
		"message": "product added",
	})
}
