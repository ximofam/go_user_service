package service

import (
	"context"
	"fmt"

	"github.com/ximofam/user-service/internal/dto"
	"github.com/ximofam/user-service/internal/model"
	"github.com/ximofam/user-service/internal/repository"
	"github.com/ximofam/user-service/internal/utils"
	"github.com/ximofam/user-service/pkg/token"
)

type AuthService interface {
	Login(ctx context.Context, input *dto.UserLoginInput) (*dto.UserLoginOutput, error)
	Register(ctx context.Context, input *dto.UserResgisterInput) error
	ChangePassword(ctx context.Context, input *dto.UserChangePasswordInput) error
	RefreshToken(ctx context.Context, input *dto.UserRefreshTokenInput) (*dto.UserLoginOutput, error)
	Logout(ctx context.Context, input *dto.UserLogoutInput) error
}

type authService struct {
	userRepo     repository.UserRepository
	tokenService token.TokenService
}

func NewUserService(userRepo repository.UserRepository, tokenService token.TokenService) AuthService {
	return &authService{
		userRepo:     userRepo,
		tokenService: tokenService,
	}
}

func (s *authService) Login(ctx context.Context, input *dto.UserLoginInput) (*dto.UserLoginOutput, error) {
	user, err := s.userRepo.GetByEmail(ctx, input.Email)
	if err != nil {
		return nil, utils.ErrUnauthorized("Invalid email or password", err)
	}

	if !utils.CompareHashAndPassword(user.Password, input.Password) {
		return nil, utils.ErrUnauthorized("Invalid email or password", nil)
	}

	accessToken, err := s.tokenService.GenerateAccessToken(ctx, user)
	if err != nil {
		return nil, utils.ErrInternal("Failed to generate access token", err)
	}

	refreshToken, err := s.tokenService.GenerateRefreshToken(ctx, user)
	if err != nil {
		return nil, utils.ErrInternal("Failed to generate refresh token", err)
	}

	return &dto.UserLoginOutput{
		AccessToken:  accessToken,
		RefreshToken: refreshToken.Token,
	}, nil
}

func (s *authService) Register(ctx context.Context, input *dto.UserResgisterInput) error {
	_, err := s.userRepo.GetByEmail(ctx, input.Email)
	if err == nil {
		return utils.ErrBadRequest(fmt.Sprintf("User already exists with email: %s", input.Email), nil)
	}

	hashPassword := utils.HashPassword(input.Password)
	if hashPassword == "" {
		return utils.ErrInternal("Failed to hash password", nil)
	}

	user := model.User{
		Email:    input.Email,
		Password: hashPassword,
	}

	if err := s.userRepo.Create(ctx, &user); err != nil {
		return utils.ErrInternal("Fail to create user", err)
	}

	return nil
}

func (s *authService) ChangePassword(ctx context.Context, input *dto.UserChangePasswordInput) error {
	user, err := s.userRepo.GetByID(ctx, input.UserID)
	if err != nil {
		return utils.ErrUnauthorized("Missing or invalid user id", err)
	}

	if !utils.CompareHashAndPassword(user.Password, input.CurrentPassword) {
		return utils.ErrUnauthorized("Invalid current password", nil)
	}

	hashNewPassword := utils.HashPassword(input.NewPassword)
	if hashNewPassword == "" {
		return utils.ErrInternal("Failed to hash password", nil)
	}

	if err := s.userRepo.UpdatePassword(ctx, input.UserID, hashNewPassword); err != nil {
		return utils.ErrInternal("Failed to update password", err)
	}

	return nil
}

func (s *authService) RefreshToken(ctx context.Context, input *dto.UserRefreshTokenInput) (*dto.UserLoginOutput, error) {
	oldRefreshToken, err := s.tokenService.ValidateRefreshToken(ctx, input.Token)
	if err != nil {
		return nil, utils.ErrBadRequest(err.Error(), nil)
	}

	if err := s.tokenService.RevokeRefreshToken(ctx, oldRefreshToken.Token); err != nil {
		return nil, utils.ErrInternal("Failed to revoke token", err)
	}

	user := &model.User{
		ID: oldRefreshToken.UserID,
	}

	accessToken, err := s.tokenService.GenerateAccessToken(ctx, user)
	if err != nil {
		return nil, utils.ErrInternal("Failed to generate access token", err)
	}

	refreshToken, err := s.tokenService.GenerateRefreshToken(ctx, user)
	if err != nil {
		return nil, utils.ErrInternal("Failed to generate refresh token", err)
	}

	return &dto.UserLoginOutput{
		AccessToken:  accessToken,
		RefreshToken: refreshToken.Token,
	}, nil
}

func (s *authService) Logout(ctx context.Context, input *dto.UserLogoutInput) error {
	refreshToken, err := s.tokenService.ValidateRefreshToken(ctx, input.RefreshToken)
	if err != nil {
		return utils.ErrBadRequest("The user already logout", err)
	}

	if input.UserID != refreshToken.UserID {
		return utils.ErrForbidden("You are not alowed", nil)
	}

	return s.tokenService.RevokeRefreshToken(ctx, refreshToken.Token)
}
