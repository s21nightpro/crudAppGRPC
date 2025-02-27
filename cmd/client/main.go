package main

import (
	"context"
	"fmt"
	grpc2 "github.com/s21nightpro/crudAppGRPC/cmd/api"
	//"github.com/s21nightpro/crudAppGRPC/internal/logger"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"time"
)

var logger *zap.Logger

func Init() {
	var err error
	logger, err = zap.NewProduction()
	if err != nil {
		panic(fmt.Sprintf("failed to initialize zap logger: %v", err))
	}
	defer logger.Sync()
}

func Get() *zap.Logger {
	return logger
}

func createUser(ctx context.Context, c grpc2.UserServiceClient, name, email string) {
	r, err := c.CreateUser(ctx, &grpc2.CreateUserRequest{
		Name:  name,
		Email: email,
	})
	if err != nil {
		if status.Code(err) == codes.Unknown {
			logger.Error("Could not create user", zap.Error(err))
		} else {
			logger.Fatal("Could not create user", zap.Error(err))
		}
	} else {
		logger.Info("User created", zap.String("email", r.GetEmail()), zap.String("id", r.GetId()))
	}
}

func getUser(ctx context.Context, c grpc2.UserServiceClient, id string) {
	r, err := c.GetUser(ctx, &grpc2.GetUserRequest{Id: id})
	if err != nil {
		if status.Code(err) == codes.Unknown {
			logger.Error("Could not get user", zap.Error(err))
		} else {
			logger.Fatal("Could not get user", zap.Error(err))
		}
	} else if r != nil {
		logger.Info("User retrieved", zap.String("id", r.GetId()), zap.String("email", r.GetEmail()), zap.String("name", r.GetName()))
	} else {
		logger.Info("User retrieval response is nil")
	}
}

func updateUser(ctx context.Context, c grpc2.UserServiceClient, id, newEmail, newName string) (*grpc2.User, error) {
	req := &grpc2.UpdateUserRequest{
		Id:    id,
		Name:  newName,
		Email: newEmail,
	}
	updatedUser, err := c.UpdateUser(ctx, req)
	if err != nil {
		logger.Error("Error updating user", zap.Error(err))
		return nil, fmt.Errorf("could not update user: %v", err)
	}
	if updatedUser == nil {
		logger.Error("Update user response is nil")
		return nil, fmt.Errorf("update user response is nil")
	}
	logger.Info("User updated", zap.String("email", updatedUser.GetEmail()), zap.String("id", updatedUser.GetId()))
	//logger.Info("User updated", zap.String("email", updatedUser.GetEmail()))
	return updatedUser, nil
}

func deleteUser(ctx context.Context, c grpc2.UserServiceClient, id string) {
	a := &grpc2.User{Id: id}
	r, err := c.DeleteUser(ctx, &grpc2.DeleteUserRequest{Id: id})
	if err != nil {
		if status.Code(err) == codes.Unknown {
			logger.Error("Could not delete user", zap.Error(err))
		}
	} else if r != nil {
		logger.Info("User deleted", zap.String("id", a.Id))
	} else {
		logger.Info("User deletion response is nil")
	}
}

func main() {
	Init()
	conn, err := grpc.Dial("localhost:50057", grpc.WithInsecure(), grpc.WithBlock())
	defer conn.Close()
	if err != nil {
		logger.Fatal("did not connect", zap.Error(err))
	}

	c := grpc2.NewUserServiceClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	//createUser(ctx, c, "John Cock", "john.cock@example.com")
	//getUser(ctx, c, "67e44047-6eef-45fd-845a-ec53abc89b55")
	updateUser(ctx, c, "5339d08c-4825-4b9e-ae71-ac2f194b8a8b", "john.dicks@example.com", "John Dicks")
	//getUser(ctx, c, "1dbf99bf-f560-45b9-895b-e5c945ad6b46")
	//deleteUser(ctx, c, "67e44047-6eef-45fd-845a-ec53abc89b55")
	//getUser(ctx, c, "67e44047-6eef-45fd-845a-ec53abc89b55")
	//deleteUser(ctx, c, "1dbf99bf-f560-45b9-895b-e5c945ad6b46")
}
