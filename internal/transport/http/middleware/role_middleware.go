package middleware

import (
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/ximofam/user-service/internal/utils"
)

func RoleMiddleware(roles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		userRole := c.GetString(utils.UserRoleKey)

		if userRole == "" {
			c.AbortWithStatusJSON(403, gin.H{"error": "You are not allowed to access this service"})
			return
		}

		for _, role := range roles {
			if strings.EqualFold(userRole, role) {
				c.Next()
				return
			}
		}

		c.AbortWithStatusJSON(403, gin.H{"error": "You are not allowed to access this service"})
	}
}
