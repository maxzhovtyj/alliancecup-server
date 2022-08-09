package handler

import (
	"github.com/gin-gonic/gin"
	server "github.com/zh0vtyj/allincecup-server"
	_ "github.com/zh0vtyj/allincecup-server/docs"
	"net/http"
	"net/mail"
	"time"
)

const refreshTokenTTL = 1440 * time.Hour

// signUp godoc
// @Summary      SignUp
// @Tags         auth
// @Description  registers a new user
// @ID create account
// @Accept       json
// @Produce      json
// @Param        input body server.User true "account info"
// @Success      200  {integer} integer 2
// @Failure      400  {object}  Error
// @Failure      404  {object}  Error
// @Failure      500  {object}  Error
// @Router       /auth/sign-up [post]
func (h *Handler) signUp(ctx *gin.Context) {
	var input server.User

	if err := ctx.BindJSON(&input); err != nil {
		newErrorResponse(ctx, http.StatusBadRequest, err.Error())
		return
	}

	// email, password, phone_number validation
	_, err := mail.ParseAddress(input.Email)
	if err != nil {
		newErrorResponse(ctx, http.StatusBadRequest, "non valid email")
		return
	}
	if len(input.Password) < 4 {
		newErrorResponse(ctx, http.StatusBadRequest, "non valid password")
		return
	}
	if len(input.PhoneNumber) < 10 {
		newErrorResponse(ctx, http.StatusBadRequest, "non valid phone_number")
		return
	}

	id, roleId, err := h.services.Authorization.CreateUser(input)
	if err != nil {
		newErrorResponse(ctx, http.StatusInternalServerError, err.Error())
		return
	}
	ctx.JSON(http.StatusCreated, map[string]interface{}{
		"id":      id,
		"role_id": roleId,
	})
}

// createModerator godoc
// @Summary      CreateModerator
// @Security ApiKeyAuth
// @Tags         api/admin
// @Description  registers a new moderator
// @ID create account for moderator
// @Accept       json
// @Produce      json
// @Param        input body server.User true "account info"
// @Success      200  {integer} integer 2
// @Failure      400  {object}  Error
// @Failure      404  {object}  Error
// @Failure      500  {object}  Error
// @Router       /api/admin/new-moderator [post]
func (h *Handler) createModerator(ctx *gin.Context) {
	var input server.User

	if err := ctx.BindJSON(&input); err != nil {
		newErrorResponse(ctx, http.StatusBadRequest, err.Error())
		return
	}

	// email, password, phone_number validation
	_, err := mail.ParseAddress(input.Email)
	if err != nil {
		newErrorResponse(ctx, http.StatusBadRequest, "non valid email")
		return
	}
	if len(input.Password) < 4 {
		newErrorResponse(ctx, http.StatusBadRequest, "non valid password")
		return
	}
	if len(input.PhoneNumber) < 10 {
		newErrorResponse(ctx, http.StatusBadRequest, "non valid phone_number")
		return
	}

	id, roleId, err := h.services.Authorization.CreateModerator(input)
	if err != nil {
		newErrorResponse(ctx, http.StatusInternalServerError, err.Error())
		return
	}
	ctx.JSON(http.StatusCreated, map[string]interface{}{
		"id":      id,
		"role_id": roleId,
	})
}

// signIn godoc
// @Summary      SignIn
// @Tags         auth
// @Description  signs in account
// @ID sign-in account
// @Accept       json
// @Produce      json
// @Param        input body server.SignInInput true "sign in account info"
// @Success      200  {integer}   string 4
// @Failure      400  {object}  Error
// @Failure      404  {object}  Error
// @Failure      500  {object}  Error
// @Router       /auth/sign-in [post]
func (h *Handler) signIn(c *gin.Context) {
	var input server.SignInInput

	if err := c.BindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	accessToken, refreshToken, err := h.services.Authorization.GenerateTokens(input.Email, input.Password)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	userId, userRoleId, err := h.services.ParseToken(accessToken)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	newSession, err := h.services.CreateNewSession(&server.Session{
		UserId:       userId,
		RoleId:       userRoleId,
		RefreshToken: refreshToken,
		ClientIp:     c.ClientIP(),
		UserAgent:    c.Request.UserAgent(),
		ExpiresAt:    time.Now().Add(refreshTokenTTL),
		CreatedAt:    time.Now(),
	})

	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, "unable to create new session: "+err.Error())
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"access_token":       accessToken,
		"refresh_token":      refreshToken,
		"session_id":         newSession.Id,
		"refresh_expires_at": newSession.ExpiresAt,
	})

}

// logout godoc
// @Summary      Logout
// @Security 	 ApiKeyAuth
// @Tags         api/client
// @Description  ends session
// @ID logout from account
// @Produce      json
// @Param        input body server.User true "account info"
// @Success      200  {string}  string
// @Failure      400  {object}  Error
// @Failure      404  {object}  Error
// @Failure      500  {object}  Error
// @Router       /api/client/logout [delete]
func (h *Handler) logout(ctx *gin.Context) {
	id, err := getUserId(ctx)
	if err != nil {
		newErrorResponse(ctx, http.StatusUnauthorized, "user id not found: "+err.Error())
		return
	}

	ctx.Set(userCtx, 0)
	ctx.Set(userRoleIdCtx, 0)

	err = h.services.Authorization.Logout(id)
	if err != nil {
		newErrorResponse(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	ctx.JSON(http.StatusOK, map[string]interface{}{
		"message": "logged out, session deleted",
	})
}

type RefreshTokensInput struct {
	RefreshToken string `json:"refresh_token"`
}

// refresh godoc
// @Summary      Refresh
// @Security ApiKeyAuth
// @Tags         auth
// @Description  Gets a new access using refreshToken
// @ID refreshToken from account
// @Accept       json
// @Produce      json
// @Param        input body handler.RefreshTokensInput true "refresh token"
// @Success      200  {integer}   string 1
// @Failure      400  {object}  Error
// @Failure      404  {object}  Error
// @Failure      500  {object}  Error
// @Router       /auth/refresh [post]
func (h *Handler) refresh(ctx *gin.Context) {
	var input RefreshTokensInput

	if err := ctx.BindJSON(&input); err != nil {
		newErrorResponse(ctx, http.StatusBadRequest, err.Error())
		return
	}

	accessToken, err := h.services.Authorization.RefreshAccessToken(input.RefreshToken)
	if err != nil {
		ctx.Set(userCtx, nil)
		ctx.Set(userRoleIdCtx, nil)
		newErrorResponse(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	err = h.services.Authorization.ParseRefreshToken(input.RefreshToken)
	if err != nil {
		newErrorResponse(ctx, http.StatusBadRequest, err.Error())
		return
	}

	ctx.JSON(http.StatusOK, map[string]interface{}{
		"access_token": accessToken,
	})
}
