package models

type User struct {
	UserID   string    `gorm:"column:user_id;primaryKey"`
	Email    string    `gorm:"column:email;unique"`
	Name     string    `gorm:"column:name"`
	Password string    `gorm:"column:password"`
	Accounts []Account `gorm:"foreignKey:UserID"`
}
