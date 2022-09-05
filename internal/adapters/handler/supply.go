package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/zh0vtyj/allincecup-server/internal/domain/supply"
	"net/http"
	"strconv"
)

func (h *Handler) getAllSupply(ctx *gin.Context) {
	createdAt := ctx.Query("createAt")

	s, err := h.services.Supply.GetAll(createdAt)
	if err != nil {
		newErrorResponse(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	ctx.JSON(http.StatusOK, s)
}

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

	ctx.JSON(http.StatusOK, map[string]any{
		"id":      id,
		"message": "deleted",
	})
}
