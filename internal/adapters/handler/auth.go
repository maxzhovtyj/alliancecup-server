package handler

import (
	"fmt"
	"github.com/gin-gonic/gin"
	_ "github.com/zh0vtyj/allincecup-server/docs"
	"github.com/zh0vtyj/allincecup-server/internal/domain/models"
	"github.com/zh0vtyj/allincecup-server/internal/domain/user"
	"net/http"
	"net/mail"
	"os"
	"time"
)

const refreshTokenTTL = 1440 * time.Hour

type SignInResponse struct {
	AccessToken string `json:"accessToken" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2NjA5MDI0NzAsImlhdCI6MTY2MDg5NTI3MCwidXNlcl9pZCI6MSwidXNlcl9yb2xlX2lkIjozfQ.OTiwDdjjCkYkN7LfyOL6VWF7maKvuIpXWH2XWKFzZEo"`
	SessionId   int    `json:"sessionId" example:"15"`
	UserId      int    `json:"userId" example:"5"`
	UserRoleId  int    `json:"userRoleId" example:"1"`
}

type SignInInput struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type ChangePasswordInput struct {
	OldPassword string `json:"oldPassword" binding:"required"`
	NewPassword string `json:"newPassword" binding:"required"`
}

// signUp godoc
// @Summary      SignUp
// @Tags         auth
// @Description  registers a new user
// @ID create account
// @Accept       json
// @Produce      json
// @Param        input body user.User true "account info"
// @Success      200  {integer} integer 2
// @Failure      400  {object}  Error
// @Failure      404  {object}  Error
// @Failure      500  {object}  Error
// @Router       /auth/sign-up [post]
func (h *Handler) signUp(ctx *gin.Context) {
	var input user.User

	if err := ctx.BindJSON(&input); err != nil {
		newErrorResponse(ctx, http.StatusBadRequest, err.Error())
		return
	}

	// email, password, phone_number validation
	_, err := mail.ParseAddress(input.Email)
	if err != nil || input.Email == "" {
		newErrorResponse(ctx, http.StatusBadRequest, "invalid email")
		return
	}
	if len(input.Password) < 4 {
		newErrorResponse(ctx, http.StatusBadRequest, "invalid password")
		return
	}

	// 068 306 29 75
	if len(input.PhoneNumber) < 10 {
		newErrorResponse(ctx, http.StatusBadRequest, "invalid phone_number")
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
// @Security 	 ApiKeyAuth
// @Tags         api/admin
// @Description  registers a new moderator
// @ID 			 create account for moderator
// @Accept       json
// @Produce      json
// @Param        input body user.User true "account info"
// @Success      200  {object}  object
// @Failure      400  {object}  Error
// @Failure      404  {object}  Error
// @Failure      500  {object}  Error
// @Router       /api/admin/moderator [post]
func (h *Handler) createModerator(ctx *gin.Context) {
	var input user.User

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
// @Param        input body handler.SignInInput true "sign in account info"
// @Success      200  {object} handler.SignInResponse
// @Failure      400  {object}  Error
// @Failure      404  {object}  Error
// @Failure      500  {object}  Error
// @Router       /auth/sign-in [post]
func (h *Handler) signIn(c *gin.Context) {
	var input SignInInput

	if err := c.BindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	accessToken, refreshToken, err := h.services.Authorization.GenerateTokens(input.Email, input.Password)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	userId, userRoleId, err := h.services.Authorization.ParseToken(accessToken)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	newSession, err := h.services.Authorization.CreateNewSession(&models.Session{
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

	dm := os.Getenv(domain)
	c.SetCookie(refreshTokenCookie, refreshToken, 60*60*24*60, "/", dm, false, true)

	c.JSON(http.StatusOK, SignInResponse{
		AccessToken: accessToken,
		UserId:      userId,
		UserRoleId:  userRoleId,
		SessionId:   newSession.Id,
	})
}

// logout godoc
// @Summary      Logout
// @Security 	 ApiKeyAuth
// @Tags         api/client
// @Description  ends session
// @ID logout from account
// @Produce      json
// @Success      200  {string}  string
// @Failure      400  {object}  Error
// @Failure      404  {object}  Error
// @Failure      500  {object}  Error
// @Router       /api/client/logout [delete]
func (h *Handler) logout(ctx *gin.Context) {
	id, err := getUserId(ctx)
	if err != nil {
		newErrorResponse(ctx, http.StatusUnauthorized, fmt.Errorf("user id not found: %v", err).Error())
		return
	}

	ctx.Set(userCtx, 0)
	ctx.Set(userRoleIdCtx, 0)

	err = h.services.Authorization.Logout(id)
	if err != nil {
		newErrorResponse(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	dm := os.Getenv(domain)
	ctx.SetCookie(refreshTokenCookie, "", -1, "/", dm, false, true)
	ctx.SetCookie(userCartCookie, "", -1, "/", dm, false, true)

	ctx.JSON(http.StatusOK, map[string]interface{}{
		"message": "logged out, session deleted",
	})
}

// refresh godoc
// @Summary      Refresh
// @Tags         auth
// @Description  Gets a new access token using refreshToken
// @ID refreshes token from account
// @Produce      json
// @Success      200  {object}  string
// @Failure      400  {object}  Error
// @Failure      404  {object}  Error
// @Failure      500  {object}  Error
// @Router       /auth/refresh [post]
func (h *Handler) refresh(ctx *gin.Context) {
	cookieToken, err := ctx.Cookie(refreshTokenCookie)
	if err != nil {
		ctx.Set(userCtx, 0)
		ctx.Set(userRoleIdCtx, 0)
		newErrorResponse(ctx, http.StatusUnauthorized, "refresh_token was not found "+err.Error())
		return
	}

	clientIp := ctx.ClientIP()
	userAgent := ctx.Request.UserAgent()

	err = h.services.Authorization.ParseRefreshToken(cookieToken)
	if err != nil {
		ctx.Set(userCtx, 0)
		ctx.Set(userRoleIdCtx, 0)
		newErrorResponse(ctx, http.StatusUnauthorized, err.Error())
		return
	}

	accessToken, newRefreshToken, userId, userRoleId, err := h.services.Authorization.RefreshTokens(cookieToken, clientIp, userAgent)
	if err != nil {
		ctx.Set(userCtx, 0)
		ctx.Set(userRoleIdCtx, 0)
		newErrorResponse(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	dm := os.Getenv(domain)
	ctx.SetCookie(refreshTokenCookie, newRefreshToken, 60*60*24*60, "/", dm, false, true)

	ctx.JSON(http.StatusOK, map[string]interface{}{
		"accessToken": accessToken,
		"userId":      userId,
		"userRoleId":  userRoleId,
	})
}

// changePassword godoc
// @Summary Client change password
// @Security ApiKeyAuth
// @Tags api/client
// @Description Changes user password
// @ID change user password
// @Accept json
// @Produce json
// @Param input body handler.ChangePasswordInput true "Order to change password"
// @Success 200  {object} object
// @Failure 400  {object} Error
// @Failure 401  {object} Error
// @Failure 500  {object} Error
// @Router /api/client/change-password [put]
func (h *Handler) changePassword(ctx *gin.Context) {
	id, err := getUserId(ctx)
	if err != nil {
		newErrorResponse(ctx, http.StatusBadRequest, err.Error())
		return
	}

	var input ChangePasswordInput
	if err = ctx.BindJSON(&input); err != nil {
		newErrorResponse(ctx, http.StatusBadRequest, err.Error())
		return
	}

	if len(input.NewPassword) < 4 {
		newErrorResponse(ctx, http.StatusBadRequest, fmt.Errorf("invalid password lenght").Error())
		return
	}

	err = h.services.Authorization.ChangePassword(id, input.OldPassword, input.NewPassword)
	if err != nil {
		newErrorResponse(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	ctx.JSON(http.StatusOK, map[string]any{
		"message": "password changed",
	})
}

// TODO
func (h *Handler) forgotPassword(ctx *gin.Context) {
	var input string

	if err := ctx.BindJSON(&input); err != nil {
		newErrorResponse(ctx, http.StatusBadRequest, err.Error())
		return
	}

	err := h.services.Authorization.UserForgotPassword(input)
	if err != nil {
		newErrorResponse(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	ctx.JSON(http.StatusOK, "email was successfully send")
}

// TODO
func (h *Handler) personalInfo(ctx *gin.Context) {
	id, err := getUserId(ctx)
	if err != nil {
		newErrorResponse(ctx, http.StatusBadRequest, err.Error())
		return
	}

	userInfo, err := h.services.Authorization.UserInfo(id)
	if err != nil {
		newErrorResponse(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	ctx.JSON(http.StatusOK, userInfo)
}

// TODO
func (h *Handler) updatePersonalInfo(ctx *gin.Context) {
	id, err := getUserId(ctx)
	if err != nil {
		newErrorResponse(ctx, http.StatusBadRequest, err.Error())
		return
	}

	var userInput user.InfoDTO
	if err = ctx.BindJSON(&userInput); err != nil {
		newErrorResponse(ctx, http.StatusBadRequest, err.Error())
		return
	}

	err = h.services.Authorization.ChangePersonalInfo(userInput, id)
	if err != nil {
		newErrorResponse(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	ctx.JSON(http.StatusOK, "user info updated")
}
