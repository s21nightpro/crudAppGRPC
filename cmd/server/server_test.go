package main

import (
	"context"
	_go "github.com/s21nightpro/crudAppGRPC/internal/grpc/user"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

func TestCreateUser(t *testing.T) {
	// Создаем мок для DB
	dbMock, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer dbMock.Close()

	// Создаем мок для Cache
	cacheMock := new(Cache)
	cacheMock.On("Set", mock.Anything, mock.Anything, mock.Anything).Return(nil)

	// Инициализируем сервер с моками
	s := &server{
		cache: cacheMock,
		db:    dbMock,
	}

	// Создаем запрос
	req := &_go.CreateUserRequest{Name: "Test User", Email: "test@example.com"}

	// Выполняем функцию CreateUser
	user, err := s.CreateUser(context.Background(), req)

	// Проверяем, что нет ошибок
	assert.NoError(t, err)

	// Проверяем, что пользователь был создан
	assert.NotNil(t, user)

	// Проверяем, что запрос INSERT был выполнен
	mock.ExpectExec("INSERT INTO users (id, name, email) VALUES ($1, $2, $3)").WithArgs(mock.AnyArg(), req.Name, req.Email).WillReturnResult(sqlmock.NewResult(1, 1))

	// Проверяем, что данные были сохранены в кэше
	cacheMock.AssertExpectations(t)
}
