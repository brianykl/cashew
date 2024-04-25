package client

import (
	"context"
	"log"

	cryptopb "github.com/brianykl/cashew/services/crypto/pb"
	"google.golang.org/grpc"
)

type CryptoClient struct {
	service cryptopb.CryptoServiceClient
	// conn    *grpc.ClientConn
}

type Argon2IdParams struct {
	Memory      uint32
	Iterations  uint32
	Parallelism uint8
	SaltLength  uint32
	KeyLength   uint32
}

func NewCryptoClient(cc *grpc.ClientConn) *CryptoClient {
	service := cryptopb.NewCryptoServiceClient(cc)
	return &CryptoClient{service: service}
}

func (c *CryptoClient) Encrypt(ctx context.Context, plaintext string, key []byte) (*cryptopb.EncryptResponse, error) {
	req := &cryptopb.EncryptRequest{
		Plaintext: plaintext,
		Key:       key,
	}
	log.Printf("we did it!")
	return c.service.Encrypt(ctx, req)
}

func (c *CryptoClient) Decrypt(ctx context.Context, ciphertext string, key []byte) (*cryptopb.DecryptResponse, error) {
	req := &cryptopb.DecryptRequest{
		Ciphertext: ciphertext,
		Key:        key,
	}

	return c.service.Decrypt(ctx, req)
}

func (c *CryptoClient) HashPassword(ctx context.Context, naked_pw string, params *Argon2IdParams) (*cryptopb.HashPasswordResponse, error) {
	pbParams := &cryptopb.Argon2IdParams{
		Memory:      params.Memory,
		Iterations:  params.Iterations,
		Parallelism: uint32(params.Parallelism),
		SaltLength:  params.SaltLength,
		KeyLength:   params.KeyLength,
	}

	req := &cryptopb.HashPasswordRequest{
		Password: naked_pw,
		Params:   pbParams,
	}

	return c.service.HashPassword(ctx, req)
}

func (c *CryptoClient) VerifyPassword(ctx context.Context, password, encodedHash string, params *Argon2IdParams) (*cryptopb.VerifyPasswordResponse, error) {
	pbParams := &cryptopb.Argon2IdParams{
		Memory:      params.Memory,
		Iterations:  params.Iterations,
		Parallelism: uint32(params.Parallelism),
		SaltLength:  params.SaltLength,
		KeyLength:   params.KeyLength,
	}

	req := &cryptopb.VerifyPasswordRequest{
		Password:    password,
		EncodedHash: encodedHash,
		Params:      pbParams,
	}

	return c.service.VerifyPassword(ctx, req)
}
