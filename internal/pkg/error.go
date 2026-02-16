package pkg

import (
	"ecommerce-system/internal/domain/model"
	"errors"
	"fmt"
	"strings"

	"github.com/go-playground/validator/v10"
)

func ValidationError(err error) *ErrorCustom {
	if validationErr, ok := err.(validator.ValidationErrors); ok {
		errDetail := make(map[string]string)
		for _, validate := range validationErr {
			switch validate.Tag() {
			case "email":
				errDetail[validate.Field()] = strings.ToLower(fmt.Sprint("invalid email"))
			case "required":
				errDetail[validate.Field()] = strings.ToLower(fmt.Sprintf("%s is required", validate.Field()))
			case "min":
				errDetail[validate.Field()] = strings.ToLower(fmt.Sprintf("%s must be at least %s characters", validate.Field(), validate.Param()))
			case "max":
				errDetail[validate.Field()] = strings.ToLower(fmt.Sprintf("%s must be at most %s characters", validate.Field(), validate.Param()))
			case "eqfield":
				errDetail[validate.Field()] = strings.ToLower(fmt.Sprintf("%s doesn't match", validate.Field()))
			case "numeric":
				errDetail[validate.Field()] = strings.ToLower(fmt.Sprintf("%s must be number", validate.Field()))
			case "gt":
				errDetail[validate.Field()] = strings.ToLower(fmt.Sprintf("%s must be at least 1", validate.Field()))
			case "oneof":
				errDetail[validate.Field()] = fmt.Sprintf("%s must be one of [%s]", strings.ToLower(validate.Field()), validate.Param())
			case "required_if":
				errDetail[validate.Field()] = fmt.Sprintf("%s is required when %s is %s", strings.ToLower(validate.Field()), strings.ToLower(validate.Param()[0:strings.Index(validate.Param(), " ")]), validate.Param()[strings.Index(validate.Param(), " ")+1:])
			}
		}
		return NewError(KindValidationError, "validation errors", errDetail)
	}
	return nil
}
func MappingError(err error) error {

	switch {
	case errors.Is(err, model.ErrDuplicateEmail):
		return NewError(KindConflict, err.Error(), nil)

	case errors.Is(err, model.ErrUserNotFound),
		errors.Is(err, model.ErrCategoryNotFound),
		errors.Is(err, model.ErrProductNotFound),
		errors.Is(err, model.ErrRoleNotFound),
		errors.Is(err, model.ErrAddressNotFound),
		errors.Is(err, model.ErrCartItemNotFound),
		errors.Is(err, model.ErrCheckOutNotFound):

		return NewError(KindNotFound, err.Error(), nil)

	default:
		return err
	}

}
