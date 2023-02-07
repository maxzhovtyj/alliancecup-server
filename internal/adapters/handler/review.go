package handler

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/zh0vtyj/alliancecup-server/internal/domain/review"
	"net/http"
	"strconv"
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

// deleteReview godoc
// @Summary DeleteReview
// @Security ApiKeyAuth
// @Tags api/admin
// @Description deletes review by its id
// @ID delete review
// @Produce json
// @Param reviewId query int true "review id"
// @Success 200 {object} handler.ItemProcessedResponse
// @Failure 400 {object} Error
// @Failure 403 {object} Error
// @Failure 404 {object} Error
// @Failure 500 {object} Error
// @Router /api/review [delete]
func (h *Handler) deleteReview(ctx *gin.Context) {
	reviewId, err := strconv.Atoi(ctx.Query("id"))
	if err != nil {
		newErrorResponse(ctx, http.StatusBadRequest, err.Error())
		return
	}

	err = h.services.Review.Delete(reviewId)
	if err != nil {
		newErrorResponse(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	ctx.JSON(http.StatusOK, ItemProcessedResponse{
		Id:      reviewId,
		Message: "review successfully deleted",
	})
}

// getReviews godoc
// @Summary GetReviews
// @Tags api
// @Description gets reviews by product id
// @ID gets reviews
// @Accept json
// @Produce json
// @Param productId query int false "product id"
// @Param createAt query string false "last review createdAt"
// @Success 200 {array} review.SelectReviewsDTO
// @Failure 400 {object} Error
// @Failure 404 {object} Error
// @Failure 500 {object} Error
// @Router /api/reviews [get]
func (h *Handler) getReviews(ctx *gin.Context) {
	productId := ctx.Query("productId")

	var productIdInt int
	var err error
	if productId != "" {
		productIdInt, err = strconv.Atoi(productId)
		if err != nil {
			newErrorResponse(ctx, http.StatusBadRequest, err.Error())
			return
		}
	}

	createdAt := ctx.Query("createdAt")

	reviews, err := h.services.Review.Get(createdAt, productIdInt)
	if err != nil {
		newErrorResponse(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	ctx.JSON(http.StatusOK, reviews)
}
