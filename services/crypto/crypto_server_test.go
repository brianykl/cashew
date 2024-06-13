package main

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"testing"

	cryptopb "github.com/brianykl/cashew/services/crypto/pb"
	"github.com/stretchr/testify/assert"
)

func TestCreateHMAC(t *testing.T) {
	key := []byte("secret")
	data := "testdata"
	expectedHash := createHMAC(key, data)

	// Manually calculate the expected HMAC for comparison
	h := hmac.New(sha256.New, key)
	h.Write([]byte(data))
	expectedHashManual := hex.EncodeToString(h.Sum(nil))

	assert.Equal(t, expectedHashManual, expectedHash, "HMAC hashes should match")
}

func TestHashPII(t *testing.T) {
	server := &cryptoServer{}
	req := &cryptopb.HashPIIRequest{
		Data: "test@example.com",
		Key:  []byte("secret"),
	}
	resp, err := server.HashPII(context.Background(), req)
	assert.NoError(t, err)
	assert.NotEmpty(t, resp.EncodedHash, "Encoded hash should not be empty")
}

func TestVerifyPII(t *testing.T) {
	server := &cryptoServer{}
	reqHash := &cryptopb.HashPIIRequest{
		Data: "test@example.com",
		Key:  []byte("secret"),
	}
	hashResp, err := server.HashPII(context.Background(), reqHash)
	assert.NoError(t, err)

	reqVerify := &cryptopb.VerifyPIIRequest{
		Data:        "test@example.com",
		Key:         []byte("secret"),
		EncodedHash: hashResp.EncodedHash,
	}
	verifyResp, err := server.VerifyPII(context.Background(), reqVerify)
	assert.NoError(t, err)
	assert.True(t, verifyResp.IsValid, "Hashes should match and be valid")

}

func TestHashPassword(t *testing.T) {
	server := &cryptoServer{}
	req := &cryptopb.HashPasswordRequest{
		Password: "password123",
		Params: &cryptopb.Argon2IdParams{
			Memory:      64 * 1024,
			Iterations:  3,
			Parallelism: 2,
			SaltLength:  16,
			KeyLength:   32,
		},
	}
	resp, err := server.HashPassword(context.Background(), req)
	assert.NoError(t, err)
	assert.NotEmpty(t, resp.EncodedHash, "Encoded hash should not be empty")
}

func TestVerifyPassword(t *testing.T) {
	server := &cryptoServer{}
	reqHash := &cryptopb.HashPasswordRequest{
		Password: "password123",
		Params: &cryptopb.Argon2IdParams{
			Memory:      64 * 1024,
			Iterations:  3,
			Parallelism: 2,
			SaltLength:  16,
			KeyLength:   32,
		},
	}
	hashResp, err := server.HashPassword(context.Background(), reqHash)
	assert.NoError(t, err)

	reqVerify := &cryptopb.VerifyPasswordRequest{
		Password:    "password123",
		EncodedHash: hashResp.EncodedHash,
		Params: &cryptopb.Argon2IdParams{
			Memory:      64 * 1024,
			Iterations:  3,
			Parallelism: 2,
			SaltLength:  16,
			KeyLength:   32,
		},
	}
	verifyResp, err := server.VerifyPassword(context.Background(), reqVerify)
	assert.NoError(t, err)
	assert.True(t, verifyResp.IsValid, "Passwords should match and be valid")
}
