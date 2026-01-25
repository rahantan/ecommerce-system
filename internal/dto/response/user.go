package response

type ResUser struct {
	ID        int64     `json:"id"`
	Name      string    `json:"name" `
	Email     string    `json:"email"`
	Phone     string    `json:"phone" `
	CreatedAt string    `json:"created_at"`
	UpdatedAt string    `json:"updated_at"`
	Role      Role      `json:"role"`
	Addresses []Address `json:"addresses"`
}

// type ResUser struct {
// 	User user `json:"embed"`
// }
