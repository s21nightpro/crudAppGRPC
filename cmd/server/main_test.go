package main_test

import (
	"context"
	"database/sql"
	us "github.com/s21nightpro/crudAppGRPC/internal/grpc/user"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	_ "google.golang.org/grpc"
	"testing"
)

type mockUserServiceClient struct {
	mock.Mock
}

func (m *mockUserServiceClient) CreateUser(ctx context.Context, req *us.CreateUserRequest) (*us.User, error) {
	args := m.Called(ctx, req)
	return args.Get(0).(*us.User), args.Error(1)
}

func TestCreateUser(t *testing.T) {
	mockServiceClient := new(mockUserServiceClient)
	mockServiceClient.On("CreateUser", mock.Anything, mock.Anything).Return(&us.User{Id: "new-user-id", Name: "Test User", Email: "test@example.com"}, nil)

	server := &server{
		users: make(map[string]*us.User),
		cache: NewCache(),
		db:    &sql.DB{},
	}

	ctx := context.Background()
	req := &us.CreateUserRequest{Name: "Test User", Email: "test@example.com"}

	user, err := server.CreateUser(ctx, req)
	assert.NoError(t, err)
	assert.Equal(t, "new-user-id", user.Id)
	assert.Equal(t, "Test User", user.Name)
	assert.Equal(t, "test@example.com", user.Email)
}
