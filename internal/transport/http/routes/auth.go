package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/ximofam/user-service/internal/db"
	"github.com/ximofam/user-service/internal/repository"
	"github.com/ximofam/user-service/internal/service"
	"github.com/ximofam/user-service/internal/transport/http/handler"
	"github.com/ximofam/user-service/internal/transport/http/middleware"
	"github.com/ximofam/user-service/pkg/token"
)

func AuthRoutes(r *gin.Engine, dbProvider *db.DBProvider, tokenService token.TokenService) {
	userRepo := repository.NewUserRepository(dbProvider)
	authService := service.NewAuthService(userRepo, tokenService)
	authHandler := handler.NewAuthHandler(authService)
	auth := r.Group("/auth")
	{
		auth.POST("/login", authHandler.Login)
		auth.POST("/register", authHandler.Register)
		auth.POST("/change-password", middleware.Auth(tokenService), authHandler.ChangePassword)
		auth.POST("/refresh-token", authHandler.RefreshToken)
		auth.POST("/logout", middleware.Auth(tokenService), authHandler.Logout)
	}
}
