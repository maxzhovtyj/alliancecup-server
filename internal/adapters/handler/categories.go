package handler

import (
	"fmt"
	"github.com/gin-gonic/gin"
	server "github.com/zh0vtyj/allincecup-server/internal/domain/category"
	"net/http"
	"strconv"
)

const (
	categoryIdName       = "category_id"
	filtrationListIdName = "filtration_list_id"
)

type allCategoriesResponse struct {
	Data []server.Category `json:"data"`
}

type DeleteCategoryInput struct {
	Id            int    `json:"id"`
	CategoryTitle string `json:"category_title"`
}

// getCategories godoc
// @Summary      GetCategories
// @Tags         api
// @Description  get all categories
// @ID get categories
// @Accept       json
// @Produce      json
// @Success      200  {object}  allCategoriesResponse
// @Failure      400  {object}  Error
// @Failure      404  {object}  Error
// @Failure      500  {object}  Error
// @Router       /api/categories [get]
func (h *Handler) getCategories(ctx *gin.Context) {
	categories, err := h.services.Category.GetAll()
	if err != nil {
		newErrorResponse(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	ctx.JSON(http.StatusOK, allCategoriesResponse{Data: categories})
}

// getFiltration godoc
// @Summary GetFiltration
// @Tags api
// @Description gets filtration list for a products
// @ID get filtration
// @Accept json
// @Produce json
// @Param id query int true "parent id"
// @Param parentName query string true "parent name"
// @Success 200 {array} category.Filtration
// @Failure 400 {object} Error
// @Failure 404 {object} Error
// @Failure 500 {object} Error
// @Router /api/filtration [get]
func (h *Handler) getFiltration(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Query("id"))
	if err != nil {
		newErrorResponse(ctx, http.StatusBadRequest, fmt.Errorf("failed to parse to int id due to %v", err).Error())
		return
	}

	parentName := ctx.Query("parentName")
	if parentName != categoryIdName && parentName != filtrationListIdName {
		newErrorResponse(ctx, http.StatusBadRequest, fmt.Errorf("invalid parent name, can be either category_id or filtration_list_id").Error())
		return
	}

	if err != nil {
		newErrorResponse(ctx, http.StatusBadRequest, err.Error())
		return
	}

	filtrationList, err := h.services.Category.GetFiltration(parentName, id)
	if err != nil {
		newErrorResponse(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	ctx.JSON(http.StatusOK, filtrationList)
}

// addCategory godoc
// @Summary AddCategory
// @Security ApiKeyAuth
// @Tags api/admin
// @Description Adds a new category
// @ID adds category
// @Accept json
// @Produce json
// @Param input body category.Category true "category info"
// @Success 201 {object} handler.ItemProcessedResponse
// @Failure 400 {object} Error
// @Failure 404 {object} Error
// @Failure 500 {object} Error
// @Router  /api/admin/category [post]
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

	ctx.JSON(http.StatusCreated, ItemProcessedResponse{
		Id:      id,
		Message: "category added",
	})
}

// updateCategory godoc
// @Summary      UpdateCategory
// @Security 	 ApiKeyAuth
// @Tags         api/admin
// @Description  Updates category
// @ID 			 updates category
// @Accept 	     json
// @Produce      json
// @Param        input body category.Category true "category info"
// @Success      200  {object}  handler.ItemProcessedResponse
// @Failure      400  {object}  Error
// @Failure      404  {object}  Error
// @Failure      500  {object}  Error
// @Router       /api/admin/category [put]
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

	ctx.JSON(http.StatusOK, ItemProcessedResponse{
		Id:      id,
		Message: "category updated",
	})
}

// deleteCategory godoc
// @Summary      DeleteCategory
// @Security 	 ApiKeyAuth
// @Tags         api/admin
// @Description  Deletes category
// @ID 			 deletes category
// @Accept 	     json
// @Produce      json
// @Param        id query int true "category id"
// @Success      200  {object}  handler.ItemProcessedResponse
// @Failure      400  {object}  Error
// @Failure      404  {object}  Error
// @Failure      500  {object}  Error
// @Router       /api/admin/category [delete]
func (h *Handler) deleteCategory(ctx *gin.Context) {
	// TODO "pq: update or delete on table \"products\" violates foreign key constraint \"orders_products_product_id_fkey\" on table \"orders_products\""

	categoryId, err := strconv.Atoi(ctx.Query("id"))
	if err != nil {
		newErrorResponse(ctx, http.StatusBadRequest, err.Error())
		return
	}

	err = h.services.Category.Delete(categoryId)
	if err != nil {
		newErrorResponse(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	ctx.JSON(http.StatusOK, ItemProcessedResponse{
		Id:      categoryId,
		Message: "category deleted",
	})
}
