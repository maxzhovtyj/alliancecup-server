package handler

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

const (
	authorizationHeader     = "Authorization"
	userIdCtx               = "userId"
	userRoleCodeCtx         = "userRoleCode"
	userCartCookie          = "UserCart"
	userFavouritesCookie    = "UserFavourites"
	userCartCtx             = "userCart"
	userFavouritesCtx       = "userFavourites"
	userCartCookieTTL       = 60 * 60 * 72
	userFavouritesCookieTTL = 60 * 60 * 72
)

func (h *Handler) userIdentity(ctx *gin.Context) {
	header := ctx.GetHeader(authorizationHeader)
	if header == "" {
		ctx.Set(userIdCtx, 0)
		ctx.Set(userRoleCodeCtx, h.cfg.Roles.Guest)
		return
	}

	headerParts := strings.Split(header, " ")
	if len(headerParts) != 2 || headerParts[0] != "Bearer" {
		newErrorResponse(ctx, http.StatusBadRequest, "invalid token")
		return
	}

	// parse token
	userId, userRoleCode, err := h.services.Authorization.ParseToken(headerParts[1])
	if err != nil {
		newErrorResponse(ctx, http.StatusUnauthorized, err.Error())
		return
	}

	ctx.Set(userIdCtx, userId)
	ctx.Set(userRoleCodeCtx, userRoleCode)
}

func (h *Handler) moderatorPermission(ctx *gin.Context) {
	userRoleCode, err := getUserRole(ctx)
	if err != nil {
		newErrorResponse(ctx, http.StatusInternalServerError, "user role not found or it's wrong type")
		return
	}

	if userRoleCode != h.cfg.Roles.Moderator && userRoleCode != h.cfg.Roles.SuperAdmin {
		newErrorResponse(ctx, http.StatusForbidden, "access forbidden")
		return
	}
}

func (h *Handler) userAuthorized(ctx *gin.Context) {
	userRoleId, err := getUserRole(ctx)
	if err != nil {
		newErrorResponse(ctx, http.StatusInternalServerError, "no user role or it's wrong type: "+err.Error())
		return
	}

	if userRoleId == h.cfg.Roles.Guest {
		newErrorResponse(ctx, http.StatusUnauthorized, "user unauthorized")
		return
	}
}

func (h *Handler) superAdmin(ctx *gin.Context) {
	userRoleCode, err := getUserRole(ctx)
	if err != nil {
		newErrorResponse(ctx, http.StatusInternalServerError, "no user role or it's wrong type: "+err.Error())
		return
	}

	if userRoleCode != h.cfg.Roles.SuperAdmin {
		newErrorResponse(ctx, http.StatusForbidden, "access forbidden")
		return
	}
}

func getUserId(ctx *gin.Context) (int, error) {
	id, ok := ctx.Get(userIdCtx)
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

func getUserRole(ctx *gin.Context) (string, error) {
	roleCode, ok := ctx.Get(userRoleCodeCtx)
	if !ok {
		return "", errors.New("role id not found")
	}

	roleCodeStr, ok := roleCode.(string)
	if !ok {
		return "", errors.New("user's role id is not of type int")
	}

	return roleCodeStr, nil
}

func (h *Handler) getShoppingInfo(ctx *gin.Context) {
	userCartId, err := ctx.Cookie(userCartCookie)

	if err != nil || userCartId == "" {
		cartId, errNewCart := h.newCart(ctx)
		if errNewCart != nil {
			newErrorResponse(ctx, http.StatusInternalServerError, errNewCart.Error())
			return
		}

		ctx.SetCookie(userCartCookie, cartId, userCartCookieTTL, "/", h.cfg.Domain, false, true)
		ctx.Set(userCartCtx, cartId)
	} else {
		ctx.Set(userCartCtx, userCartId)
	}

	//userFavouritesId, favErr := ctx.Cookie(userFavouritesCookie)
	//if favErr != nil || userFavouritesId == "" {
	//	favouritesId, errNewFav := h.newFavourites(ctx)
	//	if errNewFav != nil {
	//		newErrorResponse(ctx, http.StatusInternalServerError, errNewFav.Error())
	//		return
	//	}
	//
	//	ctx.SetCookie(
	//		userFavouritesCookie,
	//		favouritesId,
	//		userFavouritesCookieTTL,
	//		"/",
	//		h.cfg.Domain,
	//		false,
	//		true,
	//	)
	//	ctx.Set(userFavouritesCtx, favouritesId)
	//} else {
	//	ctx.Set(userFavouritesCtx, userFavouritesId)
	//}
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

//func (h *Handler) newFavourites(ctx *gin.Context) (string, error) {
//	id, err := getUserId(ctx)
//	if err != nil {
//		return "", err
//	}
//
//	favouritesUUID, err := h.services.Shopping.NewFavourites(id)
//	if err != nil {
//		return "", err
//	}
//
//	return favouritesUUID.String(), err
//}

func getCartId(ctx *gin.Context) (string, error) {
	id, exists := ctx.Get(userCartCtx)
	if !exists {
		return "", fmt.Errorf("failed to find user cart id")
	}

	return id.(string), nil
}

//func getFavouritesId(ctx *gin.Context) (string, error) {
//	id, exists := ctx.Get(userFavouritesCtx)
//	if !exists {
//		return "", fmt.Errorf("failed to find user favourites id")
//	}
//
//	return id.(string), nil
//}
