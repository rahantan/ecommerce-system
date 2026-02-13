package pkg

import "net/http"

type ErrorCustom struct {
	Kind    Kind
	Message string
	Errors  any
}

func NewError(kind Kind, message string, errs any) *ErrorCustom {
	return &ErrorCustom{
		Kind:    kind,
		Message: message,
		Errors:  errs,
	}
}

func (err *ErrorCustom) Error() string {
	return string(err.Message)
}
func (err *ErrorCustom) GetStatusCode() int {

	switch err.Kind {
	case KindValidationError:
		return http.StatusBadRequest
	case KindCancelled:
		return http.StatusUnprocessableEntity
	case KindBadRequest:
		return http.StatusBadRequest
	case KindUnauthorized:
		return http.StatusUnauthorized
	case KindForbidden:
		return http.StatusForbidden
	case KindNotFound:
		return http.StatusNotFound
	case KindConflict:
		return http.StatusConflict
	default:
		return http.StatusInternalServerError
	}
}

type Kind string

const (
	KindInternal        Kind = "internal server error"
	KindBadRequest      Kind = "bad request"
	KindValidationError Kind = "validation error"
	KindUnauthorized    Kind = "unauthorized"
	KindForbidden       Kind = "access denied"
	KindNotFound        Kind = "not found"
	KindConflict        Kind = "conflict"
	KindCancelled       Kind = "cancelled"
)

// CUSTOM ERROR
var (
	ErrCustomUnauthorized   = NewError(KindUnauthorized, "authentication required", nil)
	ErrCustomForbidden      = NewError(KindForbidden, "access to this resource is denied", nil)
	ErrCustomInvalidPayload = NewError(KindBadRequest, "invalid payload", nil)

	//PARAMS ERR
	ErrCustomInvalidCategoryId = NewError(KindValidationError, "invalid category id", nil)
	ErrCustomInvalidAddressId  = NewError(KindValidationError, "invalid address id", nil)
	ErrCustomInvalidProductId  = NewError(KindValidationError, "invalid product id", nil)
	ErrCustomInvalidCartId     = NewError(KindValidationError, "invalid cart id", nil)

	ErrCustomLogin = NewError(KindUnauthorized, "invalid email or password", nil)
)
