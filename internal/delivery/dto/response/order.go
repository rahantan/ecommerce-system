package response

type ResOrderStatus struct {
	ID   int64  `json:"status_id"`
	Name string `json:"name"`
}

type ResOrderProduct struct {
	ID   int64  `json:"product_id"`
	Name string `json:"name"`
}
type ResOrderItem struct {
	ID       int64           `json:"order_item_id"`
	Qty      int             `json:"qty"`
	Price    int64           `json:"price"`
	SubTotal int64           `json:"sub_total"`
	Product  ResOrderProduct `json:"product"`
}

type ResOrder struct {
	ID         int64           `json:"order_id"`
	TotalPrice int64           `json:"total_price"`
	CreatedAt  string          `json:"created_at"`
	TotalItems int             `json:"total_items,omitempty"`
	Items      *[]ResOrderItem `json:"items,omitempty"`
	Status     ResOrderStatus  `json:"status"`
}
