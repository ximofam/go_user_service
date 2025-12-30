package token

import (
	"context"

	"github.com/golang-jwt/jwt/v5"
	"github.com/ximofam/user-service/internal/model"
)

type TokenService interface {
	GenerateAccessToken(ctx context.Context, user *model.User) (string, error)
	GenerateRefreshToken(ctx context.Context, user *model.User) (RefreshToken, error)
	RevokeRefreshToken(ctx context.Context, token string) error
	ValidateRefreshToken(ctx context.Context, token string) (RefreshToken, error)
	ParseToken(ctx context.Context, tokenStr string) (MyClaims, error)
}

type MyClaims struct {
	jwt.RegisteredClaims
	UserID   string
	UserRole string
}

type RefreshToken struct {
	Token  string
	UserID uint
}
