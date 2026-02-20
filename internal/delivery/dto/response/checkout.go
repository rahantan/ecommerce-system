package response

type ResItem struct {
	ProductID int64 `json:"product_id"`
	Qty       int   `json:"qty"`
	Price     int64 `json:"price"`
	SubTotal  int64 `json:"sub_total"`
}

type ResCheckOut struct {
	ID         int64     `json:"checkout_id"`
	Status     string    `json:"status"`
	Items      []ResItem `json:"items"`
	TotalPrice int64     `json:"total_price"`
	CreatedAt  string    `json:"created_at"`
}
