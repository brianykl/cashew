package models

import (
	"github.com/google/uuid"
)

type User struct {
	UserID   string `gorm:"column:user_id;primaryKey"`
	Email    string `gorm:"column:email;unique"`
	Name     string `gorm:"column:name"`
	Password string `gorm:"column:password"`
}

func NewUser(email, name, password string) (*User, error) {
	return &User{
		UserID:   generateUserID(),
		Email:    email,
		Name:     name,
		Password: password,
	}, nil
}

// func NewUser(userID, email, name, password string, key []byte) (*User, error) {
// 	// encrypt email
// 	encryptedEmail, err := crypto.Encrypt(email, key)
// 	if err != nil {
// 		return nil, err // handle error for email encryption
// 	}

// 	// encrypt name
// 	encryptedName, err := crypto.Encrypt(name, key)
// 	if err != nil {
// 		return nil, err // handle error for name encryption
// 	}

// 	// proceed to create the User object if no errors occurred
// 	return &User{
// 		UserID:   userID,
// 		Email:    encryptedEmail,
// 		Name:     encryptedName,
// 		Password: password,
// 	}, nil
// }

func ValidateUser(user *User) error {
	return nil
}

func VerifyPassword(hashed_password, provided_password string) bool {
	return true
}

func FindUserByEmail(email string) (*User, error) {
	return nil, nil
}

func UpdateUser(user *User) error {
	return nil
}

func ChangeUserPassword(userID, new_password string) error {
	return nil
}

func DeleteUser(userID string) error {
	return nil
}

func generateUserID() string {
	id := uuid.New()
	return id.String()
}
