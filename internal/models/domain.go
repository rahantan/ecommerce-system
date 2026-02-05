package models

import (
	"errors"

	"github.com/go-sql-driver/mysql"
)

var (
	ErrUserNotFound   = errors.New("user not found")
	ErrDuplicateEmail = errors.New("email already exist")
	ErrNoRowsAffected = errors.New("no rows affected")

	ErrProductNotFound  = errors.New("product not found")
	ErrCategoryNotFound = errors.New("category not found")
	ErrCartItemNotFound = errors.New("cart item not found")
	ErrRoleNotFound     = errors.New("role not found")
	ErrAddressNotFound  = errors.New("address not found")
	ErrCheckOutNotFound = errors.New("checkout not found")
)

func getMysqlErr(err error, mysqlCode uint16) bool {
	var mysqlErr *mysql.MySQLError
	if errors.As(err, &mysqlErr) {
		return mysqlErr.Number == mysqlCode
	}
	return false
}
func IsDuplicateKeyError(err error) bool {
	return getMysqlErr(err, 1062)
}

func IsInternalErrMysql(err error) bool {
	MySQLErrorMap := map[uint16]string{
		//  Unique & Primary
		1169: "Unique constraint violation",

		//  Data validation
		1048: "Column cannot be null",
		1364: "Field doesn't have a default value",
		1265: "Data truncated for column",
		1292: "Incorrect value for column",

		//  Query & schema
		1064: "SQL syntax error",
		1054: "Unknown column",
		1146: "Table doesn't exist",

		//  Auth & connection
		1045: "Access denied for database user",
		1049: "Unknown database",
		2002: "Cannot connect to MySQL server",

		//  Transaction & lock
		1213: "Deadlock found when trying to get lock",
		1205: "Lock wait timeout exceeded",
	}

	var mysqlErr *mysql.MySQLError
	if errors.As(err, &mysqlErr) {
		return MySQLErrorMap[mysqlErr.Number] != ""
	}
	return false
}

func ForeignKeyErr(err error) bool {
	return getMysqlErr(err, 1452)
}
