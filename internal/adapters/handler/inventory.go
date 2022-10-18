package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/zh0vtyj/allincecup-server/internal/domain/inventory"
	"net/http"
	"strconv"
	"time"
)

// getProductsToInventory godoc
// @Summary      Get products to inventory them
// @Security 	 ApiKeyAuth
// @Tags         api/admin/super
// @Product  	 Gets all products to inventory
// @ID 			 gets products to inventory
// @Produce      json
// @Success      201  {array}  inventory.CurrentProductDTO
// @Failure      400  {object}  Error
// @Failure      404  {object}  Error
// @Failure      500  {object}  Error
// @Router       /api/admin/super/inventory [get]
func (h *Handler) getProductsToInventory(ctx *gin.Context) {
	products, err := h.services.Inventory.Products()
	if err != nil {
		newErrorResponse(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	ctx.JSON(http.StatusOK, products)
}

// doInventory godoc
// @Summary      Do products inventory
// @Security 	 ApiKeyAuth
// @Tags         api/admin/super
// @Product  	 Do products inventory
// @ID 			 do products inventory
// @Accept       json
// @Produce      json
// @Param		 input body []inventory.InsertProductDTO true "products info to inventory"
// @Success      201  {object}  object
// @Failure      400  {object}  Error
// @Failure      500  {object}  Error
// @Router       /api/admin/super/inventory [post]
func (h *Handler) doInventory(ctx *gin.Context) {
	var input []inventory.InsertProductDTO

	if err := ctx.BindJSON(&input); err != nil {
		newErrorResponse(ctx, http.StatusBadRequest, err.Error())
		return
	}

	err := h.services.Inventory.New(input)
	if err != nil {
		newErrorResponse(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	ctx.JSON(http.StatusCreated, map[string]any{
		"message": "inventory created",
	})
}

// getInventories godoc
// @Summary      Get inventories
// @Security 	 ApiKeyAuth
// @Tags         api/admin/super
// @Product  	 Gets inventories
// @ID 			 gets inventories
// @Produce      json
// @Param        createdAt query string false "Last inventory created at for pagination"
// @Success      201  {array}  inventory.CurrentProductDTO
// @Failure      400  {object}  Error
// @Failure      500  {object}  Error
// @Router       /api/admin/super/inventories [get]
func (h *Handler) getInventories(ctx *gin.Context) {
	createdAt := ctx.Query("createdAt")

	if createdAt != "" {
		_, err := time.Parse("2022-10-14T13:53:56.383658Z", createdAt)
		if err != nil {
			newErrorResponse(ctx, http.StatusBadRequest, err.Error())
			return
		}
	}

	inventories, err := h.services.Inventory.GetAll(createdAt)
	if err != nil {
		newErrorResponse(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	ctx.JSON(http.StatusOK, inventories)
}

// getInventoryProducts godoc
// @Summary      Get inventory products
// @Security 	 ApiKeyAuth
// @Tags         api/admin/super
// @Product  	 Gets inventory products
// @ID 			 gets inventory products
// @Produce      json
// @Param        id query int true "Inventory id"
// @Success      201  {array}   inventory.SelectProductDTO
// @Failure      400  {object}  Error
// @Failure      500  {object}  Error
// @Router       /api/admin/super/inventory-products [get]
func (h *Handler) getInventoryProducts(ctx *gin.Context) {
	inventoryId, err := strconv.Atoi(ctx.Query("id"))
	if err != nil {
		newErrorResponse(ctx, http.StatusBadRequest, "invalid inventory id")
		return
	}

	products, err := h.services.Inventory.GetInventoryProducts(inventoryId)
	if err != nil {
		newErrorResponse(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	ctx.JSON(http.StatusOK, products)
}
