package main

import (
	"context"
	"log"
	"net"

	userpb "github.com/brianykl/cashew/services/user/pb"
	"google.golang.org/grpc"
)

type userServer struct {
	userpb.UnimplementedUserServiceServer
}

func (s *userServer) CreateUser(ctx context.Context, req *userpb.CreateUserRequest) (*userpb.UserResponse, error) {
	// create user logic
	log.Printf("creating user: %s", req.GetName())
	response := userpb.UserResponse{
		UserId:   "dinky",    // generate id
		Email:    "pinky",    // encrypt email
		Name:     "donkey",   // encrypt name
		Password: "shlonkey", // hash password
	}

	// store response in db

	return &response, nil
}

func (s *userServer) GetUser(ctx context.Context, req *userpb.GetUserRequest) (*userpb.UserResponse, error) {
	// get user logic
	log.Printf("getting user: %s", req.GetUserId())
	response := userpb.UserResponse{
		UserId:   "bumbo",
		Email:    "wumbo",
		Name:     "mumbo",
		Password: "gumbo",
	}

	return &response, nil
}

func (s *userServer) UpdateUser(ctx context.Context, req *userpb.UpdateUserRequest) (*userpb.UserResponse, error) {
	// update user logic

	log.Printf("updating user: %s", req.GetUserId())
	response := userpb.UserResponse{
		UserId:   "bumbo",
		Email:    "wumbo",
		Name:     "mumbo",
		Password: "jumbo",
	}

	return &response, nil
}

func main() {
	lis, err := net.Listen("tcp", ":50052") // Use an appropriate port for your service
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()

	userpb.RegisterUserServiceServer(s, &userServer{})

	log.Printf("User service server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
