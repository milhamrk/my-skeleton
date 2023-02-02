package main

import (
	"fazz/app/internal/applications/initiator"
	"fazz/app/internal/config"
	"fazz/app/pkg/logging"
	"fmt"

	"go.uber.org/zap"
)

func main() {
	logger, err := zap.NewProduction()
	if err != nil {
		fmt.Errorf("failed to create logger: %v", err)
	}
	defer logger.Sync() // nolint:errcheck
	appLogger := logging.NewLogger(logger, "fazz")
	logger.Info("init config")
	cfg := config.GetConfig()

	a, err := initiator.NewApp(cfg, appLogger)
	if err != nil {
		appLogger.Fatal("Error create app", zap.Error(err))
	}
	logger.Info("Running Application")
	a.Run()
}
