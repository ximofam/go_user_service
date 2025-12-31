package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/ximofam/user-service/internal/db"
	"github.com/ximofam/user-service/internal/model"
	"github.com/ximofam/user-service/internal/repository"
	"github.com/ximofam/user-service/internal/service"
	"github.com/ximofam/user-service/internal/transport/http/handler"
	"github.com/ximofam/user-service/internal/transport/http/middleware"
	"github.com/ximofam/user-service/pkg/token"
)

func UserRoutes(r *gin.Engine, dbProvider *db.DBProvider, tokenService token.TokenService) {
	userRepo := repository.NewUserRepository(dbProvider)
	userService := service.NewUserService(userRepo)
	userHandler := handler.NewUserHandler(userService)
	user := r.Group("/api/v1/users", middleware.Auth(tokenService))
	{
		user.GET("/", userHandler.ListUsers)
		user.DELETE("/:id", middleware.RoleMiddleware(model.RoleAdmin), userHandler.SoftDelete)
	}
}
