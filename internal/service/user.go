package service

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/ximofam/user-service/internal/dto"
	"github.com/ximofam/user-service/internal/repository"
	"github.com/ximofam/user-service/internal/utils"
)

type UserService interface {
	ListUsers(ctx context.Context, input *dto.ListUsersInput) ([]dto.ListUsersOutput, *dto.PagingOutput, error)
	SoftDelete(ctx context.Context, userID uint) error
}

type userService struct {
	userRepo repository.UserRepository
}

func NewUserService(userRepo repository.UserRepository) *userService {
	return &userService{userRepo: userRepo}
}

func (s *userService) ListUsers(ctx context.Context, input *dto.ListUsersInput) ([]dto.ListUsersOutput, *dto.PagingOutput, error) {
	users, pagination, err := s.userRepo.GetAll(ctx, input)
	if err != nil {
		return nil, nil, utils.ErrInternal("Failed to get list users", err)
	}

	res := []dto.ListUsersOutput{}
	if err := utils.Copy(&users, &res); err != nil {
		return nil, nil, utils.ErrInternal("Failed to map users to userDTO", err)
	}

	return res, pagination, nil
}

func (s *userService) SoftDelete(ctx context.Context, userID uint) error {
	if err := s.userRepo.SoftDelete(ctx, userID); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return utils.ErrBadRequest(fmt.Sprintf("The user with id: %d does not exists", userID), nil)
		}

		return utils.ErrInternal("Failed to soft delete user", err)
	}

	return nil
}
