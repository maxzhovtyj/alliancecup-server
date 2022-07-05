package handler

import (
	server "allincecup-server"
	"github.com/gin-gonic/gin"
	"net/http"
)

func (h *Handler) addProduct(ctx *gin.Context) {
	var input server.Product

	if err := ctx.BindJSON(&input); err != nil {
		newErrorResponse(ctx, http.StatusBadRequest, err.Error())
		return
	}

	id, err := h.services.Products.AddProduct(input)
	if err != nil {
		newErrorResponse(ctx, http.StatusInternalServerError, err.Error())
		return
	}
	ctx.JSON(http.StatusOK, map[string]interface{}{
		"id":      id,
		"message": "product added",
	})
}
