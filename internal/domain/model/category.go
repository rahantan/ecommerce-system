package model

type CategoryModel struct {
	ID   int64  `gorm:"column:id;primaryKey;autoIncrement"`
	Name string `gorm:"column:name"`
}

func (c *CategoryModel) TableName() string {
	return "categories"
}
