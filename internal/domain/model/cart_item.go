package model

import "time"

type CartItemModel struct {
	ID        int64        `gorm:"column:id;primaryKey;autoIncrement"`
	UserID    int64        `gorm:"column:user_id"`
	Qty       int          `gorm:"column:qty"`
	SubTotal  float64      `gorm:"column:subtotal"`
	CreatedAt time.Time    `gorm:"column:created_at;autoCreateTime"`
	UpdatedAt time.Time    `gorm:"column:updated_at;autoCreateTime;autoUpdateTime"`
	ProductID int64        `gorm:"column:product_id"`
	Product   ProductModel `gorm:"foreignKey:ProductID;references:ID"`
}

func (c *CartItemModel) TableName() string {
	return "cart_items"
}
