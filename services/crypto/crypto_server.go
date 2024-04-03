package main

import (
	"context"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"io"
	"log"
	"net"
	"strings"

	"golang.org/x/crypto/argon2"

	cryptopb "github.com/brianykl/cashew/services/crypto/pb"
	"google.golang.org/grpc"
)

type cryptoServer struct {
	cryptopb.UnimplementedCryptoServiceServer
}

func (s *cryptoServer) Encrypt(ctx context.Context, req *cryptopb.EncryptRequest) (*cryptopb.EncryptResponse, error) {
	plaintext := req.Plaintext
	key := req.Key
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, err
	}

	ciphertext := gcm.Seal(nonce, nonce, []byte(plaintext), nil)

	response := cryptopb.EncryptResponse{
		Ciphertext: base64.StdEncoding.EncodeToString(ciphertext),
	}
	return &response, nil
}

func (s *cryptoServer) Decrypt(ctx context.Context, req *cryptopb.DecryptRequest) (*cryptopb.DecryptResponse, error) {
	ciphertext := req.Ciphertext
	key := req.Key

	data, err := base64.StdEncoding.DecodeString(ciphertext)
	if err != nil {
		return nil, err
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	if len(data) < gcm.NonceSize() {
		return nil, err
	}

	nonce, ciphertextData := data[:gcm.NonceSize()], data[gcm.NonceSize():]
	plaintext, err := gcm.Open(nil, nonce, ciphertextData, nil)
	if err != nil {
		return nil, err
	}

	response := cryptopb.DecryptResponse{
		Plaintext: string(plaintext),
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
	encodedHash := base64.RawStdEncoding.EncodeToString(hash) + "$" + base64.RawStdEncoding.EncodeToString(salt)

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
