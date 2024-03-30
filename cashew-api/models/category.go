package models

type Category struct {
	CategoryID   string        `gorm:"column:category_id;primaryKey"`
	Transactions []Transaction `gorm:"foreignKey:CategoryID"`
}
