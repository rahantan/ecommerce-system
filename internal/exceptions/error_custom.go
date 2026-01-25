package exceptions

import (
	"errors"
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

const (
	MsgBadRequest      = "bad request"
	MsgValidationError = "validation error"
	MsgUnauthorized    = "unauthorized"
	MsgForbidden       = "access denied"
	MsgNotFound        = "not found"
	MsgConflict        = "conflict"
	MsgInternalErr     = "internal server error"
)

// CUSTOM ERROR
var (
	ErrCustomUnauthorized = NewError(MsgUnauthorized, "authentication required", http.StatusUnauthorized)
	ErrCustomTokenEmpty   = NewError(MsgUnauthorized, "token is required", http.StatusUnauthorized)
	ErrCustomInvalidToken = NewError(MsgUnauthorized, "invalid token", http.StatusUnauthorized)

	ErrCustomForbidden      = NewError(MsgForbidden, nil, http.StatusForbidden)
	ErrCustomInvalidPayload = NewError(MsgBadRequest, "invalid payload", http.StatusBadRequest)

	ErrCustomInvalidCredential = NewError(MsgUnauthorized, "invalid email or password", http.StatusUnauthorized)

	ErrCustomEmailNotFound = NewError(MsgNotFound, "email not found", http.StatusNotFound)

	ErrCustomUserNotFound = NewError(MsgNotFound, "user not found", http.StatusNotFound)

	ErrCustomEmailAlreadyExist = NewError(MsgConflict, "email already exists", http.StatusConflict)

	ErrCustomValidation = NewError(MsgValidationError, nil, http.StatusUnprocessableEntity)
)

var (
	ErrUserNotFound   = errors.New("user not found")
	ErrDuplicateEmail = errors.New("duplicate email")
	ErrNoRowsAffected = errors.New("no rows affected")
)
