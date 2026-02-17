package response

type ResPayment struct {
	ID          int64  `json:"id"`
	Token       string `json:"snap_token"`
	OrderID     int64  `json:"order_id"`
	RedirectUrl string `json:"redirect_url"`
}
