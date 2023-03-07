package handler

import (
	"fmt"
	"github.com/gin-gonic/gin"
	_ "github.com/zh0vtyj/alliancecup-server/docs"
	"github.com/zh0vtyj/alliancecup-server/internal/domain/models"
	"github.com/zh0vtyj/alliancecup-server/internal/domain/user"
	"net/http"
	"net/mail"
	"strconv"
	"time"
)

const (
	authUrl            = "/auth"
	signInUrl          = "/sign-in"
	signUpUrl          = "/sign-up"
	refreshUrl         = "/refresh"
	moderatorUrl       = "/moderator"
	logoutUrl          = "/logout"
	changePasswordUrl  = "/change-password"
	restorePasswordUrl = "/restore-password"
	personalInfoUrl    = "/personal-info"
)

func (h *Handler) initAuthRoutes(router *gin.Engine) {
	auth := router.Group(authUrl)
	{
		auth.POST(signUpUrl, h.signUp)
		auth.POST(signInUrl, h.signIn)
		auth.POST(refreshUrl, h.refresh)
	}
}

func (h *Handler) initAdminModeratorsRoutes(router *gin.RouterGroup) {
	router.GET(moderatorUrl, h.getModerators)
	router.POST(moderatorUrl, h.createModerator)
	router.DELETE(moderatorUrl, h.deleteModerator)
}

func (h *Handler) initClientRoutes(group *gin.RouterGroup) {
	group.GET(personalInfoUrl, h.personalInfo)
	group.PUT(personalInfoUrl, h.updatePersonalInfo)

	group.PUT(changePasswordUrl, h.changePassword)
	group.PUT(restorePasswordUrl, h.restorePassword)
	group.DELETE(logoutUrl, h.logout)
}

