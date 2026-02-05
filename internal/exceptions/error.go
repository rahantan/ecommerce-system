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
		return NewError(DefaultMsgNotFound, err.Error(), http.StatusNotFound)
	case errors.Is(err, models.ErrDuplicateEmail):
		return NewError(DefaultMsgConflict, err.Error(), http.StatusConflict)
	case errors.Is(err, models.ErrCategoryNotFound):
		return NewError(DefaultMsgNotFound, err.Error(), http.StatusNotFound)
	case errors.Is(err, models.ErrProductNotFound):
		return NewError(DefaultMsgNotFound, err.Error(), http.StatusNotFound)
	case errors.Is(err, models.ErrRoleNotFound):
		return NewError(DefaultMsgNotFound, err.Error(), http.StatusNotFound)
	case errors.Is(err, models.ErrAddressNotFound):
		return NewError(DefaultMsgNotFound, err.Error(), http.StatusNotFound)
	case errors.Is(err, models.ErrCartItemNotFound):
		return NewError(DefaultMsgNotFound, err.Error(), http.StatusNotFound)
	case errors.Is(err, models.ErrCheckOutNotFound):
		return NewError(DefaultMsgNotFound, err.Error(), http.StatusNotFound) //ini bakalan nil

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
