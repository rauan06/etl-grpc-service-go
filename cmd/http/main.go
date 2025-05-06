package main

import (
	"log"
	"log/slog"
	"net/url"
	"os"

	"category/internal/adapter/handler"
	"category/pkg/lib/logger"
)

const (
	apiURL = "apiLink.com"
)

func main() {
	// cfg := config.LoadConfig()

	logger := logger.SetupPrettySlog(os.Stdout)
	slog.SetDefault(logger)

	URL, err := url.Parse(apiURL)
	if err != nil {
		log.Fatalf("Error parsing url: %v", err)
	}

	httpSrv := handler.NewAPIServer(
		"0.0.0.0:8080",
		URL,
		logger,
	)
	httpSrv.Run()
}
