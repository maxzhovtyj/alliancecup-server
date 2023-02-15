package handler

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/zh0vtyj/alliancecup-server/internal/domain/category"
	"net/http"
	"strconv"
)

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

// getFiltrationAllItems godoc
// @Summary Get all filtration items
// @Tags api/admin/characteristics
// @Description gets all filtration items
// @ID get filtration items
// @Accept json
// @Produce json
// @Success 200 {array}  category.Filtration
// @Failure 400 {object} Error
// @Failure 500 {object} Error
// @Router /api/admin/characteristics [get]
func (h *Handler) getFiltrationAllItems(ctx *gin.Context) {
	filtrationItems, err := h.services.Category.GetFiltrationItems()
	if err != nil {
		newErrorResponse(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	ctx.JSON(http.StatusOK, filtrationItems)
}

// getFiltrationItem godoc
// @Summary Get filtration
// @Security ApiKeyAuth
// @Tags api/admin
// @Description Adds a filtration item to a category
// @ID add filtration
// @Accept json
// @Produce json
// @Param id query int true "filtration item id"
// @Success 200 {object} string
// @Failure 400 {object} Error
// @Failure 404 {object} Error
// @Failure 500 {object} Error
// @Router /api/admin/filtration-item [get]
func (h *Handler) getFiltrationItem(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Query("id"))
	if err != nil {
		newErrorResponse(ctx, http.StatusBadRequest, fmt.Errorf("invalid id, %v", err).Error())
		return
	}

	filtrationItem, err := h.services.Category.GetFiltrationItem(id)
	if err != nil {
		newErrorResponse(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	ctx.JSON(http.StatusOK, filtrationItem)
}

// addFiltrationItem godoc
// @Summary Add filtration
// @Security ApiKeyAuth
// @Tags api/admin
// @Description Adds a filtration item to a category
// @ID add filtration
// @Accept json
// @Produce json
// @Param input body category.Filtration true "filtration info"
// @Success 201 {object} string
// @Failure 400 {object} Error
// @Failure 404 {object} Error
// @Failure 500 {object} Error
// @Router /api/admin/filtration [post]
func (h *Handler) addFiltrationItem(ctx *gin.Context) {
	ctx.Writer.Header().Set("Content-Type", "form/json")
	err := ctx.Request.ParseMultipartForm(fileMaxSize)
	if err != nil {
		newErrorResponse(ctx, http.StatusBadRequest, err.Error())
		return
	}

	var filtration category.CreateFiltrationDTO

	categoryId := ctx.Request.Form.Get("categoryId")
	if categoryId != "" {
		categoryIdInt, err := strconv.Atoi(categoryId)
		if err != nil {
			newErrorResponse(ctx, http.StatusBadRequest, fmt.Errorf("invalid category id, %v", err).Error())
			return
		}

		filtration.CategoryId = &categoryIdInt
	}

	imgUrl := ctx.Request.Form.Get("imgUrl")
	if imgUrl != "" {
		filtration.ImgUrl = &imgUrl
	}

	filtration.SearchKey = ctx.Request.Form.Get("searchKey")
	if filtration.SearchKey == "" {
		newErrorResponse(ctx, http.StatusBadRequest, fmt.Errorf("invalid search key").Error())
		return
	}

	filtration.SearchCharacteristic = ctx.Request.Form.Get("searchCharacteristic")
	if filtration.SearchCharacteristic == "" {
		newErrorResponse(ctx, http.StatusBadRequest, fmt.Errorf("invalid search characteristic").Error())
		return
	}

	filtration.FiltrationTitle = ctx.Request.Form.Get("filtrationTitle")
	if filtration.FiltrationTitle == "" {
		newErrorResponse(ctx, http.StatusBadRequest, fmt.Errorf("invalid filtration title").Error())
		return
	}
	filtrationDescription := ctx.Request.Form.Get("filtrationDescription")
	if filtrationDescription != "" {
		filtration.FiltrationDescription = &filtrationDescription
	}

	filtrationListId := ctx.Request.Form.Get("filtrationListId")
	if filtrationListId != "" {
		filtrationListIdInt, err := strconv.Atoi(filtrationListId)
		if err != nil {
			newErrorResponse(
				ctx,
				http.StatusBadRequest,
				fmt.Errorf("failed parse to int filtration list id %v", err).Error())
			return
		}

		filtration.FiltrationListId = &filtrationListIdInt
	}

	filtration.Img, err = parseFile(ctx.Request.MultipartForm.File)
	if err != nil {
		if !errors.Is(err, ErrEmptyFile) {
			newErrorResponse(ctx, http.StatusBadRequest, err.Error())
			return
		}
	}

	if (filtration.FiltrationListId != nil && filtration.CategoryId != nil) ||
		(filtration.FiltrationListId == nil && filtration.CategoryId == nil) {
		newErrorResponse(
			ctx,
			http.StatusInternalServerError,
			fmt.Errorf("filtration item must have either list id or category id").Error(),
		)
		return
	}

	id, err := h.services.Category.AddFiltration(filtration)
	if err != nil {
		newErrorResponse(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	ctx.JSON(http.StatusCreated, map[string]any{
		"id":      id,
		"message": "filtration list created",
	})
}

// updateFiltrationItemImage godoc
// @Summary Update filtration image (Minio)
// @Security ApiKeyAuth
// @Tags api/admin
// @Description Updates filtration item image (Minio)
// @ID update filtration
// @Accept json
// @Produce json
// @Param input body category.Filtration true "filtration info"
// @Success 201 {object} string
// @Failure 400 {object} Error
// @Failure 404 {object} Error
// @Failure 500 {object} Error
// @Router /api/admin/filtration-image [put]
func (h *Handler) updateFiltrationItemImage(ctx *gin.Context) {
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

	id, err := h.services.Category.UpdateFiltrationImage(dto)
	if err != nil {
		newErrorResponse(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	ctx.JSON(http.StatusOK, ItemProcessedResponse{
		Id:      id,
		Message: "filtration item image updated",
	})
}

// deleteFiltrationItemImage godoc
// @Summary Delete filtration image (Minio)
// @Security ApiKeyAuth
// @Tags api/admin
// @Description Deletes filtration item image (Minio)
// @ID delete filtration item image (minio)
// @Accept json
// @Produce json
// @Param id query int true "filtration id"
// @Success 201 {object} string
// @Failure 400 {object} Error
// @Failure 404 {object} Error
// @Failure 500 {object} Error
// @Router /api/admin/filtration-image [delete]
func (h *Handler) deleteFiltrationItemImage(ctx *gin.Context) {
	filtrationId, err := strconv.Atoi(ctx.Query("id"))
	if err != nil {
		newErrorResponse(ctx, http.StatusBadRequest, fmt.Errorf("invalid id, %v", err).Error())
		return
	}

	err = h.services.Category.DeleteFiltrationImage(filtrationId)
	if err != nil {
		newErrorResponse(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	ctx.JSON(http.StatusOK, map[string]any{
		"id":      filtrationId,
		"message": "filtration image (Minio) deleted",
	})
}

// updateFiltrationItem godoc
// @Summary Update filtration
// @Security ApiKeyAuth
// @Tags api/admin
// @Description Updates a filtration item to a category
// @ID update filtration
// @Accept json
// @Produce json
// @Param input body category.Filtration true "filtration info"
// @Success 201 {object} string
// @Failure 400 {object} Error
// @Failure 404 {object} Error
// @Failure 500 {object} Error
// @Router /api/admin/filtration [post]
func (h *Handler) updateFiltrationItem(ctx *gin.Context) {
	var filtration category.UpdateFiltrationDTO
	if err := ctx.BindJSON(&filtration); err != nil {
		newErrorResponse(ctx, http.StatusBadRequest, err.Error())
		return
	}

	id, err := h.services.Category.UpdateFiltration(filtration)
	if err != nil {
		newErrorResponse(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	ctx.JSON(http.StatusOK, map[string]any{
		"id":      id,
		"message": "filtration item updated",
	})
}

// deleteFiltrationItem godoc
// @Summary Delete filtration
// @Security ApiKeyAuth
// @Tags api/admin
// @Description Deletes a filtration item
// @ID delete filtration
// @Accept json
// @Produce json
// @Param id query int true "filtration id"
// @Success 200 {string} string
// @Failure 400 {object} Error
// @Failure 403 {object} Error
// @Failure 500 {object} Error
// @Router /api/admin/filtration [delete]
func (h *Handler) deleteFiltrationItem(ctx *gin.Context) {
	filtrationId := ctx.Query("id")
	filtrationIdInt, err := strconv.Atoi(filtrationId)
	if err != nil {
		newErrorResponse(ctx, http.StatusBadRequest, err.Error())
		return
	}

	err = h.services.Category.DeleteFiltration(filtrationIdInt)
	if err != nil {
		newErrorResponse(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	ctx.JSON(http.StatusOK, "filtration item deleted")
}
