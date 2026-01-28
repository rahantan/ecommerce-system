package request

type ReqCreateAddress struct {
	UserID   int64
	City     string `json:"city" validate:"required,min=3"`
	Address  string `json:"address" validate:"required,min=3"`
	IsActive *bool  `json:"is_active" validate:"required"`
}
type ReqUpdateAddress struct {
	ID       int64 `json:"id" validate:"required,numeric"`
	UserID   int64
	City     string `json:"city" validate:"required,min=3"`
	Address  string `json:"address" validate:"required,min=3"`
	IsActive *bool  `json:"is_active" validate:"required"`
}
