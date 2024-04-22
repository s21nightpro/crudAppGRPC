package main

import (
	"context"
	"log"
	"net"

	"github.com/s21nightpro/crudApp"
	"google.golang.org/grpc"
)

type server struct {
	user.UnimplementedUserServiceServer
}

func (s *server) CreateUser(ctx context.Context, req *user.CreateUserRequest) (*user.User, error) {
	// Реализация создания пользователя
	return &user.User{Id: "1", Name: req.Name, Email: req.Email}, nil
}

func (s *server) GetUser(ctx context.Context, req *user.GetUserRequest) (*user.User, error) {
	// Реализация получения пользователя
	return &user.User{Id: "1", Name: "John Doe", Email: "john.doe@example.com"}, nil
}

func (s *server) UpdateUser(ctx context.Context, req *user.UpdateUserRequest) (*user.User, error) {
	// Реализация обновления пользователя
	return &user.User{Id: req.Id, Name: req.Name, Email: req.Email}, nil
}

func (s *server) DeleteUser(ctx context.Context, req *user.DeleteUserRequest) (*user.User, error) {
	// Реализация удаления пользователя
	return &user.User{Id: req.Id}, nil
}

func main() {
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	user.RegisterUserServiceServer(s, &server{})
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