// signUp godoc
// @Summary      SignUp
// @Tags         auth
// @Description  registers a new user
// @ID create account
// @Accept       json
// @Produce      json
// @Param        input body user.User true "account info"
// @Success      200  {object} 	object
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

	_, err := mail.ParseAddress(input.Email)
	if err != nil || input.Email == "" {
		newErrorResponse(ctx, http.StatusBadRequest, "invalid email")
		return
	}

	if len(input.Password) < 4 {
		newErrorResponse(ctx, http.StatusBadRequest, "invalid password")
		return
	}

	if len(input.PhoneNumber) < 10 {
		newErrorResponse(ctx, http.StatusBadRequest, "invalid phone_number")
		return
	}

	id, roleCode, err := h.services.Authorization.CreateUser(input, h.cfg.Roles.Client)
	if err != nil {
		newErrorResponse(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	ctx.JSON(http.StatusCreated, map[string]interface{}{
		"id":       id,
		"roleCode": roleCode,
	})
}

type SignInResponse struct {
	AccessToken  string `json:"accessToken" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2NjA5MDI0NzAsImlhdCI6MTY2MDg5NTI3MCwidXNlcl9pZCI6MSwidXNlcl9yb2xlX2lkIjozfQ.OTiwDdjjCkYkN7LfyOL6VWF7maKvuIpXWH2XWKFzZEo"`
	SessionId    int    `json:"sessionId" example:"15"`
	UserId       int    `json:"userId" example:"5"`
	UserRoleCode string `json:"userRoleCode" example:"5000"`
}

type SignInInput struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
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
func (h *Handler) signIn(ctx *gin.Context) {
	var input SignInInput

	if err := ctx.BindJSON(&input); err != nil {
		newErrorResponse(ctx, http.StatusBadRequest, err.Error())
		return
	}

	accessToken, refreshToken, err := h.services.Authorization.GenerateTokens(input.Email, input.Password)
	if err != nil {
		newErrorResponse(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	userId, userRoleCode, err := h.services.Authorization.ParseToken(accessToken)
	if err != nil {
		newErrorResponse(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	newSession, err := h.services.Authorization.CreateNewSession(models.Session{
		UserId:       userId,
		RoleCode:     userRoleCode,
		RefreshToken: refreshToken,
		ClientIp:     ctx.ClientIP(),
		UserAgent:    ctx.Request.UserAgent(),
		ExpiresAt:    time.Now().Add(h.cfg.Auth.JWT.RefreshTokenTTL),
		CreatedAt:    time.Now(),
	})

	if err != nil {
		newErrorResponse(ctx, http.StatusInternalServerError, "unable to create new session: "+err.Error())
		return
	}

	ctx.SetCookie(
		refreshTokenCookie,
		refreshToken,
		60*60*24*60,
		"/",
		h.cfg.Domain,
		false,
		true,
	)

	ctx.JSON(http.StatusOK, SignInResponse{
		AccessToken:  accessToken,
		UserId:       userId,
		UserRoleCode: userRoleCode,
		SessionId:    newSession.Id,
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

	ctx.Set(userIdCtx, 0)
	ctx.Set(userRoleCodeCtx, 0)

	err = h.services.Authorization.Logout(id)
	if err != nil {
		newErrorResponse(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	ctx.SetCookie(refreshTokenCookie, "", -1, "/", h.cfg.Domain, false, true)
	ctx.SetCookie(userCartCookie, "", -1, "/", h.cfg.Domain, false, true)

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
		ctx.Set(userIdCtx, 0)
		ctx.Set(userRoleCodeCtx, "")
		newErrorResponse(ctx, http.StatusUnauthorized, "refresh_token was not found "+err.Error())
		return
	}

	clientIp := ctx.ClientIP()
	userAgent := ctx.Request.UserAgent()

	err = h.services.Authorization.ParseRefreshToken(cookieToken)
	if err != nil {
		ctx.Set(userIdCtx, 0)
		ctx.Set(userRoleCodeCtx, 0)
		newErrorResponse(ctx, http.StatusUnauthorized, err.Error())
		return
	}

	accessToken, newRefreshToken, userId, userRoleCode, err := h.services.Authorization.RefreshTokens(cookieToken, clientIp, userAgent)
	if err != nil {
		ctx.Set(userIdCtx, 0)
		ctx.Set(userRoleCodeCtx, 0)
		newErrorResponse(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	ctx.SetCookie(
		refreshTokenCookie,
		newRefreshToken,
		60*60*24*60,
		"/",
		h.cfg.Domain,
		false,
		true,
	)

	ctx.JSON(http.StatusOK, map[string]interface{}{
		"accessToken":  accessToken,
		"userId":       userId,
		"userRoleCode": userRoleCode,
	})

}

type ChangePasswordInput struct {
	OldPassword string `json:"oldPassword" binding:"required"`
	NewPassword string `json:"newPassword" binding:"required"`
}

// changePassword godoc
// @Summary Client change password
// @Security ApiKeyAuth
// @Tags api/client
// @Description Changes user password
// @ID change user password
// @Accept json
// @Produce json
// @Param input body handler.ChangePasswordInput true "Info to change password"
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

type RestorePasswordInput struct {
	NewPassword string `json:"password" binding:"required"`
}

// restorePassword godoc
// @Summary Client restore password
// @Security ApiKeyAuth
// @Tags api/client
// @Description Restores user password
// @ID change user password
// @Accept json
// @Produce json
// @Param input body handler.ChangePasswordInput true "Info to change password"
// @Success 200  {object} object
// @Failure 400  {object} Error
// @Failure 401  {object} Error
// @Failure 500  {object} Error
// @Router /api/client/restore-password [put]
func (h *Handler) restorePassword(ctx *gin.Context) {
	id, err := getUserId(ctx)
	if err != nil {
		newErrorResponse(ctx, http.StatusBadRequest, err.Error())
		return
	}

	var input RestorePasswordInput
	if err = ctx.BindJSON(&input); err != nil {
		newErrorResponse(ctx, http.StatusBadRequest, err.Error())
		return
	}

	if len(input.NewPassword) < 4 {
		newErrorResponse(ctx, http.StatusBadRequest, fmt.Errorf("invalid password lenght").Error())
		return
	}

	err = h.services.Authorization.RestorePassword(id, input.NewPassword)
	if err != nil {
		newErrorResponse(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	ctx.JSON(http.StatusOK, map[string]any{
		"message": "password changed",
	})
}

type ForgotPasswordInput struct {
	Email string `json:"email" binding:"required"`
}

func (h *Handler) forgotPassword(ctx *gin.Context) {
	var input ForgotPasswordInput

	if err := ctx.BindJSON(&input); err != nil {
		newErrorResponse(ctx, http.StatusBadRequest, err.Error())
		return
	}

	err := h.services.Authorization.UserForgotPassword(input.Email)
	if err != nil {
		newErrorResponse(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	ctx.JSON(http.StatusOK, "email was successfully send")
}

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

// createModerator godoc
// @Summary      CreateModerator
// @Security 	 ApiKeyAuth
// @Tags         api/super/admin
// @Description  registers a new moderator
// @ID 			 create account for moderator
// @Accept       json
// @Produce      json
// @Param        input body user.User true "account info"
// @Success      200  {object}  object
// @Failure      400  {object}  Error
// @Failure      404  {object}  Error
// @Failure      500  {object}  Error
// @Router       /api/admin/super/moderator [post]
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

	id, roleCode, err := h.services.Authorization.CreateUser(input, h.cfg.Roles.Moderator)
	if err != nil {
		newErrorResponse(ctx, http.StatusInternalServerError, err.Error())
		return
	}
	ctx.JSON(http.StatusCreated, map[string]interface{}{
		"id":       id,
		"roleCode": roleCode,
	})
}

// getModerators godoc
// @Summary      Get moderators
// @Security 	 ApiKeyAuth
// @Tags         api/super/admin
// @Description  registers a new moderator
// @ID 			 create account for moderator
// @Accept       json
// @Produce      json
// @Param        input body user.User true "account info"
// @Success      200  {object}  object
// @Failure      400  {object}  Error
// @Failure      404  {object}  Error
// @Failure      500  {object}  Error
// @Router       /api/admin/super/moderator [get]
func (h *Handler) getModerators(ctx *gin.Context) {
	createdAt := ctx.Query("createdAt")

	moderators, err := h.services.Authorization.GetModerators(createdAt, h.cfg.Roles.Moderator)
	if err != nil {
		newErrorResponse(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	ctx.JSON(http.StatusOK, moderators)
}

// deleteModerator godoc
// @Summary      Delete moderator
// @Security 	 ApiKeyAuth
// @Tags         api/super/admin
// @Description  Deletes moderator
// @ID 			 delete moderator
// @Accept       json
// @Produce      json
// @Param        id query int true "Moderator id"
// @Success      200  {string}  string
// @Failure      400  {object}  Error
// @Failure      404  {object}  Error
// @Failure      500  {object}  Error
// @Router       /api/admin/super/moderator [delete]
func (h *Handler) deleteModerator(ctx *gin.Context) {
	id := ctx.Query("id")
	userId, err := strconv.Atoi(id)
	if err != nil {
		newErrorResponse(ctx, http.StatusInternalServerError, "invalid moderator id")
		return
	}

	err = h.services.Authorization.Delete(userId)
	if err != nil {
		newErrorResponse(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	ctx.JSON(http.StatusOK, "moderator successfully deleted")
}
