package client

import (
	"context"
	"log"

	userpb "github.com/brianykl/cashew/services/user/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type UserClient struct {
	service userpb.UserServiceClient
	conn    *grpc.ClientConn
}

func NewUserClientWithConn(cc *grpc.ClientConn) *UserClient {
	service := userpb.NewUserServiceClient(cc)
	return &UserClient{service: service}
}

func NewUserClient(serviceAddress string) (*UserClient, error) {
	conn, err := grpc.Dial(serviceAddress, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("could not connect to user service: %v", err)
	}

	service := userpb.NewUserServiceClient(conn)
	return &UserClient{service: service, conn: conn}, nil
}

func (c *UserClient) GetUser(ctx context.Context, id string) (*userpb.UserResponse, error) {
	req := &userpb.GetUserRequest{
		UserId: id,
	}
	return c.service.GetUser(ctx, req)
}

func (c *UserClient) CreateUser(ctx context.Context, name, email, password string) (*userpb.UserResponse, error) {
	req := &userpb.CreateUserRequest{
		Email:    email,
		Name:     name,
		Password: password,
	}
	return c.service.CreateUser(ctx, req)
}

func (c *UserClient) UpdateUser(ctx context.Context, id, name, email, password string) (*userpb.UserResponse, error) {
	req := &userpb.UpdateUserRequest{
		UserId:   id,
		Email:    email,
		Name:     name,
		Password: password,
	}
	return c.service.UpdateUser(ctx, req)
}
