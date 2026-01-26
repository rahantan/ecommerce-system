package models

import "time"

type ProductModel struct {
	ID         int64         `gorm:"column:id;primaryKey;autoIncrement"`
	Name       string        `gorm:"column:name"`
	Price      float64       `gorm:"column:price"`
	Stock      int           `gorm:"column:stock"`
	CreatedAt  time.Time     `gorm:"column:created_at;autoCreateTime"`
	UpdatedAt  time.Time     `gorm:"column:updated_at;autoCreateTime;autoUpdateTime"`
	CategoryID int64         `gorm:"column:category_id"`
	Category   CategoryModel `gorm:"foreignKey:CategoryID;references:ID"`
}

func (p *ProductModel) TableName() string {
	return "products"
}
