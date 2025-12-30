package response

import (
	"errors"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/ximofam/user-service/internal/utils"
)

type ErrorResponse struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

func Error(c *gin.Context, err error) {
	var appErr *utils.AppError

	status := http.StatusInternalServerError
	resp := ErrorResponse{
		Code: utils.ErrCodeInternal,
	}

	if errors.As(err, &appErr) {
		status = httpStatusFromCode(appErr.Code)
		resp.Code = appErr.Code
		resp.Message = appErr.Message

		if appErr.Err != nil {
			log.Println(appErr.Err)
		}
	} else {
		resp.Message = err.Error()
	}

	c.JSON(status, gin.H{
		"error": resp,
	})
}

func httpStatusFromCode(code string) int {
	switch code {
	case utils.ErrCodeBadRequest:
		return http.StatusBadRequest
	case utils.ErrCodeNotFound:
		return http.StatusNotFound
	case utils.ErrCodeConflict:
		return http.StatusConflict
	case utils.ErrCodeUnauthorized:
		return http.StatusUnauthorized
	case utils.ErrCodeTooManyRequests:
		return http.StatusTooManyRequests
	default:
		return http.StatusInternalServerError
	}
}
