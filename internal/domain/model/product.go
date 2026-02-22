package model

import (
	"time"

	"gorm.io/gorm"
)

type ProductModel struct {
	ID         int64          `gorm:"column:id;primaryKey;autoIncrement"`
	Name       string         `gorm:"column:name"`
	Price      int64          `gorm:"column:price"`
	Stock      int            `gorm:"column:stock"`
	Image      string         `gorm:"column:image"`
	CreatedAt  time.Time      `gorm:"column:created_at;autoCreateTime"`
	UpdatedAt  time.Time      `gorm:"column:updated_at;autoCreateTime;autoUpdateTime"`
	DeletedAt  gorm.DeletedAt `gorm:"column:deleted_at;index"`
	CategoryID int64          `gorm:"column:category_id"`
	Category   CategoryModel  `gorm:"foreignKey:CategoryID;references:ID"`
}

func (p *ProductModel) TableName() string {
	return "products"
}
