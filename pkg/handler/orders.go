package handler

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	server "github.com/zh0vtyj/allincecup-server"
	"net/http"
)

type OrderResponse struct {
	Id      uuid.UUID `json:"id"`
	Message string    `json:"message"`
}

// newOrder godoc
// @Summary      NewOrder
// @Tags         api
// @Description  creates a new order
// @ID creates an order
// @Accept       json
// @Produce      json
// @Param        input body server.OrderFullInfo true "order info"
// @Success      200  {object}  handler.OrderResponse
// @Failure      400  {object}  Error
// @Failure      404  {object}  Error
// @Failure      500  {object}  Error
// @Router       /api/new-order [post]
func (h *Handler) newOrder(ctx *gin.Context) {
	var input server.OrderFullInfo

	if err := ctx.BindJSON(&input); err != nil {
		newErrorResponse(ctx, http.StatusBadRequest, err.Error())
		return
	}

	id, err := getUserId(ctx)
	if err != nil {
		newErrorResponse(ctx, http.StatusInternalServerError, "user id not found")
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

	ctx.JSON(http.StatusCreated, OrderResponse{
		Id:      orderId,
		Message: "order created",
	})
}

// userOrders godoc
// @Summary      GetUserOrders
// @Security 	 ApiKeyAuth
// @Tags         api/client
// @Description  gets user orders
// @ID gets orders
// @Accept       json
// @Produce      json
// @Param        created_at query string false "last order created_at for pagination"
// @Success      200  {object}  server.OrderInfo
// @Failure      400  {object}  Error
// @Failure      404  {object}  Error
// @Failure      500  {object}  Error
// @Router       /api/client/user-orders [get]
func (h *Handler) userOrders(ctx *gin.Context) {
	id, err := getUserId(ctx)
	if err != nil {
		newErrorResponse(ctx, http.StatusUnauthorized, "user id not found")
	}

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

// getOrderById godoc
// @Summary      GetOrderById
// @Security 	 ApiKeyAuth
// @Tags         api/client
// @Description  gets user order full info by its id
// @ID gets order by id
// @Produce      json
// @Param        order_id query string true "order id"
// @Success      200  {object}  server.OrderInfo
// @Failure      400  {object}  Error
// @Failure      404  {object}  Error
// @Failure      500  {object}  Error
// @Router       /api/client/get-order [get]
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

// getDeliveryPaymentTypes godoc
// @Summary      Get order info types
// @Tags         api/order-info-types
// @Description  get payment and delivery types
// @ID get order info types
// @Produce      json
// @Success      200  {array}   server.DeliveryPaymentTypes
// @Failure      400  {object}  Error
// @Failure      404  {object}  Error
// @Failure      500  {object}  Error
// @Router       /api/order-info-types [get]
func (h *Handler) deliveryPaymentTypes(ctx *gin.Context) {
	deliveryTypes, err := h.services.Orders.DeliveryPaymentTypes()
	if err != nil {
		newErrorResponse(ctx, http.StatusInternalServerError, fmt.Errorf("failed to get delivery types due to: %v", err).Error())
		return
	}

	ctx.JSON(http.StatusOK, deliveryTypes)
}
