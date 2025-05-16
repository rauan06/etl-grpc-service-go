//go:generate protoc ./protos/etl/*.proto --go_out=plugins=grpc:./pb
package main

import (
	"context"
	"log"
	"net"
	"net/http"
	"net/url"
	"os"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/rauan06/etl-grpc-service-go/internal/adapter/client"
	clientGrpc "github.com/rauan06/etl-grpc-service-go/internal/adapter/client/grpc"
	clientHttp "github.com/rauan06/etl-grpc-service-go/internal/adapter/client/http"
	"github.com/rauan06/etl-grpc-service-go/internal/adapter/handler"
	"github.com/rauan06/etl-grpc-service-go/internal/adapter/repository/redis"
	"github.com/rauan06/etl-grpc-service-go/internal/adapter/repository/storage"
	"github.com/rauan06/etl-grpc-service-go/internal/core/port"
	"github.com/rauan06/etl-grpc-service-go/internal/core/service"
	"github.com/rauan06/etl-grpc-service-go/pkg/config"
	"github.com/rauan06/etl-grpc-service-go/pkg/lib/logger"
	pb "github.com/rauan06/etl-grpc-service-go/protos/etl/v1/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func isRunningInDocker() bool {
	if _, err := os.Stat("/.dockerenv"); err == nil {
		return true
	}
	return false
}

func main() {
	// Setup logger
	logger := logger.SetupPrettySlog(os.Stdout)

	// Load env variables
	cfg, err := config.New()
	if err != nil {
		logger.Error("Error loading environment variables", "error", err)
		os.Exit(1)
	}

	// Start gRPC server
	lis, err := net.Listen("tcp", ":5059")
	if err != nil {
		logger.Error("Error loading environment variables", "error", err)
		os.Exit(1)
	}
	grpcServer := grpc.NewServer()

	// Client (presentation layer)
	ctx := context.Background()

	// gRPC Client
	grpcCategoryClient, _ := clientGrpc.NewCategoryClient(ctx, cfg.Product.Env)
	grpcProductClient, _ := clientGrpc.NewProductClient(ctx, cfg.Product.Env)
	grpcCityClient, _ := clientGrpc.NewCityClient(ctx, cfg.Product.Env)
	grpcPriceClient, _ := clientGrpc.NewPriceClient(ctx, cfg.Price.Env)
	grpcStockClient, _ := clientGrpc.NewStockClient(ctx, cfg.Stock.Env)
	grpcClient := clientGrpc.NewClient(grpcCategoryClient, grpcCityClient, grpcProductClient, grpcPriceClient, grpcStockClient)

	// HTTP Client
	productURL, _ := url.Parse(cfg.Product.Env)
	storeURL, _ := url.Parse(cfg.Stock.Env)
	priceURL, _ := url.Parse(cfg.Price.Env)

	categoryHttpClient := clientHttp.NewCategoryClient(productURL)
	cityHttpClient := clientHttp.NewCityClient(productURL)
	priceHttpClient := clientHttp.NewPriceClient(priceURL)
	stockHttpClient := clientHttp.NewStockClient(storeURL)
	productHttpClient := clientHttp.NewProductClient(productURL)
	httpClient := clientHttp.NewClient(categoryHttpClient, cityHttpClient, productHttpClient, priceHttpClient, stockHttpClient)

	client := client.NewClient(grpcClient, httpClient, logger)

	// Repository layer
	// Cache (Redis)
	cache, err := redis.New(ctx, cfg.Redis)
	if err != nil {
		log.Fatalf("Failed to initialize Redis: %v", err)
	}

	// Storage
	repo := storage.NewStorage()

	// Service Layer
	stockSvc := service.NewStockService(client, repo, logger)
	collectorSvc := service.NewCollectorService(client, repo, logger)

	svcs := []port.Service{&stockSvc, &collectorSvc}

	// Handler
	h := handler.NewEtlHandler(
		cache,
		logger,
		svcs...,
	)

	// Register gRPC handler
	pb.RegisterETLServiceServer(grpcServer, h)

	// Register HTTP Gateway
	mux := runtime.NewServeMux()
	opts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}

	if err := pb.RegisterETLServiceHandlerFromEndpoint(ctx, mux, "localhost:5059", opts); err != nil {
		log.Fatalf("failed to register gateway: %v", err)
	}

	// Serve gRPC
	go func() {
		log.Println("gRPC server listening on :5059")
		if err := grpcServer.Serve(lis); err != nil {
			log.Fatalf("failed to serve gRPC: %v", err)
		}
	}()

	// Serve HTTP
	log.Println("HTTP Gateway listening on :8099")
	if err := http.ListenAndServe(":8099", mux); err != nil {
		log.Fatalf("failed to start HTTP server: %v", err)
	}
}
