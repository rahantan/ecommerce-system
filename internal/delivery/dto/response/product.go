package response

type ResProduct struct {
	ID           int64  `json:"product_id"`
	Name         string `json:"name"`
	Price        int64  `json:"price"`
	Stock        int    `json:"stock"`
	ImageUrl     string `json:"image_url,omitempty"`
	CreatedAt    string `json:"created_at,omitempty"`
	UpdatedAt    string `json:"updated_at,omitempty"`
	*ResCategory `json:"category,omitempty"`
}
