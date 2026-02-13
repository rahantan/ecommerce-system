package model

type RoleModels struct {
	ID    int64  `gorm:"column:id;pribrayKey;autoIncrement"`
	Title string `gorm:"column:title"`
}

func (r *RoleModels) TableName() string {
	return "roles"
}
