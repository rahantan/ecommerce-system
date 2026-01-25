package models

type Address struct {
	ID      int64  `gorm:"column:id;pribrayKey;autoIncrement"`
	UserID  string `gorm:"column:user_id"`
	City    string `gorm:"column:city"`
	Address string `gorm:"column:address"`
}
