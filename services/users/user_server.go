package main

import (
	"context"
	"log"
	"net"

	usermodels "github.com/brianykl/cashew/services/users/models"
	userpb "github.com/brianykl/cashew/services/users/pb"

	"errors"

	"github.com/google/uuid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/status"
	"gorm.io/gorm"
)

type userServer struct {
	userpb.UnimplementedUserServiceServer
	db *gorm.DB
}

// need to write better error handling for this and need to circle back around to implement encryption and hashing
func (s *userServer) CreateUser(ctx context.Context, req *userpb.CreateUserRequest) (*userpb.UserResponse, error) {
	log.Printf("we tried it lol")
	user, _ := usermodels.NewUser(req.Email, req.Name, req.Password)
	log.Printf("creating user: %s", req.GetName())

	// unsure if this is even necessary as a response, should this just return a bool?
	response := userpb.UserResponse{
		UserId:   uuid.New().String(), // generate id
		Email:    user.Email,          // encrypt email
		Name:     user.Name,           // encrypt name
		Password: user.Password,       // hash password
	}

	if err := s.db.Create(user).Error; err != nil {
		log.Printf("failed to insert %v", err)
		return nil, err
	}

	log.Printf("inserted successfully")

	return &response, nil
}

func (s *userServer) GetUser(ctx context.Context, req *userpb.GetUserRequest) (*userpb.UserResponse, error) {
	userID := req.UserId
	log.Printf("getting user: %s", userID)
	var user usermodels.User
	result := s.db.Where(&usermodels.User{UserID: userID}).First(&user)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			log.Printf("user with ID %s not found", userID)
			return nil, status.Errorf(codes.NotFound, "User with ID %s not found", userID)
		} else {
			log.Printf("error with ID %v not found", result.Error)
			return nil, status.Errorf(codes.Internal, "Error retrieving user: %v", result.Error)
		}
	}

	response := userpb.UserResponse{
		UserId:   user.UserID,
		Email:    user.Email,
		Name:     user.Name,
		Password: user.Password,
	}

	return &response, nil
}

func (s *userServer) UpdateUser(ctx context.Context, req *userpb.UpdateUserRequest) (*userpb.UserResponse, error) {

	userID := req.UserId
	log.Printf("updating user: %s", userID)
	var user usermodels.User
	result := s.db.Where(&usermodels.User{UserID: userID}).First(&user)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			log.Printf("user with ID %s not found", userID)
			return nil, status.Errorf(codes.NotFound, "User with ID %s not found", userID)
		} else {
			log.Printf("error with ID %v not found", result.Error)
			return nil, status.Errorf(codes.Internal, "Error retrieving user: %v", result.Error)
		}
	}

	// these need to be encrypted/hashed
	user.Email = req.GetEmail()
	user.Name = req.GetPassword()
	user.Password = req.GetPassword()

	result = s.db.Save(&user)
	if result.Error != nil {
		return nil, status.Errorf(codes.Internal, "error updating user: %v", result.Error)
	}

	// unsure if this needs to be the response
	response := userpb.UserResponse{
		UserId:   user.UserID,
		Email:    user.Email,
		Name:     user.Name,
		Password: user.Password,
	}

	return &response, nil
}

func main() {
	lis, err := net.Listen("tcp", ":5001") // Use an appropriate port for your service
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	db, err := usermodels.Connect()
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}

	s := grpc.NewServer()
	userServer := &userServer{db: db}
	userpb.RegisterUserServiceServer(s, userServer)
	reflection.Register(s)

	log.Printf("user service server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
