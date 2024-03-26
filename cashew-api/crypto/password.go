package crypto

import (
	"crypto/rand"
	"encoding/base64"
	"strings"

	"golang.org/x/crypto/argon2"
)

// Argon2idParams holds the parameters for the Argon2id hashing algorithm.
type Argon2idParams struct {
	Memory      uint32
	Iterations  uint32
	Parallelism uint8
	SaltLength  uint32
	KeyLength   uint32
}

// DefaultParams provides a reasonable default for the Argon2id parameters.
var DefaultParams = Argon2idParams{
	Memory:      64 * 1024, // 64 MB
	Iterations:  3,
	Parallelism: 2,
	SaltLength:  16,
	KeyLength:   32,
}

// GenerateSalt generates a random salt.
func GenerateSalt(length uint32) ([]byte, error) {
	salt := make([]byte, length)
	_, err := rand.Read(salt)
	// Note: Using crypto/rand for secure random salt generation
	if err != nil {
		return nil, err
	}
	return salt, nil
}

// HashPassword hashes a password using Argon2id.
func HashPassword(password string, p *Argon2idParams) (string, error) {
	// Generate a random salt
	salt, err := GenerateSalt(p.SaltLength)
	if err != nil {
		return "", err
	}

	// Hash the password with Argon2id
	hash := argon2.IDKey([]byte(password), salt, p.Iterations, p.Memory, p.Parallelism, p.KeyLength)

	// Return the encoded hash (including the salt and parameters for verification)
	encodedHash := base64.RawStdEncoding.EncodeToString(hash)
	encodedSalt := base64.RawStdEncoding.EncodeToString(salt)
	return encodedSalt + "$" + encodedHash, nil
}

// VerifyPassword compares a plaintext password against the stored hash.
func VerifyPassword(password, encodedHash string, p *Argon2idParams) (bool, error) {
	// Decode the hash and salt from the encodedHash
	parts := strings.Split(encodedHash, "$")
	if len(parts) != 2 {
		return false, nil // Incorrect hash format
	}
	salt, err := base64.RawStdEncoding.DecodeString(parts[0])
	if err != nil {
		return false, err
	}
	originalHash, err := base64.RawStdEncoding.DecodeString(parts[1])
	if err != nil {
		return false, err
	}

	// Hash the input password using the same salt and parameters
	hash := argon2.IDKey([]byte(password), salt, p.Iterations, p.Memory, p.Parallelism, p.KeyLength)

	// Compare the newly generated hash with the original hash
	return compareHashes(originalHash, hash), nil
}

// compareHashes safely compares two hashes to prevent timing attacks.
func compareHashes(a, b []byte) bool {
	if len(a) != len(b) {
		return false
	}
	var result byte
	for i := 0; i < len(a); i++ {
		result |= a[i] ^ b[i]
	}
	return result == 0
}
