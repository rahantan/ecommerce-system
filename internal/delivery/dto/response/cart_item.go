package response

type ResCartItem struct {
	CartItemID int64      `json:"cart_item_id"`
	Product    ResProduct `json:"product"`
	Qty        int        `json:"qty"`
	SubTotal   int64      `json:"sub_total"`
}
