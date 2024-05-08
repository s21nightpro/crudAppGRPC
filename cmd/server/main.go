package main

import (
	_ "github.com/lib/pq"
	_go "github.com/s21nightpro/crudAppGRPC/cmd/api"
	"github.com/s21nightpro/crudAppGRPC/internal/cache"
	"github.com/s21nightpro/crudAppGRPC/internal/db"
	"github.com/s21nightpro/crudAppGRPC/internal/server"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"log"
	"net"
)

func main() {
	logger, err := zap.NewProduction()
	if err != nil {
		log.Fatalf("failed to create zap logger: %v", err)
	}
	defer logger.Sync()
	cache := cache.NewCache()
	db, err := db.InitDB()
	if err != nil {
		logger.Fatal("failed to initialize database: %v", zap.Error(err))
	}
	defer db.Close()

	lis, err := net.Listen("tcp", ":50057")
	if err != nil {
		logger.Fatal("failed to listen: %v", zap.Error(err))
	}
	s := grpc.NewServer()
	defaultServer := server.NewServer(cache, db)

	_go.RegisterUserServiceServer(s, defaultServer)
	if err := s.Serve(lis); err != nil {
		logger.Fatal("failed to serve: %v", zap.Error(err))
	}
}
