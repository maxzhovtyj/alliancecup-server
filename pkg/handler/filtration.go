package handler

import (
	"github.com/gin-gonic/gin"
	server "github.com/zh0vtyj/allincecup-server/pkg/models"
	"net/http"
)

// addFiltrationItem godoc
// @Summary Add filtration for category
// @Tags api/admin
// @Description Adds a filtration item to a category
// @ID add filtration
// @Accept json
// @Produce json
// @Param input body server.CategoryFiltration true "filtration info"
// @Success 201 {object} string
// @Failure 400 {object} Error
// @Failure 404 {object} Error
// @Failure 500 {object} Error
// @Router /api/admin/filtration-item [post]
func (h *Handler) addFiltrationItem(ctx *gin.Context) {
	var input server.CategoryFiltration

	if err := ctx.BindJSON(&input); err != nil {
		newErrorResponse(ctx, http.StatusBadRequest, err.Error())
		return
	}

	id, err := h.services.Category.AddFiltration(input)
	if err != nil {
		newErrorResponse(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	ctx.JSON(http.StatusCreated, map[string]any{
		"id":      id,
		"message": "filtration list created",
	})
}
