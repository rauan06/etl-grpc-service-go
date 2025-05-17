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
	"github.com/rauan06/etl-grpc-service-go/internal/adapter/repository"
	"github.com/rauan06/etl-grpc-service-go/internal/core/port"
	"github.com/rauan06/etl-grpc-service-go/internal/core/service"
	"github.com/rauan06/etl-grpc-service-go/pkg/config"
	"github.com/rauan06/etl-grpc-service-go/pkg/lib/logger"
	pb "github.com/rauan06/etl-grpc-service-go/protos/etl/v1/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

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
	grpcCategoryClient, _ := clientGrpc.NewCategoryClient(ctx, cfg.Product.URL+":"+cfg.Product.PortGrpc)
	grpcProductClient, _ := clientGrpc.NewProductClient(ctx, cfg.Product.URL+":"+cfg.Product.PortGrpc)
	grpcCityClient, _ := clientGrpc.NewCityClient(ctx, cfg.Product.URL+":"+cfg.Product.PortGrpc)
	grpcPriceClient, _ := clientGrpc.NewPriceClient(ctx, cfg.Price.URL+":"+cfg.Price.PortGrpc)
	grpcStockClient, _ := clientGrpc.NewStockClient(ctx, cfg.Stock.URL+":"+cfg.Stock.PortGrpc)
	grpcClient := clientGrpc.NewClient(grpcCategoryClient, grpcCityClient, grpcProductClient, grpcPriceClient, grpcStockClient)

	// HTTP client
	// Parse URLs for HTTP
	productURL, err := url.Parse("http://" + cfg.Product.URL + ":" + cfg.Price.PortHttp)
	if err != nil {
		logger.Error("Error parsing url", "error", err, "url", cfg.Product.URL+":"+cfg.Price.PortHttp)
		os.Exit(1)
	}

	storeURL, err := url.Parse("http://" + cfg.Stock.URL + ":" + cfg.Product.PortHttp)
	if err != nil {
		logger.Error("Error parsing url", "error", err, "url", cfg.Stock.URL+":"+cfg.Product.PortHttp)
		os.Exit(1)
	}

	priceURL, err := url.Parse("http://" + cfg.Price.URL + ":" + cfg.Product.PortHttp)
	if err != nil {
		logger.Error("Error parsing url", "error", err, "url", cfg.Price.URL+":"+cfg.Product.PortHttp)
		os.Exit(1)
	}

	categoryHTTPClient := clientHttp.NewCategoryClient(productURL)
	cityHTTPClient := clientHttp.NewCityClient(productURL)
	priceHTTPClient := clientHttp.NewPriceClient(priceURL)
	stockHTTPClient := clientHttp.NewStockClient(storeURL)
	productHTTPClient := clientHttp.NewProductClient(productURL)
	HTTPClient := clientHttp.NewClient(categoryHTTPClient, cityHTTPClient, productHTTPClient, priceHTTPClient, stockHTTPClient)

	client := client.NewClient(grpcClient, HTTPClient, logger)

	// Repository layer
	// Storage
	repo := repository.NewRepositry()

	// Service Layer
	stockSvc := service.NewStockService(client, repo, logger)
	collectorSvc := service.NewCollectorService(client, repo, logger)

	svcs := []port.Service{&stockSvc, &collectorSvc}

	// Handler
	h := handler.NewEtlHandler(
		repo,
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
