package utils

import (
	"ecommerce-system/internal/exceptions"
)

func WithMessage(err error, msg string) error {
	if newErr, ok := err.(*exceptions.ErrorCustom); ok {
		return newErr.WithMessage(msg)
	}
	return err
}
