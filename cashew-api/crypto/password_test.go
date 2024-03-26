package crypto

import (
	"strings"
	"testing"
)

func TestHashAndPasswordVerification(t *testing.T) {
	password := "big-password"
	hashedPassword, err := HashPassword(password, &DefaultParams)
	if err != nil {
		t.Fatalf("Hashing the password failed: %v", err)
	}

	// Split the hashed password to verify format
	parts := strings.Split(hashedPassword, "$")
	if len(parts) != 2 {
		t.Errorf("Hashed password format incorrect, expected 2 parts got %d", len(parts))
	}

	// Test verifying the correct password
	if match, err := VerifyPassword(password, hashedPassword, &DefaultParams); err != nil || !match {
		t.Errorf("Password verification failed: %v", err)
	}

	// Test verifying an incorrect password
	if match, err := VerifyPassword("wrongpassword", hashedPassword, &DefaultParams); err != nil || match {
		t.Errorf("Password verification incorrectly succeeded for a wrong password")
	}
}

func TestGenerateSalt(t *testing.T) {
	salt, err := GenerateSalt(DefaultParams.SaltLength)
	if err != nil {
		t.Fatalf("Failed to generate salt: %v", err)
	}
	if len(salt) != int(DefaultParams.SaltLength) {
		t.Errorf("Generated salt length incorrect, expected %d, got %d", DefaultParams.SaltLength, len(salt))
	}
}
