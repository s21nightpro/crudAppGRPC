package main

import (
	"context"
	"log"
	"time"

	user "github.com/s21nightpro/crudApp/crudApp/go/user" // Импортируйте ваш пакет с определениями сервиса
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func main() {
	// Создание подключения к серверу
	conn, err := grpc.Dial("localhost:50057", grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()

	// Создание клиента
	c := user.NewUserServiceClient(conn)

	// Создание контекста с таймаутом
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	// Отправка запроса на создание пользователя
	r, err := c.CreateUser(ctx, &user.CreateUserRequest{
		Name:  "John Doe",
		Email: "john.doe@example.com",
	})
	if err != nil {
		if status.Code(err) == codes.Unknown {
			log.Printf("Could not create user: %v", err)
			// Здесь можно добавить логику обработки ошибки, например, предложить пользователю ввести другой адрес электронной почты
		} else {
			log.Fatalf("Could not create user: %v", err)
		}
	} else {
		log.Printf("User created: %s", r.GetId())
	}
}
