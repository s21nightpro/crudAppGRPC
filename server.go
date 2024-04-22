package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"sync"

	user "github.com/s21nightpro/crudApp/crudApp/go/user"
	"google.golang.org/grpc"
)

type server struct {
	user.UnimplementedUserServiceServer
	users map[string]*user.User
	mu    sync.Mutex
}

func (s *server) CreateUser(ctx context.Context, req *user.CreateUserRequest) (*user.User, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	// Проверяем, существует ли пользователь с таким именем или адресом электронной почты
	if _, exists := s.users[req.Email]; exists {
		return nil, fmt.Errorf("user with email %s already exists", req.Email)
	}

	// Если пользователь не существует, создаем нового пользователя
	newUser := &user.User{Id: fmt.Sprintf("%d", len(s.users)+1), Name: req.Name, Email: req.Email}
	s.users[req.Email] = newUser

	return newUser, nil
}

func main() {
	lis, err := net.Listen("tcp", ":50057")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	user.RegisterUserServiceServer(s, &server{users: make(map[string]*user.User)})
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
