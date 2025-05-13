package main

import (
	"context"
	"log"
	"net/url"
	"os"

	"category/internal/adapter/client/grpc"
	"category/internal/adapter/client/http"
	redis "category/internal/adapter/repository"
	"category/internal/core/service"
	"category/pkg/config"
	"category/pkg/lib/logger"
)

func main() {
	cfg := config.GetConfing()
	ctx := context.Background()

	logger := logger.SetupPrettySlog(os.Stdout)

	grpcClient, err := grpc.NewCategoryClient(ctx, "0.0.0.0:5050")
	if err != nil {
		log.Fatal(err)
	}

	URL, err := url.Parse("http://0.0.0.0:8080")
	if err != nil {
		log.Fatal(err)
	}

	// Initialize http client
	httpClient := http.NewCategoryClient(URL)

	cache, err := redis.New(ctx, cfg)
	if err != nil {
		log.Fatal(err)
	}

	svc := service.NewCategoryService(grpcClient, httpClient, cache, logger)
	svc.Stop()
}
