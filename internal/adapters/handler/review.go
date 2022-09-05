package handler

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/zh0vtyj/allincecup-server/internal/domain/review"
	"net/http"
)

// addReview godoc
// @Summary AddReview
// @Tags api
// @Description creates a new review
// @ID create review
// @Accept json
// @Produce json
// @Param input body review.CreateReviewDTO true "review info"
// @Success 201 {object} handler.ItemProcessedResponse
// @Failure 400 {object} Error
// @Failure 404 {object} Error
// @Failure 500 {object} Error
// @Router /api/review [post]
func (h *Handler) addReview(ctx *gin.Context) {
	var input review.CreateReviewDTO

	if err := ctx.BindJSON(&input); err != nil {
		newErrorResponse(ctx, http.StatusBadRequest, fmt.Errorf("invalid input data, %v", err).Error())
		return
	}

	id, err := h.services.Review.Create(input)
	if err != nil {
		newErrorResponse(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	ctx.JSON(http.StatusCreated, ItemProcessedResponse{
		Id:      id,
		Message: "review created",
	})
}
