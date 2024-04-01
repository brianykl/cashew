package client

import (
	"context"

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
