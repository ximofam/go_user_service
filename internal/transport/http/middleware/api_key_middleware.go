package middleware

import (
	"github.com/gin-gonic/gin"
)

func APIKeyAuth(validKey string) gin.HandlerFunc {
	return func(c *gin.Context) {
		apiKey := c.GetHeader("X-API-KEY")

		if apiKey == "" {
			c.AbortWithStatusJSON(401, gin.H{"error": "Missing header X-API-KEY"})
			return
		}

		if apiKey != validKey {
			c.AbortWithStatusJSON(401, gin.H{"error": "Invalid api key"})
			return
		}

		c.Next()
	}
}
