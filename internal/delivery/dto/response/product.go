package response

type ResProduct struct {
	ID           int64  `json:"product_id"`
	Name         string `json:"name"`
	Price        int64  `json:"price,omitempty"`
	Stock        int    `json:"stock,omitempty"`
	CreatedAt    string `json:"created_at,omitempty"`
	UpdatedAt    string `json:"updated_at,omitempty"`
	*ResCategory `json:"category,omitempty"`
}
