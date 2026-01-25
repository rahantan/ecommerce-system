package utils

// import (
// 	"ecommerce-system/internal/exceptions"
// 	"errors"

// 	"github.com/go-sql-driver/mysql"
// )

// var (
// 	ErrUserNotFound   = errors.New("user not found")
// 	ErrDuplicateEmail = errors.New("duplicate email")
// 	ErrNoRowsAffected = errors.New("no rows affected")
// )

// func IsDuplicateKeyError(err error) bool {
// 	var mysqlErr *mysql.MySQLError
// 	if errors.As(err, &mysqlErr) {
// 		return mysqlErr.Number == 1062
// 	}
// 	return false
// }

// func CheckError(err error) error {
// 	if err == nil {
// 		return nil
// 	}
// 	switch {

// 	// USER SERVICE CHECK
// 	case errors.Is(err, exceptions.ErrUnauthorized):
// 		return exceptions.ErrUnauthorized
// 	case errors.Is(err, ErrUserNotFound):
// 		return exceptions.ErrUserNotFound
// 	case errors.Is(err, ErrDuplicateEmail):
// 		return exceptions.ErrEmailAlreadyExist
// 	case errors.Is(err, ErrDuplicateEmail):
// 		return exceptions.ErrEmailAlreadyExist

// 	default:
// 		return exceptions.ErrInternalServer
// 	}
// }
