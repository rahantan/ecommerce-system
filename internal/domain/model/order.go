package model

import (
	"time"

	"gorm.io/gorm"
)

type OrderStatusModel struct {
	ID   int64  `gorm:"column:id;primaryKey;autoIncrement"`
	Name string `gorm:"column:name"`
}

func (osm *OrderStatusModel) TableName() string {
	return "order_status"
}

type OrderModel struct {
	ID            int64             `gorm:"column:id;primaryKey;autoIncrement"`
	UserID        int64             `gorm:"column:user_id"`
	TotalPrice    int64             `gorm:"column:total_price"`
	PaymentMethod string            `gorm:"column:payment_method"`
	Noted         string            `gorm:"column:note"`
	StatusID      int64             `gorm:"column:status_id"`
	CreatedAt     time.Time         `gorm:"column:created_at;autoCreateTime"`
	UpdatedAt     time.Time         `gorm:"column:updated_at;autoUpdateTime"`
	DeletedAt     gorm.DeletedAt    `gorm:"column:deleted_at;index"`
	OrderStatus   OrderStatusModel  `gorm:"foreignKey:StatusID;references:ID"`
	OrderItem     []*OrderItemModel `gorm:"foreignKey:OrderID"`
	AddressOrder  AddressOrderModel `gorm:"foreignKey:OrderID"`
	Payment       PaymentOrderModel `gorm:"foreignKey:OrderID"`
}

func (om *OrderModel) TableName() string {
	return "orders"
}

type OrderItemModel struct {
	ID        int64        `gorm:"column:id;primaryKey;autoIncrement"`
	OrderID   int64        `gorm:"column:order_id"`
	ProductID int64        `gorm:"column:product_id"`
	Qty       int          `gorm:"column:qty"`
	Price     int64        `gorm:"column:price"`
	SubTotal  int64        `gorm:"column:subtotal"`
	Product   ProductModel `gorm:"foreignKey:ProductID"`
}

func (oim *OrderItemModel) TableName() string {
	return "order_items"
}

type AddressOrderModel struct {
	ID        int64     `gorm:"column:id;primaryKey;autoIncrement"`
	OrderID   int64     `gorm:"column:order_id"`
	City      string    `gorm:"column:city"`
	Address   string    `gorm:"column:address"`
	CreatedAt time.Time `gorm:"column:created_at;autoCreateTime"`
}

func (a *AddressOrderModel) TableName() string {
	return "address_orders"
}

type PaymentOrderModel struct {
	ID          int64     `gorm:"column:id;primaryKey;autoIncrement"`
	OrderID     int64     `gorm:"column:order_id"`
	SnapToken   string    `gorm:"column:snap_token"`
	RedirectURL string    `gorm:"column:redirect_url"`
	Status      string    `gorm:"column:status"`
	CreatedAt   time.Time `gorm:"column:created_at;autoCreateTime"`
	UpdatedAt   time.Time `gorm:"column:updated_at;autoUpdateTime"`
}

func (b *PaymentOrderModel) TableName() string {
	return "payments"
}
