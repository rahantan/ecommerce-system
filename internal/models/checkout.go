package models

import (
	"time"
)

type CheckoutModel struct {
	ID           int64               `gorm:"column:id;primaryKey;autoIncrement"`
	TotalPrice   float64             `gorm:"column:total_price"`
	Status       string              `gorm:"column:status"`
	Source       string              `gorm:"column:source"`
	UserID       int64               `gorm:"column:user_id"`
	CheckoutItem []CheckoutItemModel `gorm:"foreignKey:CheckoutID"`
	CreatedAt    time.Time           `gorm:"column:created_at;autoCreateTime"`
	UpdatedAt    time.Time           `gorm:"column:updated_at;autoUpdateTime"`
}

func (om *CheckoutModel) TableName() string {
	return "checkout_sessions"
}

type CheckoutItemModel struct {
	ID         int64        `gorm:"column:id;primaryKey;autoIncrement"`
	CheckoutID int64        `gorm:"column:checkout_id"`
	ProductID  int64        `gorm:"column:product_id"`
	Qty        int          `gorm:"column:qty"`
	Price      float64      `gorm:"column:price"`
	SubTotal   float64      `gorm:"column:subtotal"`
	CartID     int64        `gorm:"column:cart_item_id"` //required in cart
	Product    ProductModel `gorm:"foreignKey:ProductID"`
}

func (oim *CheckoutItemModel) TableName() string {
	return "checkout_items"
}
