package models

import "time"

// type User struct {
// }
type UserModel struct {
	ID        int64          `gorm:"column:id;primaryKey;autoIncrement"`
	Name      string         `gorm:"column:name"`
	Email     string         `gorm:"column:email;unique"`
	Password  string         `gorm:"column:password"`
	Phone     string         `gorm:"column:phone"`
	CreatedAt time.Time      `gorm:"column:created_at;autoCreateTime"`
	UpdatedAt time.Time      `gorm:"column:updated_at;autoCreateTime;autoUpdateTime"`
	RoleID    int64          `gorm:"column:role_id"`
	Address   []AddressModel `gorm:"foreignKey:UserID;references:ID"`
	Role      RoleModels     `gorm:"foreignKey:RoleID;references:ID"`
}

func (u *UserModel) TableName() string {
	return "users"
}
