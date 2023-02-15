package handler

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/zh0vtyj/alliancecup-server/internal/domain/category"
	"net/http"
	"strconv"
)

const (
	categoryIdName       = "category_id"
	filtrationListIdName = "filtration_list_id"
)

type DeleteCategoryInput struct {
	Id            int    `json:"id"`
	CategoryTitle string `json:"category_title"`
}

// getCategory godoc
// @Summary      Get category
// @Tags         api
// @Description  get category
// @ID get categories
// @Accept       json
// @Produce      json
// @Param        id query int true "Category id"
// @Success      200  {object}  category.Category
// @Failure      500  {object}  Error
// @Router       /api/category [get]
func (h *Handler) getCategory(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Query("id"))
	if err != nil {
		newErrorResponse(ctx, http.StatusBadRequest, fmt.Errorf("invalid id, %v", err).Error())
		return
	}

	categoryItem, err := h.services.Category.Get(id)
	if err != nil {
		newErrorResponse(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	ctx.JSON(http.StatusOK, categoryItem)
}

// getCategories godoc
// @Summary      GetCategories
// @Tags         api
// @Description  get all categories
// @ID get categories
// @Accept       json
// @Produce      json
// @Success      200  {array}  category.Category
// @Failure      500  {object}  Error
// @Router       /api/categories [get]
func (h *Handler) getCategories(ctx *gin.Context) {
	categories, err := h.services.Category.GetAll()
	if err != nil {
		newErrorResponse(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	ctx.JSON(http.StatusOK, categories)
}

// addCategory godoc
// @Summary AddCategory
// @Security ApiKeyAuth
// @Tags api/admin
// @Description Adds a new category
// @ID adds category
// @Accept json
// @Produce json
// @Param input body category.Category true "category info" // TODO
// @Success 201 {object} handler.ItemProcessedResponse
// @Failure 400 {object} Error
// @Failure 500 {object} Error
// @Router  /api/admin/category [post]
func (h *Handler) addCategory(ctx *gin.Context) {
	ctx.Writer.Header().Set("Content-Type", "form/json")
	err := ctx.Request.ParseMultipartForm(fileMaxSize)
	if err != nil {
		newErrorResponse(ctx, http.StatusBadRequest, err.Error())
		return
	}

	var dto category.CreateDTO

	dto.CategoryTitle = ctx.Request.Form.Get("categoryTitle")
	if dto.CategoryTitle == "" {
		newErrorResponse(ctx, http.StatusBadRequest, "category title is empty")
		return
	}

	description := ctx.Request.Form.Get("description")
	if description != "" {
		dto.CategoryDescription = &description
	}

	imgUrl := ctx.Request.Form.Get("imgUrl")
	if imgUrl != "" {
		dto.ImgUrl = &imgUrl
	}

	dto.Img, err = parseFile(ctx.Request.MultipartForm.File)
	if err != nil {
		if !errors.Is(err, ErrEmptyFile) {
			newErrorResponse(ctx, http.StatusBadRequest, err.Error())
			return
		}
	}

	id, err := h.services.Category.Create(dto)
	if err != nil {
		newErrorResponse(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	ctx.JSON(http.StatusCreated, ItemProcessedResponse{
		Id:      id,
		Message: "category successfully created",
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
// @Failure      500  {object}  Error
// @Router       /api/admin/category [put]
func (h *Handler) updateCategory(ctx *gin.Context) {
	var input category.Category

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

// updateCategoryImage godoc
// @Summary      Update category image
// @Security 	 ApiKeyAuth
// @Tags         api/admin
// @Description  Updates category image
// @ID 			 updates category image
// @Accept 	     json
// @Produce      json
// @Param        input body category.Category true "category info"
// @Success      200  {object}  handler.ItemProcessedResponse
// @Failure      400  {object}  Error
// @Failure      500  {object}  Error
// @Router       /api/admin/category [put]
func (h *Handler) updateCategoryImage(ctx *gin.Context) {
	ctx.Writer.Header().Set("Content-Type", "form/json")
	err := ctx.Request.ParseMultipartForm(fileMaxSize)
	if err != nil {
		newErrorResponse(ctx, http.StatusBadRequest, err.Error())
		return
	}

	var dto category.UpdateImageDTO

	dto.Id, err = strconv.Atoi(ctx.Request.Form.Get("id"))
	if err != nil {
		newErrorResponse(ctx, http.StatusBadRequest, fmt.Sprintf("invalid id: %d", dto.Id))
		return
	}

	dto.Img, err = parseFile(ctx.Request.MultipartForm.File)
	if err != nil {
		newErrorResponse(ctx, http.StatusBadRequest, err.Error())
		return
	}

	id, err := h.services.Category.UpdateImage(dto)
	if err != nil {
		newErrorResponse(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	ctx.JSON(http.StatusOK, ItemProcessedResponse{
		Id:      id,
		Message: "category image updated",
	})
}

// deleteCategoryImage godoc
// @Summary      Delete category image (Minio)
// @Security 	 ApiKeyAuth
// @Tags         api/admin
// @Description  Deletes category image (Minio)
// @ID 			 deletes category image (Minio)
// @Accept 	     json
// @Produce      json
// @Param        input body category.Category true "category info"
// @Success      200  {object}  handler.ItemProcessedResponse
// @Failure      400  {object}  Error
// @Failure      500  {object}  Error
// @Router       /api/admin/category-image [delete]
func (h *Handler) deleteCategoryImage(ctx *gin.Context) {
	categoryId, err := strconv.Atoi(ctx.Query("id"))
	if err != nil {
		newErrorResponse(ctx, http.StatusBadRequest, fmt.Errorf("invalid id, %v", err).Error())
		return
	}

	err = h.services.Category.DeleteImage(categoryId)
	if err != nil {
		newErrorResponse(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	ctx.JSON(http.StatusOK, map[string]any{
		"id":      categoryId,
		"message": "category image (Minio) deleted",
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
