package main

import (
	"context"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestDatabaseOperation(t *testing.T) {
	// Создание мока базы данных
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	// Настройка ожидаемого запроса и ответа
	rows := sqlmock.NewRows([]string{"id", "name", "email"}).
		AddRow(1, "John Doe", "john.doe@example.com")

	mock.ExpectQuery("^SELECT (.+) FROM users$").WillReturnRows(rows)

	// Выполнение запроса к моку базы данных
	var id int
	var name string
	var email string
	err = db.QueryRow(context.Background(), "SELECT id, name, email FROM users WHERE id = $1", 1).Scan(&id, &name, &email)
	assert.NoError(t, err)
	assert.Equal(t, 1, id)
	assert.Equal(t, "John Doe", name)
	assert.Equal(t, "john.doe@example.com", email)

	// Проверка, что все ожидания выполнены
	err = mock.ExpectationsWereMet()
	if err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
	clientRun()
}
