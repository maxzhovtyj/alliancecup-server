package handler

import (
	"github.com/gin-gonic/gin"
)

type Error struct {
	Message string `json:"message"`
}

func newErrorResponse(ctx *gin.Context, statusCode int, message string) {
	ctx.AbortWithStatusJSON(statusCode, Error{message})
}

type ItemProcessedResponse struct {
	Id      int    `json:"id"`
	Message string `json:"message"`
}
