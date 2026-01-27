package exceptions

import (
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/go-sql-driver/mysql"
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
	case errors.Is(err, ErrUserNotFound):
		return ErrCustomUserNotFound
	case errors.Is(err, ErrDuplicateEmail):
		return ErrCustomEmailAlreadyExist
	case errors.Is(err, ErrCategoryNotFound):
		return ErrCustomCategoryNotFound
	case errors.Is(err, ErrProductNotFound):
		return ErrCustomProductNotFound
	case errors.Is(err, ErrRoleNotFound):
		return ErrCustomRoleNotFound
	case errors.Is(err, ErrAddressNotFound):
		return ErrCustomAddressNotFound
	default:
		return NewError(DefaultMsgInternalErr, err.Error(), http.StatusInternalServerError)
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
			}
		}
		return NewError(DefaultMsgValidationError, errDetail, http.StatusUnprocessableEntity)
	}
	return nil
}

func IsDuplicateKeyError(err error) bool {
	var mysqlErr *mysql.MySQLError

	if errors.As(err, &mysqlErr) {
		return mysqlErr.Number == 1062
	}

	return false
}

func CheckContainError(err error) bool {
	if strings.Contains(err.Error(), "Error 1452") {
		return true
	}
	return false
}
