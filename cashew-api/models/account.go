package models

type Account struct {
	AccountID   string `gorm:"column:account_id;primaryKey"`
	UserID      string `gorm:"column:user_id"`
	AccountType string `gorm:"column:account_type"`
	User        User   `gorm:"foreignKey:UserID;references:UserID"`
}
