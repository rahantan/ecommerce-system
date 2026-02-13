package model

type AddressModel struct {
	ID       int64  `gorm:"column:id;primaryKey;autoIncrement"`
	UserID   int64  `gorm:"column:user_id"`
	City     string `gorm:"column:city"`
	IsActive bool   `gorm:"column:is_active"`
	Address  string `gorm:"column:address"`
}

func (a *AddressModel) TableName() string {
	return "addresses"
}
