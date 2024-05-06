package logger

import (
	"fmt"
	"go.uber.org/zap"
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
