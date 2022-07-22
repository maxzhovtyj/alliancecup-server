package handler

import (
	server "allincecup-server"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

func (h *Handler) addProduct(ctx *gin.Context) {
	var input server.ProductInfoDescription

	if err := ctx.BindJSON(&input); err != nil {
		newErrorResponse(ctx, http.StatusBadRequest, err.Error())
		return
	}

	id, err := h.services.Products.AddProduct(input.Info, input.Description)
	if err != nil {
		newErrorResponse(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	ctx.JSON(http.StatusCreated, map[string]interface{}{
		"id":      id,
		"message": "product added",
	})
}

func (h *Handler) getProductById(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Query("id"))

	if err != nil {
		newErrorResponse(ctx, http.StatusBadRequest, "id was not found in request path")
		return
	}

	product, err := h.services.Products.GetProductById(id)
	if err != nil {
		newErrorResponse(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	ctx.JSON(http.StatusOK, product)
}

func (h *Handler) updateProduct(ctx *gin.Context) {
	var input server.ProductInfoDescription

	if err := ctx.BindJSON(&input); err != nil {
		newErrorResponse(ctx, http.StatusBadRequest, err.Error())
		return
	}

	id, err := h.services.Products.Update(input)
	if err != nil {
		newErrorResponse(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	ctx.JSON(http.StatusOK, map[string]interface{}{
		"product_id": id,
		"message":    "product updated",
	})
}

func (h *Handler) deleteProduct(ctx *gin.Context) {
	type ProductIdInput struct {
		Id int `json:"id"`
	}

	var input ProductIdInput

	err := ctx.BindJSON(&input)
	if err != nil {
		newErrorResponse(ctx, http.StatusBadRequest, err.Error())
		return
	}

	err = h.services.Products.Delete(input.Id)
	if err != nil {
		newErrorResponse(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	ctx.JSON(http.StatusOK, map[string]interface{}{
		"id":      input.Id,
		"message": "deleted",
	})
}
