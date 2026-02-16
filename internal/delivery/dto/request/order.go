package request

type ReqItem struct {
	ProductID int64 `json:"product_id"`
	Qty       int   `json:"qty"`
}

type ReqCheckout struct {
	Source  string    `json:"source" validate:"required,oneof=cart direct"`
	Items   []ReqItem `json:"items,omitempty" validate:"required_if=Source direct,dive"`
	CartIDs []int64   `json:"cart_ids,omitempty" validate:"required_if=Source cart,dive"`
}

type ReqConfirmCheckout struct {
	// CheckoutID    int64  `json:"checkout_id" validate:"required"`
	Note          string `json:"note"`
	PaymentMethod string `json:"payment_method" validate:"required"`
	AddressID     int64  `json:"address_id,omitempty"` //default address active
}

/*

























order
- paymentmetod
order item
- productId
- orderID
- qty


*/

// type ReqCreateOrder struct {
// 	Noted         string     `json:"noted"`
// 	PaymentMethod string     `json:"payment_method"`
// 	Items         []*ReqItem `json:"items"`
// }
