package token

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/ximofam/user-service/internal/model"
	"github.com/ximofam/user-service/pkg/cache"
)

var refreshTokenKey = "refresh_token"

type jwtService struct {
	accessTokenTTL  time.Duration
	refreshTokenTTL time.Duration
	secretKey       []byte
	cacheService    cache.CacheService
}

func NewTokenService(accessTokenTTL, refreshTokenTTL time.Duration, secretKey []byte, cacheService cache.CacheService) *jwtService {
	return &jwtService{
		accessTokenTTL:  accessTokenTTL,
		refreshTokenTTL: refreshTokenTTL,
		secretKey:       secretKey,
		cacheService:    cacheService,
	}
}

func (s *jwtService) GenerateAccessToken(ctx context.Context, user *model.User) (string, error) {
	claims := MyClaims{
		UserID:   strconv.Itoa(int(user.ID)),
		UserRole: user.Role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(s.accessTokenTTL)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(s.secretKey)
}

func (s *jwtService) GenerateRefreshToken(ctx context.Context, user *model.User) (RefreshToken, error) {
	b := make([]byte, 32) // 256-bit
	if _, err := rand.Read(b); err != nil {
		return RefreshToken{}, err
	}

	token := base64.URLEncoding.EncodeToString(b)
	refreshToken := RefreshToken{
		Token:  token,
		UserID: user.ID,
	}

	cacheKey := fmt.Sprintf("%s:%s", refreshTokenKey, token)
	s.cacheService.Set(ctx, cacheKey, refreshToken, s.refreshTokenTTL)

	return refreshToken, nil
}

func (s *jwtService) RevokeRefreshToken(ctx context.Context, token string) error {
	cacheKey := fmt.Sprintf("%s:%s", refreshTokenKey, token)
	// var refreshToken RefreshToken
	// if err := s.cacheService.Get(ctx, cacheKey, &refreshToken); err != nil {
	// 	return errors.New("the token has been revoked or does not exists")
	// }

	return s.cacheService.Del(ctx, cacheKey)
}

func (s *jwtService) ValidateRefreshToken(ctx context.Context, token string) (RefreshToken, error) {
	cacheKey := fmt.Sprintf("%s:%s", refreshTokenKey, token)
	var refreshToken RefreshToken
	if err := s.cacheService.Get(ctx, cacheKey, &refreshToken); err != nil {
		return RefreshToken{}, errors.New("the token has been revoked or does not exists")
	}

	return refreshToken, nil
}

func (s *jwtService) ParseToken(ctx context.Context, tokenStr string) (MyClaims, error) {
	claims := MyClaims{}

	token, err := jwt.ParseWithClaims(tokenStr, &claims, func(t *jwt.Token) (interface{}, error) {
		return s.secretKey, nil
	})
	if err != nil {
		return MyClaims{}, err
	}

	if !token.Valid {
		return MyClaims{}, errors.New("invalid token")
	}

	return claims, nil
}
