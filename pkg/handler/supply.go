package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/zh0vtyj/allincecup-server/pkg/models"
	"net/http"
)

func (h *Handler) getAllSupply(ctx *gin.Context) {
	createdAt := ctx.Query("createAt")

	supply, err := h.services.Supply.GetAll(createdAt)
	if err != nil {
		newErrorResponse(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	ctx.JSON(http.StatusOK, supply)
}

func (h *Handler) newSupply(ctx *gin.Context) {
	var input models.SupplyDTO

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

func (h *Handler) updateSupply(ctx *gin.Context) {
	panic("implement me")
}

func (h *Handler) deleteSupply(ctx *gin.Context) {
	panic("implement me")
}
