package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	server "github.com/zh0vtyj/allincecup-server"
	"net/http"
)

func (h *Handler) newOrder(ctx *gin.Context) {
	var input server.OrderFullInfo

	if err := ctx.BindJSON(&input); err != nil {
		newErrorResponse(ctx, http.StatusBadRequest, err.Error())
		return
	}

	id, err := getUserId(ctx)
	if err != nil {
		newErrorResponse(ctx, http.StatusInternalServerError, "user role id not found")
		return
	}

	if id != 0 {
		input.Info.UserId = id
	}

	orderId, err := h.services.Orders.New(input)
	if err != nil {
		newErrorResponse(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	ctx.JSON(http.StatusCreated, map[string]interface{}{
		"order_id": orderId,
		"message":  "order created",
	})
}

func (h *Handler) userOrders(ctx *gin.Context) {
	id, err := getUserId(ctx)

	createdAt := ctx.Query("created_at")

	if err != nil {
		newErrorResponse(ctx, http.StatusInternalServerError, "no user's id")
		return
	}

	orders, err := h.services.Orders.GetUserOrders(id, createdAt)
	if err != nil {
		newErrorResponse(ctx, http.StatusInternalServerError, err.Error())
		return
	}
	ctx.JSON(http.StatusOK, map[string]interface{}{
		"orders": orders,
	})
}

func (h *Handler) getOrderById(ctx *gin.Context) {
	orderId := ctx.Query("order_id")

	orderUUID, err := uuid.Parse(orderId)
	if err != nil {
		newErrorResponse(ctx, http.StatusBadRequest, "cannot parse uuid: "+err.Error())
		return
	}

	orderInfo, err := h.services.Orders.GetOrderById(orderUUID)
	if err != nil {
		newErrorResponse(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	ctx.JSON(http.StatusOK, map[string]interface{}{
		"data": orderInfo,
	})
}

// adminGetOrders godoc
// @Security 	 ApiKeyAuth
// @Summary      Get Orders
// @Tags         api/admin
// @Description  get orders by status
// @ID get orders
// @Accept       json
// @Produce      json
// @Param created_at query string false "Last item created at for pagination"
// @Param order_status query string true "Sort by orders status"
// @Success      200  {array} server.Order
// @Failure      400  {object}  Error
// @Failure      404  {object}  Error
// @Failure      500  {object}  Error
// @Router       /api/admin/all-orders [get]
func (h *Handler) adminGetOrders(ctx *gin.Context) {
	createdAt := ctx.Query("created_at")
	orderStatus := ctx.Query("order_status")

	if orderStatus != statusCompleted && orderStatus != statusProcessed && orderStatus != statusInProgress {
		newErrorResponse(ctx, http.StatusBadRequest, "invalid order status")
		return
	}

	orders, err := h.services.Orders.GetAdminOrders(orderStatus, createdAt)
	if err != nil {
		newErrorResponse(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	ctx.JSON(http.StatusOK, map[string]interface{}{
		"data": orders,
	})
}
