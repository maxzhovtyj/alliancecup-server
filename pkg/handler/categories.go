package handler

import (
	server "allincecup-server"
	"github.com/gin-gonic/gin"
	"net/http"
)

type allCategoriesResponse struct {
	Data []server.Category `json:"data"`
}

func (h *Handler) getCategories(ctx *gin.Context) {
	categories, err := h.services.Category.GetAll()
	if err != nil {
		newErrorResponse(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	ctx.JSON(http.StatusOK, allCategoriesResponse{Data: categories})
}

func (h *Handler) addCategory(ctx *gin.Context) {
	var input server.Category

	if err := ctx.BindJSON(&input); err != nil {
		newErrorResponse(ctx, http.StatusBadRequest, err.Error())
		return
	}

	id, err := h.services.Category.Create(input)
	if err != nil {
		newErrorResponse(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	ctx.JSON(http.StatusCreated, map[string]interface{}{
		"category_id": id,
	})
}

func (h *Handler) updateCategory(ctx *gin.Context) {
	var input server.Category

	if err := ctx.BindJSON(&input); err != nil {
		newErrorResponse(ctx, http.StatusBadRequest, err.Error())
		return
	}

	id, err := h.services.Category.Update(input)
	if err != nil {
		newErrorResponse(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	ctx.JSON(http.StatusOK, map[string]interface{}{
		"category_id": id,
	})
}

func (h *Handler) deleteCategory(ctx *gin.Context) {}
