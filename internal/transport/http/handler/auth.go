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
	input, ok := bindJSON[dto.UserLoginInput](c)
	if !ok {
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
	input, ok := bindJSON[dto.UserResgisterInput](c)
	if !ok {
		return
	}

	if err := h.authService.Register(c.Request.Context(), &input); err != nil {
		response.Error(c, err)
		return
	}

	response.Message(c, 201, "Create user successfully")
}

func (h *AuthHandler) ChangePassword(c *gin.Context) {
	input, ok := bindJSON[dto.UserChangePasswordInput](c)
	if !ok {
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
	input, ok := bindJSON[dto.UserRefreshTokenInput](c)
	if !ok {
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
	input, ok := bindJSON[dto.UserLogoutInput](c)
	if !ok {
		return
	}

	input.UserID = c.GetUint(utils.UserIDKey)

	if err := h.authService.Logout(c.Request.Context(), &input); err != nil {
		response.Error(c, err)
		return
	}

	response.Message(c, 200, "Logout user successfully")
}
