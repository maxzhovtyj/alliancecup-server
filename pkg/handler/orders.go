package handler

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	server "github.com/zh0vtyj/allincecup-server/internal/order"
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
// @Param        input body server.Order true "order info"
// @Success      200  {object}  handler.OrderResponse
// @Failure      400  {object}  Error
// @Failure      404  {object}  Error
// @Failure      500  {object}  Error
// @Router       /api/new-order [post]
func (h *Handler) newOrder(ctx *gin.Context) {
	var input server.Info

	if err := ctx.BindJSON(&input); err != nil {
		newErrorResponse(ctx, http.StatusBadRequest, err.Error())
		return
	}

	if input.Order.OrderSumPrice < 400 {
		newErrorResponse(ctx, http.StatusBadRequest, fmt.Errorf("failed to create order, minimal order price is 400hrn").Error())
	}

	id, err := getUserId(ctx)
	if err != nil {
		newErrorResponse(ctx, http.StatusInternalServerError, "user id not found")
		return
	}

	if id != 0 {
		input.Order.UserId = id
	}

	orderId, err := h.services.Order.New(input)
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
// @Description  gets user order
// @ID gets order
// @Accept       json
// @Produce      json
// @Param        created_at query string false "last order created_at for pagination"
// @Success      200  {object}  server.FullInfo
// @Failure      400  {object}  Error
// @Failure      404  {object}  Error
// @Failure      500  {object}  Error
// @Router       /api/client/user-order [get]
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

	orders, err := h.services.Order.GetUserOrders(id, createdAt)
	if err != nil {
		newErrorResponse(ctx, http.StatusInternalServerError, err.Error())
		return
	}
	ctx.JSON(http.StatusOK, map[string]interface{}{
		"order": orders,
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
// @Success      200  {object}  server.FullInfo
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

	orderInfo, err := h.services.Order.GetOrderById(orderUUID)
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
// @Description  get order by status
// @ID get order
// @Produce      json
// @Param created_at query string false "Last item created at for pagination"
// @Param order_status query string true "Sort by order status"
// @Success      200  {array} server.Order
// @Failure      400  {object}  Error
// @Failure      404  {object}  Error
// @Failure      500  {object}  Error
// @Router       /api/admin/all-order [get]
func (h *Handler) adminGetOrders(ctx *gin.Context) {
	createdAt := ctx.Query("created_at")
	orderStatus := ctx.Query("order_status")

	if orderStatus != StatusCompleted && orderStatus != StatusProcessed && orderStatus != StatusInProgress {
		orderStatus = ""
		fmt.Println("order status either empty or invalid")
	}

	orders, err := h.services.Order.GetAdminOrders(orderStatus, createdAt)
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
// @Tags         api
// @Description  get payment and delivery types
// @ID get order info types
// @Produce      json
// @Success      200  {array}   server.DeliveryPaymentTypes
// @Failure      400  {object}  Error
// @Failure      404  {object}  Error
// @Failure      500  {object}  Error
// @Router       /api/order-info-types [get]
func (h *Handler) deliveryPaymentTypes(ctx *gin.Context) {
	deliveryTypes, err := h.services.Order.DeliveryPaymentTypes()
	if err != nil {
		newErrorResponse(ctx, http.StatusInternalServerError, fmt.Errorf("failed to get delivery types due to: %v", err).Error())
		return
	}

	ctx.JSON(http.StatusOK, deliveryTypes)
}

type ProcessedOrderStatus struct {
	OrderId  uuid.UUID `json:"orderId" binding:"required"`
	ToStatus string    `json:"toStatus" binding:"required"`
}

// processedOrder godoc
// @Summary      Processed order by id
// @Tags         api/admin
// @Description  handler for admin/moderator to processed order by id
// @ID processed order
// @Accept 	  	 json
// @Produce      json
// @Param 		 input body handler.ProcessedOrderStatus true "order status"
// @Success      200  {array}   server.DeliveryPaymentTypes
// @Failure      400  {object}  Error
// @Failure      404  {object}  Error
// @Failure      500  {object}  Error
// @Router       /api/order-info-types [get]
func (h *Handler) processedOrder(ctx *gin.Context) {
	var orderInput ProcessedOrderStatus

	if err := ctx.BindJSON(&orderInput); err != nil {
		newErrorResponse(ctx, http.StatusBadRequest, fmt.Errorf("failed to bind input data due to: %v", err).Error())
		return
	}

	if orderInput.ToStatus != StatusCompleted &&
		orderInput.ToStatus != StatusProcessed &&
		orderInput.ToStatus != StatusInProgress {
		newErrorResponse(ctx, http.StatusBadRequest, fmt.Errorf("failed to update order status, invalid input status").Error())
		return
	}

	err := h.services.Order.ProcessedOrder(orderInput.OrderId, orderInput.ToStatus)
	if err != nil {
		newErrorResponse(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	ctx.JSON(http.StatusOK, map[string]any{
		"orderId": orderInput.OrderId,
		"message": "order status successfully updated",
	})
}
