package handler

import (
	"github.com/gin-gonic/gin"
	server "github.com/zh0vtyj/allincecup-server"
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

func (h *Handler) deleteCategory(ctx *gin.Context) {
	type DeleteCategoryInput struct {
		Id            int    `json:"id"`
		CategoryTitle string `json:"category_title"`
	}

	var input DeleteCategoryInput

	if err := ctx.BindJSON(&input); err != nil {
		newErrorResponse(ctx, http.StatusBadRequest, err.Error())
		return
	}

	err := h.services.Category.Delete(input.Id, input.CategoryTitle)
	if err != nil {
		newErrorResponse(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	ctx.JSON(http.StatusOK, map[string]interface{}{
		"id":      input.Id,
		"message": "deleted",
	})
}
