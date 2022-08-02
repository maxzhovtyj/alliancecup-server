package handler

import (
	"github.com/gin-gonic/gin"
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
	id := 1

	createdAt := ctx.Query("created_at")

	//if err != nil {
	//	newErrorResponse(ctx, http.StatusInternalServerError, "no user's id")
	//	return
	//}

	orders, err := h.services.Orders.GetUserOrders(id, createdAt)
	if err != nil {
		newErrorResponse(ctx, http.StatusInternalServerError, err.Error())
		return
	}
	ctx.JSON(http.StatusOK, map[string]interface{}{
		"orders": orders,
	})
}
