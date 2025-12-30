package handler

import (
	"context"

	"github.com/gin-gonic/gin"
	"github.com/ximofam/user-service/internal/dto"
	"github.com/ximofam/user-service/internal/service"
	"github.com/ximofam/user-service/internal/transport/http/response"
	"github.com/ximofam/user-service/internal/utils"
)

type AuthHandler struct {
	authService service.AuthService
}

func NewUserHandler(authService service.AuthService) *AuthHandler {
	return &AuthHandler{authService: authService}
}

func (h *AuthHandler) Login(c *gin.Context) {
	var input dto.UserLoginInput
	if err := c.ShouldBindJSON(&input); err != nil {
		response.Error(c, utils.ErrBadRequest("Invalid request data", err))
		return
	}

	output, err := h.authService.Login(c.Request.Context(), &input)
	if err != nil {
		response.Error(c, err)
		return
	}

	response.OK(c, output)
}

func (h *AuthHandler) Register(c *gin.Context) {
	var input dto.UserResgisterInput
	if err := c.ShouldBindJSON(&input); err != nil {
		response.Error(c, utils.ErrBadRequest("Invalid request data", err))
		return
	}

	if err := h.authService.Register(c.Request.Context(), &input); err != nil {
		response.Error(c, err)
		return
	}

	response.Message(c, 201, "Create user successfully")
}

func (h *AuthHandler) ChangePassword(c *gin.Context) {
	var input dto.UserChangePasswordInput
	if err := c.ShouldBindJSON(&input); err != nil {
		response.Error(c, utils.ErrBadRequest("Invalid request data", err))
		return
	}

	input.UserID = c.GetUint(utils.UserIDKey)

	if err := h.authService.ChangePassword(c.Request.Context(), &input); err != nil {
		response.Error(c, err)
		return
	}

	response.Message(c, 200, "Change password successfully")
}

func (h *AuthHandler) RefreshToken(c *gin.Context) {
	var input dto.UserRefreshTokenInput
	if err := c.ShouldBindJSON(&input); err != nil {
		response.Error(c, utils.ErrBadRequest("Invalid request data", err))
		return
	}

	output, err := h.authService.RefreshToken(context.Background(), &input)
	if err != nil {
		response.Error(c, err)
		return
	}

	response.OK(c, output)
}

func (h *AuthHandler) Logout(c *gin.Context) {
	var input dto.UserLogoutInput
	if err := c.ShouldBindJSON(&input); err != nil {
		response.Error(c, utils.ErrBadRequest("Invalid request data", err))
		return
	}

	input.UserID = c.GetUint(utils.UserIDKey)

	if err := h.authService.Logout(c.Request.Context(), &input); err != nil {
		response.Error(c, err)
		return
	}

	response.Message(c, 200, "Logout user successfully")
}
