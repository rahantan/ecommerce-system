package request

type ReqCreateUser struct {
	Name            string `json:"name" validate:"required"`
	Email           string `json:"email" validate:"required,email"`
	Password        string `json:"password" validate:"required,min=8"`
	ConfirmPassword string `json:"confirm_password" validate:"required,min=8,eqfield=Password"`
	Phone           string `json:"phone" validate:"required"`
	RoleID          int64  `json:"role_id" validate:"required"`
}
