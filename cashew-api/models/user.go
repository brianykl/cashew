package models

type User struct {
	UserID   string `gorm:"column:user_id"`
	Email    string `gorm:"column:email;unique"`
	Name     string `gorm:"column:name"`
	Password string `gorm:"column:password"`
}
