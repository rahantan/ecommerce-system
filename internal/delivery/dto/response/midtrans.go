package response

type ResPayment struct {
	Token       string `json:"snap_token"`
	OrderID     int64  `json:"order_id"`
	RedirectUrl string `json:"redirect_url"`
}
