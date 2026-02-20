package response

type ResCartItem struct {
	CartID   int64      `json:"cart_id"`
	Product  ResProduct `json:"product"`
	Qty      int        `json:"qty"`
	SubTotal int64      `json:"sub_total"`
}
