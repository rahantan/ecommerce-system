package order

import (
	"crypto/sha512"
	"encoding/hex"
)

func (orderHandler *OrderHandlerImpl) midTransVerify(orderId, statusCode, grossMount, signature string) bool {
	payLoad := orderId + statusCode + grossMount + orderHandler.Midtrans.ServerKey
	hash := sha512.Sum512([]byte(payLoad))
	return hex.EncodeToString(hash[:]) == signature
}
