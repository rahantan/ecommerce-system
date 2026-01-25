package response

type Address struct {
	ID      int64  `json:"id"`
	City    string `json:"city"`
	Address string `json:"address"`
}
