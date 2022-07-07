package handler

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

const (
	authorizationHeader = "Authorization"
	userCtx             = "userId"
	userRoleIdCtx       = "userRoleId"
)

func (h *Handler) userIdentity(ctx *gin.Context) {
	header := ctx.GetHeader(authorizationHeader)
	if header == "" {
		return
	}

	headerParts := strings.Split(header, " ")
	if len(headerParts) != 2 {
		return
	}

	//parse token
	userId, userRoleId, err := h.services.Authorization.ParseToken(headerParts[1])
	if err != nil {
		newErrorResponse(ctx, http.StatusUnauthorized, err.Error())
		return
	}

	ctx.Set(userCtx, userId)
	ctx.Set(userRoleIdCtx, userRoleId)
}

func (h *Handler) userAuthorized(ctx *gin.Context) {
	header := ctx.GetHeader(authorizationHeader)
	if header == "" {
		newErrorResponse(ctx, http.StatusUnauthorized, "empty auth header")
		return
	}

	headerParts := strings.Split(header, " ")
	if len(headerParts) != 2 || headerParts[0] != "Bearer" {
		newErrorResponse(ctx, http.StatusUnauthorized, "invalid auth header")
		return
	}

	if len(headerParts[1]) == 0 {
		newErrorResponse(ctx, http.StatusUnauthorized, "token is empty")
		return
	}

	//parse token
	_, _, err := h.services.Authorization.ParseToken(headerParts[1])
	if err != nil {
		newErrorResponse(ctx, http.StatusUnauthorized, err.Error())
		return
	}
}

func (h *Handler) userHasPermission(ctx *gin.Context) {
	userRoleId, ok := ctx.Get(userRoleIdCtx)
	if !ok {
		newErrorResponse(ctx, http.StatusInternalServerError, "user not found")
		return
	}

	userRoleIdInt, ok := userRoleId.(int)
	if !ok {
		newErrorResponse(ctx, http.StatusInternalServerError, "user id is of invalid type")
		return
	}

	if userRoleIdInt == 1 {
		newErrorResponse(ctx, http.StatusForbidden, "access forbidden")
	}
}
