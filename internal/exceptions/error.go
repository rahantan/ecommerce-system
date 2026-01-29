package exceptions

import (
	"ecommerce-system/internal/models"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/go-playground/validator/v10"
)

// func getNewMsg(messages ...string) string {
// 	msg := ""
// 	for _, msgVal := range messages {
// 		msg += " " + msgVal
// 	}
// 	return msg
// }

// Check Error Service Layer
func CheckError(err error) error {
	if err == nil {
		return nil
	}
	switch {
	case errors.Is(err, models.ErrUserNotFound):
		return ErrCustomUserNotFound
	case errors.Is(err, models.ErrDuplicateEmail):
		return ErrCustomEmailAlreadyExist
	case errors.Is(err, models.ErrCategoryNotFound):
		return ErrCustomCategoryNotFound
	case errors.Is(err, models.ErrProductNotFound):
		return ErrCustomProductNotFound
	case errors.Is(err, models.ErrRoleNotFound):
		return ErrCustomRoleNotFound
	case errors.Is(err, models.ErrAddressNotFound):
		return ErrCustomAddressNotFound
	case errors.Is(err, models.ErrCartItemNotFound):
		return NewError("", "not found cart item", http.StatusFound)
	default:
		return err
	}
}
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
			}
		}
		return NewError(DefaultMsgValidationError, errDetail, http.StatusUnprocessableEntity)
	}
	return nil
}
