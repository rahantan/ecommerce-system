package utils

import (
	"ecommerce-system/internal/exceptions"
	"ecommerce-system/internal/models"
	"errors"
)

func MappingError(err error) error {

	switch {
	case errors.Is(err, models.ErrDuplicateEmail):
		return exceptions.NewError(exceptions.KindConflict, err.Error(), nil)

	case errors.Is(err, models.ErrUserNotFound),
		errors.Is(err, models.ErrCategoryNotFound),
		errors.Is(err, models.ErrProductNotFound),
		errors.Is(err, models.ErrRoleNotFound),
		errors.Is(err, models.ErrAddressNotFound),
		errors.Is(err, models.ErrCartItemNotFound),
		errors.Is(err, models.ErrCheckOutNotFound):

		return exceptions.NewError(exceptions.KindNotFound, err.Error(), nil)

	default:
		return err
	}

}
