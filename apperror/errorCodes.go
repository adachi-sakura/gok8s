package apperror

import (
	"fmt"
	"net/http"
)

type ErrorCode string

const (
	NotAuthenticatedError		= "NotAuthenticated"
	InternalServerError			= "InternalServerError"
	InvalidHeader				= "InvalidHeader"
	InvalidParameter			= "InvalidParameter"
	AuthorizationError			= "AuthorizationError"
	InvalidRequestBody			= "InvalidRequestBody"
	ResourceNotFount			= "ResourceNotFound"
)

var errorCodesMap = map[ErrorCode]int {
	InvalidHeader:		http.StatusBadRequest,
	InvalidParameter:	http.StatusBadRequest,
	InvalidRequestBody:	http.StatusBadRequest,
	NotAuthenticatedError:	http.StatusUnauthorized,
	AuthorizationError:		http.StatusUnauthorized,
	ResourceNotFount:		http.StatusNotFound,
	InternalServerError:	http.StatusInternalServerError,
}

func NewInvalidHeaderError(name string) *AppError {
	appError := &AppError{
		ErrorCode:	InvalidHeader,
		Message:	"header is invalid or missing "+name,
	}
	return appError
}

func NewParameterRequiredError(name string) *AppError {
	appError := &AppError{
		ErrorCode:	InvalidParameter,
		Message:	"parameter is required "+name,
	}
	return appError
}

func NewHeaderRequiredError(name string) *AppError {
	appError := &AppError{
		ErrorCode:	InvalidHeader,
		Message:	"header is required "+name,
	}
	return appError
}

func NewInvalidParameterError(name string, cause error) *AppError {
	appError := &AppError{
		ErrorCode:	InvalidParameter,
		Message:	fmt.Sprintf("parameter is invalid: %s, cause: %s", name, cause.Error()),
	}
	return appError
}

func NewInvalidRequestBodyError(cause error) *AppError {
	appError := &AppError{
		ErrorCode:	InvalidRequestBody,
		Message:	"request body is invalid: "+cause.Error(),
	}
	return appError
}

func NewNotAuthenticatedError() *AppError {
	appError := &AppError{
		ErrorCode:	NotAuthenticatedError,
		Message: "not authenticated",
	}
	return appError
}

func NewAuthorizationError() *AppError {
	appError := &AppError{
		ErrorCode:	AuthorizationError,
		Message: "authorization error",
	}
	return appError
}

func NewResourceNotFoundError(target string) *AppError {
	appError := &AppError{
		ErrorCode:	ResourceNotFount,
		Message:	"resource not found: "+target,
	}
	return appError
}

func NewInternalServerError(cause string) *AppError {
	appError := &AppError{
		ErrorCode:	InternalServerError,
		Message:	cause,
	}
	return appError
}
