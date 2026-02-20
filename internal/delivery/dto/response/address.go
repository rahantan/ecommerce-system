package response

type ResAddress struct {
	ID       int64  `json:"address_id"`
	City     string `json:"city"`
	Address  string `json:"address"`
	IsActive bool   `json:"is_active"`
}
