package handler

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/zh0vtyj/allincecup-server/internal/domain/order"
	"net/http"
	"os"
	"strconv"
)

type OrderResponse struct {
	Id      int    `json:"id"`
	Message string `json:"message"`
}

type ProcessedOrderStatus struct {
	OrderId  int    `json:"orderId" binding:"required"`
	ToStatus string `json:"toStatus" binding:"required"`
}

// newOrder godoc
// @Summary      NewOrder
// @Tags         api
// @Description  creates a new order
// @ID creates an order
// @Accept       json
// @Produce      json
// @Param        input body order.CreateDTO true "order info"
// @Success      200  {object}  handler.OrderResponse
// @Failure      400  {object}  Error
// @Failure      500  {object}  Error
// @Router       /api/order [post]
func (h *Handler) newOrder(ctx *gin.Context) {
	var input order.CreateDTO

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
		input.Order.UserId = &id
	}

	cartId, err := getCartId(ctx)
	if err != nil {
		newErrorResponse(ctx, http.StatusBadRequest, "user cart id not found")
		return
	}

	orderId, err := h.services.Order.New(input, cartId)
	if err != nil {
		newErrorResponse(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	ctx.SetCookie(cartId, "", -1, "/", h.cfg.Domain, false, true)

	ctx.JSON(http.StatusCreated, OrderResponse{
		Id:      orderId,
		Message: "order created",
	})
}

// adminNewOrder godoc
// @Summary      New order
// @Tags         api
// @Description  Execute new order by admin
// @ID creates an order by admin
// @Accept       json
// @Produce      json
// @Param        input body order.CreateDTO true "order info"
// @Success      200  {object}  handler.OrderResponse
// @Failure      400  {object}  Error
// @Failure      500  {object}  Error
// @Router       /api/order [post]
func (h *Handler) adminNewOrder(ctx *gin.Context) {
	var input order.CreateDTO

	if err := ctx.BindJSON(&input); err != nil {
		newErrorResponse(ctx, http.StatusBadRequest, err.Error())
		return
	}

	id, err := getUserId(ctx)
	if err != nil {
		newErrorResponse(ctx, http.StatusInternalServerError, "user id not found")
		return
	}
	input.Order.ExecutedBy = &id

	orderId, err := h.services.Order.AdminNew(input)
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
// @Success      200  {array}  	order.SelectDTO
// @Failure      401  {object}  Error
// @Failure      500  {object}  Error
// @Router       /api/client/user-order [get]
func (h *Handler) userOrders(ctx *gin.Context) {
	id, err := getUserId(ctx)
	if err != nil {
		newErrorResponse(ctx, http.StatusUnauthorized, "user id not found")
		return
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

	ctx.JSON(http.StatusOK, orders)
}

// getOrderById godoc
// @Summary      GetOrderById
// @Security 	 ApiKeyAuth
// @Tags         api/client
// @Description  gets user order full info by its id
// @ID gets order by id
// @Produce      json
// @Param        id query string true "order id"
// @Success      200  {array}  order.SelectDTO "get order"
// @Failure      400  {object}  Error
// @Failure      401  {object}  Error
// @Failure      500  {object}  Error
// @Router       /api/admin/order [get]
func (h *Handler) getOrderById(ctx *gin.Context) {
	orderId, err := strconv.Atoi(ctx.Query("id"))
	if err != nil {
		newErrorResponse(ctx, http.StatusBadRequest, err.Error())
		return
	}

	orderInfo, err := h.services.Order.GetOrderById(orderId)
	if err != nil {
		newErrorResponse(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	ctx.JSON(http.StatusOK, orderInfo)
}

// adminGetOrders godoc
// @Security ApiKeyAuth
// @Summary Get Orders
// @Tags api/admin
// @Description get order by status
// @ID get order
// @Produce json
// @Param created_at query string false "Last item created at for pagination"
// @Param order_status query string true "Sort by order status"
// @Success 200 {array} order.Order
// @Failure 400 {object} Error
// @Failure 404 {object} Error
// @Failure 500 {object} Error
// @Router /api/admin/orders [get]
func (h *Handler) adminGetOrders(ctx *gin.Context) {
	createdAt := ctx.Query("created_at")
	orderStatus := ctx.Query("order_status")
	search := ctx.Query("search")

	if orderStatus != "" {
		if orderStatus != StatusCompleted && orderStatus != StatusProcessed && orderStatus != StatusInProgress {
			newErrorResponse(ctx, http.StatusBadRequest, "order status either empty or invalid")
			return
		}
	}

	orders, err := h.services.Order.AdminGetOrders(orderStatus, createdAt, search)
	if err != nil {
		newErrorResponse(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	ctx.JSON(http.StatusOK, orders)
}

// getDeliveryPaymentTypes godoc
// @Summary      Get order info types
// @Tags         api
// @Description  get payment and delivery types
// @ID get order info types
// @Produce      json
// @Success      200  {array}   shopping.DeliveryPaymentTypes
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

// processedOrder godoc
// @Summary      Processed order by id
// @Security 	 ApiKeyAuth
// @Tags         api/admin
// @Description  handler for admin/moderator to processed order by id
// @ID processed order
// @Accept 	  	 json
// @Produce      json
// @Param 		 input body handler.ProcessedOrderStatus true "order status"
// @Success      200  {object}  ItemProcessedResponse
// @Failure      400  {object}  Error
// @Failure      500  {object}  Error
// @Router       /api/admin/processed-order [put]
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

	err := h.services.Order.ProcessedOrder(orderInput.OrderId, StatusProcessed)
	if err != nil {
		newErrorResponse(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	ctx.JSON(http.StatusOK, ItemProcessedResponse{
		Id:      orderInput.OrderId,
		Message: "order processed",
	})
}

// getOrderInvoice godoc
// @Summary      Get pdf invoice by order id
// @Security 	 ApiKeyAuth
// @Tags         api
// @Description  Get pdf file with order invoice
// @ID pdf order invoice
// @Produce      json
// @Param 		 id query int true "order id"
// @Success      200  {attachment} File
// @Failure      400  {object}  Error
// @Failure      500  {object}  Error
// @Router       /api/invoice [put]
func (h *Handler) getOrderInvoice(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Query("id"))
	if err != nil {
		newErrorResponse(ctx, http.StatusBadRequest, err.Error())
		return
	}

	pwd, err := os.Getwd()
	if err != nil {
		newErrorResponse(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	fileName := fmt.Sprintf("invoice-%d.pdf", id)
	filePath := fmt.Sprintf("%s/tmp/%s", pwd, fileName)

	invoice, err := h.services.Order.GetInvoice(id)
	if err != nil {
		newErrorResponse(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	err = invoice.OutputFileAndClose(filePath)
	if err != nil {
		newErrorResponse(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	ContentDispositionValue := fmt.Sprintf("attachment; filename=%s", fileName)
	ctx.Header("Content-Disposition", ContentDispositionValue)
	//ctx.Header("Content-Type", "application/octet-stream")
	ctx.FileAttachment(filePath, fileName)

	err = os.Remove(filePath)
	if err != nil {
		return
	}
}
