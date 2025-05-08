package main

import (
	"log"
	"log/slog"
	"net"
	"net/url"
	"os"

	handler "category/internal/adapter/handler/grpc"
	pb "category/internal/adapter/handler/grpc/product/v1"
	"category/internal/adapter/repository"
	"category/internal/core/service"
	"category/pkg/lib/logger"

	"google.golang.org/grpc"
)

const (
	apiURL = "apiLink.com"
)

func main() {
	// cfg := config.LoadConfig()

	logger := logger.SetupPrettySlog(os.Stdout)
	slog.SetDefault(logger)

	grpcServer := grpc.NewServer()

	URL, err := url.Parse(apiURL)
	if err != nil {
		log.Fatalf("Error parsing url: %v", err)
	}

	// #######################
	// Repository Layer
	// #######################
	categoryRepository := repository.NewCategoryRepository(URL)

	// #######################
	// Service Layer
	// #######################
	categoryService := service.NewCategoryService(categoryRepository)

	// #######################
	// Presentation Layer
	// #######################
	categoryHandler := handler.NewCategoryGrpcHandler(categoryService, logger)

	// #######################
	// Registering Endpoints
	// #######################
	pb.RegisterCategoryServer(grpcServer, categoryHandler)

	list, err := net.Listen("tcp", "localhost:9090")
	if err != nil {
		log.Fatal("couldn't start grpc server")
	}

	logger.Info("started grpc server on localhost:9090")

	if err := grpcServer.Serve(list); err != nil {
		log.Fatal("couldn't start grpc server")
	}
}
