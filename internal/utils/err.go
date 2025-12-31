package utils

const (
	ErrCodeBadRequest      = "BAD_REQUEST"
	ErrCodeNotFound        = "NOT_FOUND"
	ErrCodeConflict        = "CONFLICT"
	ErrCodeInternal        = "INTERNAL"
	ErrCodeUnauthorized    = "UNAUTHORIZED"
	ErrCodeForbidden       = "FORBIDDEN"
	ErrCodeTooManyRequests = "TOO_MANY_REQUESTS"
)

type AppError struct {
	Code    string
	Message string
	Err     any
}

func (a *AppError) Error() string {
	return a.Message
}

func NewAppError(code, message string, err any) *AppError {
	return &AppError{
		Code:    code,
		Message: message,
		Err:     err,
	}
}

func ErrBadRequest(message string, err any) *AppError {
	return NewAppError(ErrCodeBadRequest, message, err)
}

func ErrNotFound(message string, err any) *AppError {
	return NewAppError(ErrCodeNotFound, message, err)
}

func ErrConflict(message string, err any) *AppError {
	return NewAppError(ErrCodeConflict, message, err)
}

func ErrInternal(message string, err any) *AppError {
	return NewAppError(ErrCodeInternal, message, err)
}

func ErrUnauthorized(message string, err any) *AppError {
	return NewAppError(ErrCodeUnauthorized, message, err)
}

func ErrForbidden(message string, err any) *AppError {
	return NewAppError(ErrCodeForbidden, message, err)
}

func ErrTooManyRequests(message string, err any) *AppError {
	return NewAppError(ErrCodeTooManyRequests, message, err)
}
