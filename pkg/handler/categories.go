package handler

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func (h *Handler) addCategory(ctx *gin.Context) {
	type CategoryInput struct {
		Title string `json:"title"`
	}

	var input CategoryInput

	if err := ctx.BindJSON(&input); err != nil {
		newErrorResponse(ctx, http.StatusBadRequest, err.Error())
		return
	}

	id, err := h.services.Category.Create(input.Title)
	if err != nil {
		newErrorResponse(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	ctx.JSON(http.StatusOK, map[string]interface{}{
		"category_id": id,
	})
}

func (h *Handler) deleteCategory(ctx *gin.Context) {
}
