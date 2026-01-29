package exceptions

import (
	"net/http"
)

type ErrorCustom struct {
	Message    string
	Errors     any
	StatusCode int
}

func NewError(message string, err any, statusCode int) *ErrorCustom {
	return &ErrorCustom{
		Message:    message,
		Errors:     err,
		StatusCode: statusCode,
	}
}

func (err *ErrorCustom) Error() string {
	return err.Message
}

func (err *ErrorCustom) WithMessage(msg string) *ErrorCustom {
	return &ErrorCustom{
		Message:    msg,
		Errors:     err.Errors,
		StatusCode: err.StatusCode,
	}
}

const (
	DefaultMsgBadRequest      = "bad request"
	DefaultMsgValidationError = "validation error"
	DefaultMsgUnauthorized    = "unauthorized"
	DefaultMsgForbidden       = "access denied"
	DefaultMsgNotFound        = "not found"
	DefaultMsgConflict        = "conflict"
)

// CUSTOM ERROR
var (
	ErrCustomInvalidPayload = NewError(DefaultMsgBadRequest, "invalid payload", http.StatusBadRequest)

	ErrCustomValidation        = NewError(DefaultMsgValidationError, nil, http.StatusUnprocessableEntity)
	ErrCustomInvalidCategoryId = NewError(DefaultMsgValidationError, "invalid category id", http.StatusUnprocessableEntity)
	ErrCustomInvalidAddressId  = NewError(DefaultMsgValidationError, "invalid address id", http.StatusUnprocessableEntity)
	ErrCustomInvalidProductId  = NewError(DefaultMsgValidationError, "invalid product id", http.StatusUnprocessableEntity)
	ErrCustomInvalidCartId     = NewError(DefaultMsgValidationError, "invalid cart id", http.StatusUnprocessableEntity)

	ErrCustomUnauthorized      = NewError(DefaultMsgUnauthorized, "authentication required", http.StatusUnauthorized)
	ErrCustomTokenEmpty        = NewError(DefaultMsgUnauthorized, "token is required", http.StatusUnauthorized)
	ErrCustomInvalidToken      = NewError(DefaultMsgUnauthorized, "invalid token or expired", http.StatusUnauthorized)
	ErrCustomInvalidCredential = NewError(DefaultMsgUnauthorized, "invalid email or password", http.StatusUnauthorized)

	ErrCustomForbidden = NewError(DefaultMsgForbidden, "access to this resource is denied", http.StatusForbidden)

	ErrCustomEmailNotFound    = NewError(DefaultMsgNotFound, "email not found", http.StatusNotFound)
	ErrCustomProductNotFound  = NewError(DefaultMsgNotFound, "product not found", http.StatusNotFound)
	ErrCustomCategoryNotFound = NewError(DefaultMsgNotFound, "category not found", http.StatusNotFound)
	ErrCustomRoleNotFound     = NewError(DefaultMsgNotFound, "role not found", http.StatusNotFound)
	ErrCustomAddressNotFound  = NewError(DefaultMsgNotFound, "address not found", http.StatusNotFound)
	ErrCustomUserNotFound     = NewError(DefaultMsgNotFound, "user not found", http.StatusNotFound)

	ErrCustomEmailAlreadyExist = NewError(DefaultMsgConflict, "email already exists", http.StatusConflict)
)
