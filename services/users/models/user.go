package models

import (
	"context"
	"encoding/hex"
	"fmt"
	"log"

	"github.com/brianykl/cashew/services/crypto/client"
	"github.com/google/uuid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type User struct {
	UserID   string `gorm:"column:user_id;primaryKey"`
	Email    string `gorm:"column:email;unique"`
	Name     string `gorm:"column:name"`
	Password string `gorm:"column:password"`
}

func NewUser(email, name, password string) (*User, error) {
	cryptoServiceAddress := "localhost:5002"
	conn, err := grpc.NewClient(cryptoServiceAddress, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("could not connect to crypto service: %v", err)
	}
	defer conn.Close()
	ctx := context.Background()
	cryptoClient := client.NewCryptoClient(conn)

	params := client.Argon2IdParams{
		Memory:      64 * 1024, // 64 MiB of RAM
		Iterations:  3,         // 3 Iterations
		Parallelism: 2,         // Utilize 2 CPU cores (adjust if needed)
		SaltLength:  16,        // 16-byte salt
		KeyLength:   32,        // 32-byte output hash
	}

	// example encryption key, going to figure out how to securely generate and store it
	hexKey := "f13fd7ee2c6346b67aae8863ec68c170d26766a6fe216485ca5bfdfa1c25b233"

	key, err := hex.DecodeString(hexKey)
	if err != nil {
		log.Printf("error decoding hex key: %v", err)
		return nil, fmt.Errorf("error decoding hex key: %v", err)
	}

	hashedEmail, err := cryptoClient.HashPII(ctx, email, key) // THIS NEEDS TO BE CHANGED TO HMAC
	if err != nil {
		log.Printf("error encrypting email: %v", err)
		return nil, fmt.Errorf("error encrypting email: %v", err)
	}

	hashedName, err := cryptoClient.HashPII(ctx, name, key) // THIS NEEDS TO BE CHANGED TO HMAC
	if err != nil {
		log.Printf("error encrypting name: %v", err)
		return nil, fmt.Errorf("error encrypting name: %v", err)
	}
	encodedHash, _ := cryptoClient.HashPassword(ctx, password, &params)
	if err != nil {
		return nil, fmt.Errorf("error hashing password: %v", err)
	}

	return &User{
		UserID:   generateUserID(),
		Email:    hashedEmail.EncodedHash,
		Name:     hashedName.EncodedHash,
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
	// decryptedEmail := user.Email
	// go search db for email
	// if email does not exist, return invalid credentials
	// verify provided password against hashed password
	// if password does not exist, return invalid credentials
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
