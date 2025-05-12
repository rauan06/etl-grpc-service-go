package main

import (
	"context"
	"log"
	"net/url"
	"os"

	"category/internal/adapter/client/grpc"
	"category/internal/adapter/client/http"
	"category/internal/core/service"
	"category/pkg/lib/logger"
)

func main() {
	ctx := context.Background()

	logger := logger.SetupPrettySlog(os.Stdout)

	grpcClient, _ := grpc.NewCategoryClient(ctx, "0.0.0.0:5050")

	URL, err := url.Parse("http://0.0.0.0:8080")
	if err != nil {
		log.Fatal(err.Error())
	}

	// Initialize http client
	httpClient := http.NewCategoryClient(URL)

	svc := service.NewCategoryService(grpcClient, httpClient)
	svc.Run(ctx, logger)
}
