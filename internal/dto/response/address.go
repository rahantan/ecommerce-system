package response

type ResAddress struct {
	ID       int64  `json:"id"`
	City     string `json:"city"`
	Address  string `json:"address"`
	IsActive bool   `json:"is_active"`
}
