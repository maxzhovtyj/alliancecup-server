package handler

import (
	"errors"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

const (
	authorizationHeader = "Authorization"
	userCartCookie      = "UserCart"
	userCtx             = "userId"
	userRoleIdCtx       = "userRoleId"
	userCartCtx         = "userCart"
)

func (h *Handler) userIdentity(ctx *gin.Context) {
	header := ctx.GetHeader(authorizationHeader)
	if header == "" {
		ctx.Set(userCtx, 0)
		ctx.Set(userRoleIdCtx, 0)
		return
	}

	headerParts := strings.Split(header, " ")
	if len(headerParts) != 2 || headerParts[0] != "Bearer" {
		newErrorResponse(ctx, http.StatusBadRequest, "invalid token")
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

func (h *Handler) userHasPermission(ctx *gin.Context) {
	userRoleIdInt, err := getUserRoleId(ctx)
	if err != nil {
		newErrorResponse(ctx, http.StatusInternalServerError, "user role not found or it's wrong type")
		return
	}

	if userRoleIdInt == 0 || userRoleIdInt == 1 {
		newErrorResponse(ctx, http.StatusForbidden, "access forbidden")
		return
	}
}

func (h *Handler) userAuthorized(ctx *gin.Context) {
	userRoleId, err := getUserRoleId(ctx)
	if err != nil {
		newErrorResponse(ctx, http.StatusInternalServerError, "no user role or it's wrong type: "+err.Error())
		return
	}

	if userRoleId == 0 {
		newErrorResponse(ctx, http.StatusUnauthorized, "user unauthorized")
		return
	}
}

func (h *Handler) superAdmin(ctx *gin.Context) {
	userRoleId, err := getUserRoleId(ctx)
	if err != nil {
		newErrorResponse(ctx, http.StatusInternalServerError, "no user role or it's wrong type: "+err.Error())
		return
	}

	if userRoleId != 3 {
		newErrorResponse(ctx, http.StatusUnauthorized, "user unauthorized")
		return
	}
}

func getUserId(ctx *gin.Context) (int, error) {
	id, ok := ctx.Get(userCtx)
	if !ok {
		newErrorResponse(ctx, http.StatusInternalServerError, "userId not found")
		return 0, errors.New("user id not found")
	}

	idInt, ok := id.(int)
	if !ok {
		newErrorResponse(ctx, http.StatusInternalServerError, "user id is not of type int")
		return 0, errors.New("user id is not of type int")
	}

	return idInt, nil
}

func getUserRoleId(ctx *gin.Context) (int, error) {
	id, ok := ctx.Get(userRoleIdCtx)
	if !ok {
		newErrorResponse(ctx, http.StatusInternalServerError, "role id not found")
		return 0, errors.New("role id not found")
	}

	idInt, ok := id.(int)
	if !ok {
		newErrorResponse(ctx, http.StatusInternalServerError, "user's role id is not of type int")
		return 0, errors.New("user's role id is not of type int")
	}

	return idInt, nil
}

func (h *Handler) getShoppingInfo(ctx *gin.Context) {
	userCartId, err := ctx.Cookie(userCartCookie)
	if err != nil {
		cartId, errNewCart := h.newCart(ctx)
		if errNewCart != nil {
			newErrorResponse(ctx, http.StatusInternalServerError, err.Error())
			return
		}

		ctx.Set(userCartCtx, cartId)
	} else {
		// TODO check whether such cart uuid exist in redis cache
		ctx.Set(userCartCtx, userCartId)
	}
}

func (h *Handler) newCart(ctx *gin.Context) (string, error) {
	id, err := getUserId(ctx)
	if err != nil {
		return "", err
	}

	cartUUID, err := h.services.Shopping.NewCart(id)
	if err != nil {
		newErrorResponse(ctx, http.StatusInternalServerError, err.Error())
		return "", err
	}

	return cartUUID.String(), err
}
