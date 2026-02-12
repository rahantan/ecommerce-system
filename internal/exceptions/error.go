package exceptions

import (
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
			}
		}
		return NewError(KindValidationError, "validation errors", errDetail)
	}
	return nil
}
