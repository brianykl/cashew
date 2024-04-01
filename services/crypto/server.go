package main

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"log"
	"net"

	"golang.org/x/crypto/argon2"

	cryptopb "github.com/brianykl/cashew/services/crypto/pb"
	"google.golang.org/grpc"
)

type cryptoServer struct {
	cryptopb.UnimplementedCryptoServiceServer
}

func (s *cryptoServer) HashPassword(ctx context.Context, req *cryptopb.HashPasswordRequest) (*cryptopb.HashPasswordResponse, error) {
	password := req.Password
	params := req.Params
	salt, err := GenerateSalt(params.SaltLength)
	if err != nil {
		return nil, err
	}

	if params.Parallelism > 0xFF {
		log.Fatalf("Parallelism parameter out of range: %d", params.Parallelism)
	}

	hash := argon2.IDKey([]byte(password), salt, params.Iterations, params.Memory, uint8(params.Parallelism), params.KeyLength)
	encodedHash := base64.RawStdEncoding.EncodeToString(hash) + "$" + base64.RawStdEncoding.EncodeToString(salt)

	result := cryptopb.HashPasswordResponse{
		EncodedHash: encodedHash,
	}

	return &result, nil
}

func GenerateSalt(length uint32) ([]byte, error) {
	salt := make([]byte, length)
	_, err := rand.Read(salt)
	if err != nil {
		return nil, err
	}
	return salt, nil
}

// func compareHashes(a, b []byte) bool {
// 	if len(a) != len(b) {
// 		return false
// 	}
// 	var result byte
// 	for i := 0; i < len(a); i++ {
// 		result |= a[i] ^ b[i]
// 	}
// 	return result == 0
// }

func main() {
	lis, err := net.Listen("tcp", ":50052")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()
	cryptoServer := &cryptoServer{}
	cryptopb.RegisterCryptoServiceServer(s, cryptoServer)

	log.Printf("crypto service server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
