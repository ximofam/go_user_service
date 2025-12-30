package response

import "github.com/gin-gonic/gin"

type APIResponse[T any] struct {
	Data T   `json:"data"`
	Meta any `json:"meta,omitempty"`
}

func OK[T any](c *gin.Context, data T) {
	c.JSON(200, APIResponse[T]{
		Data: data,
	})
}

func OKWithMeta[T any](c *gin.Context, data T, meta any) {
	c.JSON(200, APIResponse[T]{
		Data: data,
		Meta: meta,
	})
}

func Message(c *gin.Context, statusCode int, message string) {
	c.JSON(statusCode, gin.H{"message": message})
}
