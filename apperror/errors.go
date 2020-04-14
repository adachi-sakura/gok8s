package apperror

import (
	"fmt"
	"net/http"
)

type AppError struct {
	ErrorCode ErrorCode		`json:"errorCode"`
	Message   string		`json:"message"`
}

func (err *AppError) StatusCode() int {
	if val, exists := errorCodesMap[err.ErrorCode]; exists {
		return val
	} else {
		return http.StatusInternalServerError
	}
}

func Wrap(err interface{}) *AppError {
	if IsAppError(err) {
		return err.(*AppError)
	}
	return &AppError{
		ErrorCode:	InternalServerError,
		Message:	"Non AppError Occurred",
	}
}

func IsAppError(err interface{}) bool {
	if err == nil {
		return false
	}
	switch err.(type) {
	case *AppError:
		return true
	case error:
		return false
	default:
		return false
	}
}

func (err *AppError) Error() string {
	result := fmt.Sprintf("[%s] %s", string(err.ErrorCode), err.Message)
	return result
}
