package models

import (
	"context"
	"time"

	"github.com/brianykl/cashew/services/crypto/client"
	"github.com/google/uuid"
	"google.golang.org/grpc"
)

type User struct {
	UserID   string `gorm:"column:user_id;primaryKey"`
	Email    string `gorm:"column:email;unique"`
	Name     string `gorm:"column:name"`
	Password string `gorm:"column:password"`
}

func NewUser(email, name, password string) (*User, error) {
	cryptoServiceAddress := "localhost:50051"
	conn, err := grpc.Dial(cryptoServiceAddress, grpc.WithInsecure())
	if err != nil {
		// Handle connection errors
	}
	defer conn.Close()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	cryptoClient := client.NewCryptoClient(conn)

	params := client.Argon2IdParams{
		Memory:      64 * 1024, // 64 MiB of RAM
		Iterations:  3,         // 3 Iterations
		Parallelism: 2,         // Utilize 2 CPU cores (adjust if needed)
		SaltLength:  16,        // 16-byte salt
		KeyLength:   32,        // 32-byte output hash
	}

	encodedHash, _ := cryptoClient.HashPassword(ctx, password, &params)
	return &User{
		UserID:   generateUserID(),
		Email:    email,
		Name:     name,
		Password: encodedHash.EncodedHash,
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
