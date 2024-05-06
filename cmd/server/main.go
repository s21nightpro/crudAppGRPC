package main

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/google/uuid"

	_ "github.com/lib/pq"
	_ "github.com/s21nightpro/crudAppGRPC/internal/cache"
	"github.com/s21nightpro/crudAppGRPC/internal/grpc/user"
	"go.uber.org/zap"
	_ "go.uber.org/zap"
	"google.golang.org/grpc"
	"log"
	"net"
	"sync"
	"time"
)

type server struct {
	_go.UnimplementedUserServiceServer
	users map[string]*_go.User
	cache *Cache
	db    *sql.DB
	mu    sync.Mutex
}

func initDB() (*sql.DB, error) {
	connStr := "user=biba dbname=postgres password=boba host=localhost sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}
	if err = db.Ping(); err != nil {
		return nil, err
	}
	return db, nil
}

func (s *server) userExists(id string) (bool, error) {
	var exists bool
	err := s.db.QueryRow("SELECT EXISTS(SELECT 1 FROM users WHERE id = $1)", id).Scan(&exists)
	if err != nil {
		return false, err
	}
	return exists, nil
}

func (s *server) CreateUser(ctx context.Context, req *_go.CreateUserRequest) (*_go.User, error) {
	// Генерация UUID для нового пользователя
	userID := uuid.New().String()

	// Вставляем нового пользователя в базу данных с сгенерированным UUID
	_, err := s.db.Exec("INSERT INTO users (id, name, email) VALUES ($1, $2, $3)", userID, req.Name, req.Email)
	if err != nil {
		return nil, fmt.Errorf("error inserting user into database: %v", err)
	}

	// Возвращаем созданного пользователя с сгенерированным UUID
	createdUser := &_go.User{Id: userID, Name: req.Name, Email: req.Email}
	return createdUser, nil
}

func (s *server) UpdateUser(ctx context.Context, req *_go.UpdateUserRequest) (*_go.User, error) {
	exists, err := s.userExists(req.Id)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, fmt.Errorf("user with ID %s does not exist", req.Id)
	}

	// Продолжайте с обновлением пользователя, если ID существует
	_, err = s.db.Exec("UPDATE users SET name = $1, email = $2 WHERE id = $3", req.Name, req.Email, req.Id)
	if err != nil {
		return nil, err
	}

	updatedUser := &_go.User{Id: req.Id, Name: req.Name, Email: req.Email}
	return updatedUser, nil
}

func (s *server) DeleteUser(ctx context.Context, req *_go.DeleteUserRequest) (*_go.User, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	result, err := s.db.Exec("DELETE FROM users WHERE id = $1", req.Id)
	if err != nil {
		return nil, err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return nil, err
	}
	if rowsAffected == 0 {
		return nil, fmt.Errorf("user with ID %s does not exist", req.Id)
	}

	s.cache.Delete(req.Id)

	return nil, nil
}

func (s *server) GetUser(ctx context.Context, req *_go.GetUserRequest) (*_go.User, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	item, found := s.cache.Get(req.Id)
	if found {
		userToGet, ok := item.(*_go.User)
		if !ok {
			return nil, fmt.Errorf("failed to assert type of cached item")
		}
		return userToGet, nil
	}

	var userToGet _go.User
	err := s.db.QueryRowContext(ctx, "SELECT id, name, email FROM users WHERE id = $1", req.Id).Scan(&userToGet.Id, &userToGet.Name, &userToGet.Email)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("user with ID %s does not exist", req.Id)
		}
		return nil, fmt.Errorf("error querying user from database: %v", err)
	}

	s.cache.Set(req.Id, &userToGet, time.Minute*5)

	return &userToGet, nil
}

var logger *zap.Logger

func init() {
	var err error
	logger, err = zap.NewProduction()
	if err != nil {
		panic(fmt.Sprintf("failed to initialize zap logger: %v", err))
	}
	defer logger.Sync()
}

func main() {
	logger, err := zap.NewProduction()
	if err != nil {
		log.Fatalf("failed to create zap logger: %v", err)
	}
	defer logger.Sync()
	cache := NewCache()
	db, err := initDB()
	if err != nil {
		logger.Fatal("failed to initialize database: %v", zap.Error(err))
	}
	defer db.Close()
	lis, err := net.Listen("tcp", ":50057")
	if err != nil {
		logger.Fatal("failed to listen: %v", zap.Error(err))
	}
	s := grpc.NewServer()
	_go.RegisterUserServiceServer(s, &server{users: make(map[string]*_go.User), cache: cache, db: db})
	if err := s.Serve(lis); err != nil {
		logger.Fatal("failed to serve: %v", zap.Error(err))
	}
}
