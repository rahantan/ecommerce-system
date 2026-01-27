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
	DefaultMsgInternalErr     = "internal server error"
)

const (
	// Product
	MsgFailGetProduct    = "failed to get product"
	MsgFailGetAllProduct = "failed to get products"
	MsgFailCreateProduct = "failed to create product"
	MsgFailUpdateProduct = "failed to update product"

	// Category
	MsgFailGetCategory    = "failed to get category"
	MsgFailGetAllCategory = "failed to get categories"
	MsgFailCreateCategory = "failed to create category"
	MsgFailUpdateCategory = "failed to update category"
	MsgFailDeleteCategory = "failed to delete category"

	// Category
	MsgFailGetAddress      = "failed to get address"
	MsgFailGetAllAddresses = "failed to get addresses"
	MsgFailCreateAddress   = "failed to create address"
	MsgFailUpdateAddress   = "failed to update address"
	MsgFailDeleteAddress   = "failed to delete address"

	// User
	MsgFailGetUser    = "failed to get user"
	MsgFailGetAllUser = "failed to get users"
	MsgFailCreateUser = "failed to create user"
	MsgFailUpdateUser = "failed to update user"
	MsgFailDeleteUser = "failed to delete user"

	// Auth
	MsgFailRegister = "failed to register"
	MsgFailLogin    = "failed to login"
)

// CUSTOM ERROR
var (
	ErrCustomUnauthorized = NewError(DefaultMsgUnauthorized, "authentication required", http.StatusUnauthorized)
	ErrCustomTokenEmpty   = NewError(DefaultMsgUnauthorized, "token is required", http.StatusUnauthorized)
	ErrCustomInvalidToken = NewError(DefaultMsgUnauthorized, "invalid token", http.StatusUnauthorized)

	ErrCustomForbidden      = NewError(DefaultMsgForbidden, "access to this resource is denied", http.StatusForbidden)
	ErrCustomInvalidPayload = NewError(DefaultMsgBadRequest, "invalid payload", http.StatusBadRequest)

	ErrCustomInvalidCredential = NewError(DefaultMsgUnauthorized, "invalid email or password", http.StatusUnauthorized)

	ErrCustomEmailNotFound    = NewError(DefaultMsgNotFound, "email not found", http.StatusNotFound)
	ErrCustomProductNotFound  = NewError(DefaultMsgNotFound, "product not found", http.StatusNotFound)
	ErrCustomCategoryNotFound = NewError(DefaultMsgNotFound, "category not found", http.StatusNotFound)
	ErrCustomRoleNotFound     = NewError(DefaultMsgNotFound, "role not found", http.StatusNotFound)
	ErrCustomAddressNotFound  = NewError(DefaultMsgNotFound, "address not found", http.StatusNotFound)

	ErrCustomUserNotFound = NewError(DefaultMsgNotFound, "user not found", http.StatusNotFound)

	ErrCustomEmailAlreadyExist = NewError(DefaultMsgConflict, "email already exists", http.StatusConflict)

	ErrCustomValidation        = NewError(DefaultMsgValidationError, nil, http.StatusUnprocessableEntity)
	ErrCustomInvalidCategoryId = NewError(DefaultMsgValidationError, "invalid category id", http.StatusUnprocessableEntity)
	ErrCustomInvalidAddressId  = NewError(DefaultMsgValidationError, "invalid address id", http.StatusUnprocessableEntity)
	ErrCustomInvalidProductId  = NewError(DefaultMsgValidationError, "invalid product id", http.StatusUnprocessableEntity)
)

// repository layer
var (
	ErrUserNotFound   = errors.New("user not found")
	ErrDuplicateEmail = errors.New("duplicate email")
	ErrNoRowsAffected = errors.New("no rows affected")

	ErrProductNotFound  = errors.New("product not found")
	ErrCategoryNotFound = errors.New("category not found")
	ErrRoleNotFound     = errors.New("role not found")
	ErrAddressNotFound  = errors.New("address not found")
)
