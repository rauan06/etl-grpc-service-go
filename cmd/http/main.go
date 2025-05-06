package main

import (
	"category/internal/adapter/handler"
	"category/pkg/lib/logger"
	"log/slog"
	"os"
)

func main() {
	// cfg := config.LoadConfig()

	logger := logger.SetupPrettySlog(os.Stdout)
	slog.SetDefault(logger)

	httpSrv := handler.NewAPIServer(
		"0.0.0.0:8080",
		"apiLink.com",
		logger,
	)
	httpSrv.Run()
}
