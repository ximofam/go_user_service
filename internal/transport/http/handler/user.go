package handler

import (
	"context"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/ximofam/user-service/internal/dto"
	"github.com/ximofam/user-service/internal/service"
	"github.com/ximofam/user-service/internal/transport/http/response"
	"github.com/ximofam/user-service/internal/utils"
)

type UserHandler struct {
	userService service.UserService
}

func NewUserHandler(userService service.UserService) *UserHandler {
	return &UserHandler{userService: userService}
}

func (h *UserHandler) ListUsers(c *gin.Context) {
	input, ok := bindQuery[dto.ListUsersInput](c)
	if !ok {
		return
	}

	users, pagination, err := h.userService.ListUsers(context.Background(), &input)
	if err != nil {
		response.Error(c, err)
		return
	}

	response.OKWithMeta(c, users, pagination)
}

func (h *UserHandler) SoftDelete(c *gin.Context) {
	userRole := c.GetString(utils.UserRoleKey)
	if userRole == "" {
		response.Error(c, utils.ErrForbidden("You are not allowed to use this service", nil))
		return
	}

	DeleteUserID, err := strconv.Atoi(c.Param("id"))
	if err != nil || DeleteUserID == 0 {
		response.Error(c, utils.ErrBadRequest("Missing or invalid user id to delete", err))
		return
	}

	UserID := c.GetUint(utils.UserIDKey)
	if UserID == 0 {
		response.Error(c, utils.ErrUnauthorized("Missing user id", nil))
		return
	}

	if DeleteUserID == int(UserID) {
		response.Error(c, utils.ErrForbidden("You can't delete yourself", nil))
		return
	}

	if err := h.userService.SoftDelete(c.Request.Context(), uint(DeleteUserID)); err != nil {
		response.Error(c, err)
		return
	}

	response.Message(c, 200, "Soft delete user successfullty")
}
