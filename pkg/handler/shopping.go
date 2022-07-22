package handler

import (
	server "allincecup-server"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

type AddToCartInput struct {
	UserId  int `json:"user_id"`
	Product server.ProductOrder
}

func (h *Handler) addToCart(ctx *gin.Context) {
	var input AddToCartInput

	if err := ctx.BindJSON(&input); err != nil {
		newErrorResponse(ctx, http.StatusBadRequest, err.Error())
		return
	}

	price, err := h.services.Shopping.AddToCart(input.UserId, input.Product)
	if err != nil {
		newErrorResponse(ctx, http.StatusInternalServerError, err.Error())
	}

	ctx.JSON(http.StatusCreated, map[string]interface{}{
		"price_for_quantity": price,
		"message":            "product added",
	})
}

func (h *Handler) getFromCartById(ctx *gin.Context) {
	userId, err := strconv.Atoi(ctx.Query("user_id"))

	if err != nil {
		newErrorResponse(ctx, http.StatusBadRequest, err.Error())
		return
	}

	products, sum, err := h.services.Shopping.GetProductsInCart(userId)
	if err != nil {
		newErrorResponse(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	ctx.JSON(http.StatusOK, map[string]interface{}{
		"products": products,
		"sum":      sum,
	})
}
