package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"sync"
	"time"

	user "github.com/s21nightpro/crudApp/crudApp/go/user"
	"google.golang.org/grpc"
)

type server struct {
	user.UnimplementedUserServiceServer
	users map[string]*user.User
	cache *Cache
	mu    sync.Mutex
}
type Cache struct {
	items map[string]Item
	mu    sync.Mutex
}

type Item struct {
	Value      interface{}
	Expiration int64
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
	s.cache.Set(req.Email, newUser, time.Minute*5)
	return newUser, nil
}

func (s *server) UpdateUser(ctx context.Context, req *user.UpdateUserRequest) (*user.User, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	// Проверяем, существует ли пользователь
	existingUser, exists := s.users[req.Email]
	if !exists {
		return nil, fmt.Errorf("user with email %s does not exist", req.Email)
	}

	// Обновляем информацию о пользователе
	existingUser.Name = req.Name
	existingUser.Email = req.Email
	s.cache.Set(req.Email, existingUser, time.Minute*5)

	return existingUser, nil
}

func (s *server) DeleteUser(ctx context.Context, req *user.DeleteUserRequest) (*user.User, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	userToDelete, exists := s.users[req.Id]
	if !exists {
		return nil, fmt.Errorf("user with ID %s does not exist", req.Id)
	}
	delete(s.users, req.Id)
	s.cache.Set(req.Id, userToDelete, time.Minute*5)

	return &user.User{Id: req.Id, Name: "Deleted User", Email: "deleted@example.com"}, nil
}

func (s *server) GetUser(ctx context.Context, req *user.GetUserRequest) (*user.User, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	// Проверяем кэш на наличие данных
	item, found := s.cache.Get(req.Id)
	if found {
		// Если данные найдены в кэше, возвращаем их
		// Добавляем утверждение типа для item
		userToGet, ok := item.(*user.User)
		if !ok {
			return nil, fmt.Errorf("failed to assert type of cached item")
		}
		return userToGet, nil
	}

	// Если данных нет в кэше, ищем в базе данных
	userToGet, exists := s.users[req.Id]
	if !exists {
		return nil, fmt.Errorf("user with ID %s does not exist", req.Id)
	}

	// Кэшируем данные
	s.cache.Set(req.Id, userToGet, time.Minute*5) // Кэшируем на 5 минут

	return userToGet, nil
}

func NewCache() *Cache {
	return &Cache{
		items: make(map[string]Item),
	}
}

func (c *Cache) Set(key string, value interface{}, duration time.Duration) {
	c.mu.Lock()
	defer c.mu.Unlock()
	expiration := time.Now().Add(duration).UnixNano()
	c.items[key] = Item{Value: value, Expiration: expiration}
}

func (c *Cache) Get(key string) (interface{}, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()
	item, found := c.items[key]
	if !found {
		return nil, false
	}
	if time.Now().UnixNano() > item.Expiration {
		delete(c.items, key)
		return nil, false
	}
	return item.Value, true
}

func (c *Cache) Delete(key string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	delete(c.items, key)
}

func main() {
	cache := NewCache() // Инициализация кэша

	lis, err := net.Listen("tcp", ":50057")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	user.RegisterUserServiceServer(s, &server{users: make(map[string]*user.User), cache: cache})
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
