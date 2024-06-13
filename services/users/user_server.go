package main

import (
	"context"
	"encoding/hex"
	"log"
	"net"

	"github.com/brianykl/cashew/services/crypto/client"
	usermodels "github.com/brianykl/cashew/services/users/models"
	userpb "github.com/brianykl/cashew/services/users/pb"

	"github.com/google/uuid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/reflection"
	"gorm.io/gorm"
)

type userServer struct {
	userpb.UnimplementedUserServiceServer
	database     *gorm.DB
	cryptoClient *client.CryptoClient
	context      context.Context
}

// need to write better error handling for this and need to circle back around to implement encryption and hashing
func (s *userServer) CreateUser(ctx context.Context, req *userpb.CreateUserRequest) (*userpb.UserResponse, error) {
	log.Printf("we tried it lol")
	user, _ := usermodels.NewUser(req.Email, req.Name, req.Password)
	log.Printf("creating user: %s", req.GetName())

	// unsure if this is even necessary as a response, should this just return a bool?
	response := userpb.UserResponse{
		UserId:   uuid.New().String(),
		Email:    user.Email,    // encrypt email
		Name:     user.Name,     // encrypt name
		Password: user.Password, // hash password
	}

	if err := s.database.Create(user).Error; err != nil {
		log.Printf("failed to insert %v", err)
		return nil, err
	}

	log.Printf("inserted successfully")

	return &response, nil
}

// should this function be returning a bool instead of user?
func (s *userServer) VerifyUser(ctx context.Context, req *userpb.VerifyUserRequest) (*userpb.LoginResponse, error) {
	loginEmail := req.Email
	log.Printf("verifying user: %s", loginEmail)
	hexKey := "f13fd7ee2c6346b67aae8863ec68c170d26766a6fe216485ca5bfdfa1c25b233" // need to store somewhere safe
	key, _ := hex.DecodeString(hexKey)
	hashedLoginEmailResponse, err := s.cryptoClient.HashPII(s.context, loginEmail, key)
	if err != nil {
		log.Printf("failed to encrypt %v", err)
		return nil, err
	}

	hashedLoginEmail := hashedLoginEmailResponse.EncodedHash
	var user usermodels.User
	err = s.database.Where("email = ?", hashedLoginEmail).First(&user).Error
	if err != nil {
		log.Printf("failed to find user with this email %v", err)
		return nil, err
	}

	loginPassword := req.Password
	hashedPassword := user.Password

	params := client.Argon2IdParams{
		Memory:      64 * 1024, // 64 MiB of RAM
		Iterations:  3,         // 3 Iterations
		Parallelism: 2,         // Utilize 2 CPU cores (adjust if needed)
		SaltLength:  16,        // 16-byte salt
		KeyLength:   32,        // 32-byte output hash
	}

	passwordMatchResponse, err := s.cryptoClient.VerifyPassword(s.context, loginPassword, hashedPassword, &params)
	if err != nil {
		log.Printf("error comparing hashes %v", err)
		return nil, err
	}

	if !passwordMatchResponse.IsValid {
		log.Printf("incorrect password")
	}
	if passwordMatchResponse.IsValid {
		log.Printf("correct credentials. logging in...")
	}

	response := userpb.LoginResponse{
		IsValid: passwordMatchResponse.IsValid,
	}

	return &response, nil
}

// func (s *userServer) GetUser(ctx context.Context, req *userpb.GetUserRequest) (*userpb.UserResponse, error) {
// 	userID := req.UserId
// 	log.Printf("getting user: %s", userID)
// 	var user usermodels.User
// 	result := s.database.Where(&usermodels.User{UserID: userID}).First(&user)
// 	if result.Error != nil {
// 		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
// 			log.Printf("user with ID %s not found", userID)
// 			return nil, status.Errorf(codes.NotFound, "User with ID %s not found", userID)
// 		} else {
// 			log.Printf("error with ID %v not found", result.Error)
// 			return nil, status.Errorf(codes.Internal, "Error retrieving user: %v", result.Error)
// 		}
// 	}

// 	response := userpb.UserResponse{
// 		UserId:   user.UserID,
// 		Email:    user.Email,
// 		Name:     user.Name,
// 		Password: user.Password,
// 	}

// 	return &response, nil
// }

// func (s *userServer) UpdateUser(ctx context.Context, req *userpb.UpdateUserRequest) (*userpb.UserResponse, error) {

// 	userID := req.UserId
// 	log.Printf("updating user: %s", userID)
// 	var user usermodels.User
// 	result := s.database.Where(&usermodels.User{UserID: userID}).First(&user)
// 	if result.Error != nil {
// 		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
// 			log.Printf("user with ID %s not found", userID)
// 			return nil, status.Errorf(codes.NotFound, "User with ID %s not found", userID)
// 		} else {
// 			log.Printf("error with ID %v not found", result.Error)
// 			return nil, status.Errorf(codes.Internal, "Error retrieving user: %v", result.Error)
// 		}
// 	}

// 	// these need to be encrypted/hashed
// 	user.Email = req.GetEmail()
// 	user.Name = req.GetPassword()
// 	user.Password = req.GetPassword()

// 	result = s.database.Save(&user)
// 	if result.Error != nil {
// 		return nil, status.Errorf(codes.Internal, "error updating user: %v", result.Error)
// 	}

// 	// unsure if this needs to be the response
// 	response := userpb.UserResponse{
// 		UserId:   user.UserID,
// 		Email:    user.Email,
// 		Name:     user.Name,
// 		Password: user.Password,
// 	}

// 	return &response, nil
// }

func main() {
	lis, err := net.Listen("tcp", ":5001") // Use an appropriate port for your service
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	cryptoServiceAddress := "localhost:5002"
	conn, err := grpc.NewClient(cryptoServiceAddress, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("could not connect to crypto service: %v", err)
	}
	defer conn.Close()
	ctx := context.Background()
	cc := client.NewCryptoClient(conn)

	db, err := usermodels.Connect()
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}

	s := grpc.NewServer()
	userServer := &userServer{database: db, cryptoClient: cc, context: ctx}
	userpb.RegisterUserServiceServer(s, userServer)
	reflection.Register(s)

	log.Printf("user service server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
