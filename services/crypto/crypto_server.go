package main

import (
	"context"
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"log"
	"net"
	"strings"

	"golang.org/x/crypto/argon2"

	cryptopb "github.com/brianykl/cashew/services/crypto/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type cryptoServer struct {
	cryptopb.UnimplementedCryptoServiceServer
}

func createHMAC(key []byte, data string) string {
	h := hmac.New(sha256.New, key)
	h.Write([]byte(data))
	return hex.EncodeToString(h.Sum(nil))
}

func (s *cryptoServer) HashPII(ctx context.Context, req *cryptopb.HashPIIRequest) (*cryptopb.HashPIIResponse, error) {
	if len(req.Key) == 0 {
		return nil, status.Error(codes.InvalidArgument, "HashPIIRequest: missing key")
	}

	encodedHash := createHMAC(req.Key, req.Data)
	response := cryptopb.HashPIIResponse{
		EncodedHash: encodedHash,
	}
	return &response, nil
}

func (s *cryptoServer) VerifyPII(ctx context.Context, req *cryptopb.VerifyPIIRequest) (*cryptopb.VerifyPIIResponse, error) {
	if len(req.Key) == 0 {
		return nil, status.Error(codes.InvalidArgument, "VerifyPIIRequest: missing key")
	}

	expectedEncodedHash := createHMAC(req.Key, req.Data)
	isValid := hmac.Equal([]byte(expectedEncodedHash), []byte(req.EncodedHash))
	response := cryptopb.VerifyPIIResponse{
		IsValid: isValid,
	}
	return &response, nil
}

func (s *cryptoServer) HashPassword(ctx context.Context, req *cryptopb.HashPasswordRequest) (*cryptopb.HashPasswordResponse, error) {
	password := req.Password
	params := req.Params
	salt, err := generateSalt(params.SaltLength)
	if err != nil {
		return nil, err
	}

	if params.Parallelism > 0xFF {
		log.Fatalf("Parallelism parameter out of range: %d", params.Parallelism)
	}

	hash := argon2.IDKey([]byte(password), salt, params.Iterations, params.Memory, uint8(params.Parallelism), params.KeyLength)
	encodedHash := base64.RawStdEncoding.EncodeToString(salt) + "$" + base64.RawStdEncoding.EncodeToString(hash)

	result := cryptopb.HashPasswordResponse{
		EncodedHash: encodedHash,
	}

	return &result, nil
}

func (s *cryptoServer) VerifyPassword(ctx context.Context, req *cryptopb.VerifyPasswordRequest) (*cryptopb.VerifyPasswordResponse, error) {
	password := req.Password
	encodedHash := req.EncodedHash
	params := req.Params

	if params.Parallelism > 0xFF {
		log.Fatalf("Parallelism parameter out of range: %d", params.Parallelism)
	}

	parts := strings.Split(encodedHash, "$")
	result := cryptopb.VerifyPasswordResponse{
		IsValid: false,
	}
	if len(parts) != 2 {
		return &result, nil // incorrect hash format
	}
	salt, err := base64.RawStdEncoding.DecodeString(parts[0])
	if err != nil {
		return &result, err
	}
	originalHash, err := base64.RawStdEncoding.DecodeString(parts[1])
	if err != nil {
		return &result, err
	}

	hash := argon2.IDKey([]byte(password), salt, params.Iterations, params.Memory, uint8(params.Parallelism), params.KeyLength)
	result.IsValid = compareHashes(originalHash, hash)
	return &result, nil
}

func generateSalt(length uint32) ([]byte, error) {
	salt := make([]byte, length)
	_, err := rand.Read(salt)
	if err != nil {
		return nil, err
	}
	return salt, nil
}

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

func main() {
	lis, err := net.Listen("tcp", ":5002")
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
