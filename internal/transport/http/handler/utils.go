package handler

import (
	"errors"
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/ximofam/user-service/internal/transport/http/response"
	"github.com/ximofam/user-service/internal/utils"
)

func bindJSON[T any](c *gin.Context) (T, bool) {
	var res T
	if err := c.ShouldBindJSON(&res); err != nil {
		var vErrs validator.ValidationErrors
		if errors.As(err, &vErrs) {
			out := make([]string, len(vErrs))
			for i, v := range vErrs {
				out[i] = getValidateErrorMsg(v)
			}

			response.Error(c, utils.ErrBadRequest("Invalid request data", out))
			return res, false
		}

		response.Error(c, utils.ErrBadRequest("Invalid JSON format", err))
	}

	return res, true
}

func bindQuery[T any](c *gin.Context) (T, bool) {
	var res T
	if err := c.ShouldBindQuery(&res); err != nil {
		var vErrs validator.ValidationErrors
		if errors.As(err, &vErrs) {
			out := make([]string, len(vErrs))
			for i, v := range vErrs {
				out[i] = getValidateErrorMsg(v)
			}

			response.Error(c, utils.ErrBadRequest("Invalid request data", out))
			return res, false
		}

		response.Error(c, utils.ErrBadRequest("Invalid JSON format", err))
	}

	return res, true
}

func getValidateErrorMsg(fe validator.FieldError) string {
	switch fe.Tag() {
	case "required":
		return fmt.Sprintf("Field '%s' is required", fe.Field())
	case "email":
		return "Invalid email format"
	case "min":
		return fmt.Sprintf("Field '%s' must be at least %s characters", fe.Field(), fe.Param())
	case "max":
		return fmt.Sprintf("Field '%s' must be at most %s characters", fe.Field(), fe.Param())
	case "len":
		return fmt.Sprintf("Field '%s' must be exactly %s characters", fe.Field(), fe.Param())
	case "alphanum":
		return fmt.Sprintf("Field '%s' must contain only alphanumeric characters", fe.Field())
	}

	return fmt.Sprintf("Validation failed on field '%s' (condition: %s)", fe.Field(), fe.Tag())
}
