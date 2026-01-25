package exceptions

import (
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/go-sql-driver/mysql"
)

func CheckError(err error) error {
	if err == nil {
		return nil
	}
	switch {
	case errors.Is(err, ErrUserNotFound):
		return ErrCustomUserNotFound
	case errors.Is(err, ErrDuplicateEmail):
		return ErrCustomEmailAlreadyExist
	case errors.Is(err, ErrDuplicateEmail):
		return ErrCustomEmailAlreadyExist
	default:
		return NewError(MsgInternalErr, err.Error(), http.StatusInternalServerError)
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
			}
		}
		return NewError(MsgValidationError, errDetail, http.StatusUnprocessableEntity)
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
