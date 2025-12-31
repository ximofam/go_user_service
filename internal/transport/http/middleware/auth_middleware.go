package middleware

import (
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/ximofam/user-service/internal/utils"
	"github.com/ximofam/user-service/pkg/token"
)

func Auth(tokenService token.TokenService) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.AbortWithStatusJSON(400, gin.H{"error": "Missing Authorization header"})
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.AbortWithStatusJSON(400, gin.H{"error": "Authorization format must be Bearer <token>"})
			return
		}

		accessToken := parts[1]

		claims, err := tokenService.ParseToken(c.Request.Context(), accessToken)
		if err != nil {
			c.AbortWithStatusJSON(400, gin.H{"error": "Invalid token"})
			return
		}

		userID, err := strconv.Atoi(claims.UserID)
		if err != nil {
			c.AbortWithStatusJSON(400, gin.H{"error": "Missing user id"})
			return
		}

		c.Set(utils.UserIDKey, uint(userID))
		c.Set(utils.UserRoleKey, claims.UserRole)

		c.Next()
	}
}
