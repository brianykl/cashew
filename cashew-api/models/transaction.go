package models

type Transaction struct {
	TransactionID string  `gorm:"column:transaction_id;primaryKey"`
	AccountID     string  `gorm:"column:user_id"`
	CategoryID    string  `gorm:"column:category_id"`
	VendorName    string  `gorm:"column:vendor_name"`
	Amount        float32 `gorm:"column:amount"`
	Account       Account `gorm:"foreignKey:AccountID;references:AccountID"`
}
