package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/zh0vtyj/alliancecup-server/internal/domain/supply"
	"net/http"
	"strconv"
)

// getAllSupply godoc
// @Summary      Get Supplies
// @Security 	 ApiKeyAuth
// @Tags         api/admin
// @Description  Get supplies
// @ID 			 get supplies
// @Produce      json
// @Param        createdAt query string false "Last item createdAt for pagination"
// @Success      200  {array}  supply.InfoDTO
// @Failure      400  {object}  Error
// @Failure      404  {object}  Error
// @Failure      500  {object}  Error
// @Router       /api/admin/supply [get]
func (h *Handler) getAllSupply(ctx *gin.Context) {
	createdAt := ctx.Query("createdAt")

	supplies, err := h.services.Supply.GetAll(createdAt)
	if err != nil {
		newErrorResponse(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	ctx.JSON(http.StatusOK, supplies)
}

// getSupplyProducts godoc
// @Summary      Get Supply Products
// @Security 	 ApiKeyAuth
// @Tags         api/admin
// @Description  Get supply products
// @ID 			 get supply products
// @Produce      json
// @Param        id query int true "Supply id"
// @Param        createdAt query string false "Last item createdAt for pagination"
// @Success      200  {array}   supply.ProductDTO
// @Failure      400  {object}  Error
// @Failure      500  {object}  Error
// @Router       /api/admin/supply-products [get]
func (h *Handler) getSupplyProducts(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Query("id"))
	if err != nil {
		newErrorResponse(ctx, http.StatusBadRequest, err.Error())
		return
	}

	createdAt := ctx.Query("createdAt")

	products, err := h.services.Supply.Products(id, createdAt)
	if err != nil {
		newErrorResponse(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	ctx.JSON(http.StatusOK, products)
}

// newSupply godoc
// @Summary      Create new supply
// @Security 	 ApiKeyAuth
// @Tags         api/admin
// @Description  Creates new supply
// @ID 			 creates new supply
// @Accept       json
// @Produce      json
// @Param        supply body supply.Supply true "Supply info"
// @Success      200  {object}  object
// @Failure      400  {object}  Error
// @Failure      500  {object}  Error
// @Router       /api/admin/supply [post]
func (h *Handler) newSupply(ctx *gin.Context) {
	var input supply.Supply

	if err := ctx.BindJSON(&input); err != nil {
		newErrorResponse(ctx, http.StatusBadRequest, err.Error())
		return
	}

	err := h.services.Supply.New(input)
	if err != nil {
		newErrorResponse(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	ctx.JSON(http.StatusOK, map[string]any{
		"message": "supply added",
	})
}

// deleteSupply  godoc
// @Summary      Delete supply by id
// @Security 	 ApiKeyAuth
// @Tags         api/admin
// @Description  Deletes supply
// @ID 			 deletes supply
// @Produce      json
// @Param        id query int true "Supply id"
// @Success      200  {object}  object
// @Failure      400  {object}  Error
// @Failure      500  {object}  Error
// @Router       /api/admin/supply [delete]
func (h *Handler) deleteSupply(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Query("id"))
	if err != nil {
		newErrorResponse(ctx, http.StatusBadRequest, err.Error())
		return
	}

	err = h.services.Supply.Delete(id)
	if err != nil {
		newErrorResponse(ctx, http.StatusBadRequest, err.Error())
		return
	}

	ctx.JSON(http.StatusOK, ItemProcessedResponse{
		Id:      id,
		Message: "supply successfully deleted",
	})
}
