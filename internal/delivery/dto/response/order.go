package response

type ResItem struct {
	ProductID int64 `json:"product_id"`
	Qty       int   `json:"qty"`
	Price     int64 `json:"price"`
	SubTotal  int64 `json:"sub_total"`
}

type ResCheckOut struct {
	ID         int64     `json:"id"`
	Status     string    `json:"status"`
	Items      []ResItem `json:"items"`
	TotalPrice int64     `json:"total_price"`
	CreatedAt  string    `json:"created_at"`
}
type ResOrderStatus struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}

type ResOrderProduct struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}
type ResOrderItem struct {
	ID       int64           `json:"id"`
	Qty      int             `json:"qty"`
	Price    int64           `json:"price"`
	SubTotal int64           `json:"sub_total"`
	Product  ResOrderProduct `json:"product"`
}

type ResOrder struct {
	ID         int64           `json:"id"`
	TotalPrice int64           `json:"total_price"`
	CreatedAt  string          `json:"created_at"`
	TotalItems int             `json:"total_items,omitempty"`
	Items      *[]ResOrderItem `json:"items,omitempty"`
	Status     ResOrderStatus  `json:"status"`
}
