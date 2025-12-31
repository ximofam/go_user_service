package response

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/ximofam/user-service/internal/utils"
)

type ErrorResponse struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	Details any    `json:"details,omitempty"`
}

func Error(c *gin.Context, err error) {
	status := http.StatusInternalServerError
	resp := ErrorResponse{
		Code: utils.ErrCodeInternal,
	}

	var appErr *utils.AppError
	if errors.As(err, &appErr) {
		status = httpStatusFromCode(appErr.Code)
		resp.Code = appErr.Code
		resp.Message = appErr.Message

		if appErr.Err != nil {
			if e, ok := appErr.Err.(error); ok {
				resp.Details = e.Error()
			} else {
				resp.Details = appErr.Err
			}
		}
	} else {
		resp.Details = err.Error()
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
