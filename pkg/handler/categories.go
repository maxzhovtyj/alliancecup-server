package handler

import (
	"github.com/gin-gonic/gin"
	server "github.com/zh0vtyj/allincecup-server"
	"net/http"
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
// @Success      200  {object}   allCategoriesResponse
// @Failure      400  {object}  Error
// @Failure      404  {object}  Error
// @Failure      500  {object}  Error
// @Router       /api/all-categories [get]
func (h *Handler) getCategories(ctx *gin.Context) {
	categories, err := h.services.Category.GetAll()
	if err != nil {
		newErrorResponse(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	ctx.JSON(http.StatusOK, allCategoriesResponse{Data: categories})
}

// addCategory godoc
// @Summary      AddCategory
// @Security 	 ApiKeyAuth
// @Tags         api/admin
// @Description  Adds a new category
// @ID 			 adds category
// @Accept 	     json
// @Produce      json
// @Param        input body server.Category true "category info"
// @Success      200  {object}  handler.ItemProcessedResponse
// @Failure      400  {object}  Error
// @Failure      404  {object}  Error
// @Failure      500  {object}  Error
// @Router       /api/admin/add-category [post]
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
// @Param        input body server.Category true "category info"
// @Success      200  {object}  handler.ItemProcessedResponse
// @Failure      400  {object}  Error
// @Failure      404  {object}  Error
// @Failure      500  {object}  Error
// @Router       /api/admin/update-category [put]
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
// @Param        input body DeleteCategoryInput true "category info"
// @Success      200  {object}  handler.ItemProcessedResponse
// @Failure      400  {object}  Error
// @Failure      404  {object}  Error
// @Failure      500  {object}  Error
// @Router       /api/admin/delete-category [delete]
func (h *Handler) deleteCategory(ctx *gin.Context) {
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

	ctx.JSON(http.StatusOK, ItemProcessedResponse{
		Id:      input.Id,
		Message: "category deleted",
	})
}
