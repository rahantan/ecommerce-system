package response

type ResItem struct {
	ProductID int64   `json:"product_id"`
	Qty       int     `json:"qty"`
	Price     float64 `json:"price"`
	SubTotal  float64 `json:"sub_total"`
}

type ResCheckOut struct {
	ID         int64     `json:"id"`
	Status     string    `json:"status"`
	Items      []ResItem `json:"items"`
	TotalPrice float64   `json:"total_price"`
	CreatedAt  string    `json:"created_at"`
}
