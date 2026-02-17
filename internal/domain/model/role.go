package model

type RoleModels struct {
	ID    int64  `gorm:"column:id;primaryKey;autoIncrement"`
	Title string `gorm:"column:title"`
}

func (r *RoleModels) TableName() string {
	return "roles"
}
