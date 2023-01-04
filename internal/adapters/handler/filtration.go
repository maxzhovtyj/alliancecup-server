package handler

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/zh0vtyj/allincecup-server/internal/domain/category"
	"github.com/zh0vtyj/allincecup-server/internal/domain/models"
	"net/http"
	"strconv"
)

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
	err := ctx.Request.ParseMultipartForm(32 << 20)
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

	files, ok := ctx.Request.MultipartForm.File["file"]
	if len(files) != 0 {
		if !ok {
			newErrorResponse(ctx, http.StatusBadRequest, "something wrong with file you provided")
			return
		}

		fileInfo := files[0]
		fileReader, err := fileInfo.Open()
		if err != nil {
			newErrorResponse(ctx, http.StatusBadRequest, err.Error())
			return
		}

		filtration.Img = &models.FileDTO{
			Name:   fileInfo.Filename,
			Size:   fileInfo.Size,
			Reader: fileReader,
		}
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
