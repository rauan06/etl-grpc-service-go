//go:generate protoc ./protos/etl/*.proto --go_out=plugins=grpc:./pb
package main

import (
	client "category/internal/adapter/client/grpc"
	clientHttp "category/internal/adapter/client/http"
	"category/internal/adapter/handler"
	redis "category/internal/adapter/repository"
	"category/internal/core/service"
	"category/pkg/config"
	"category/pkg/lib/logger"
	etlPb "category/protos/etl/v1/pb"
	"context"
	"fmt"
	"log"
	"net"
	"net/http"
	"net/url"
	"os"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

const (
	productHttpPort = 8080
	priceHttpPort   = 8082
	storeHttpPort   = 8081

	productGrpcPort = 5050
	priceGrpcPort   = 5052
	storeGrpcPort   = 5051
)

func main() {
	ctx := context.Background()
	cfg := config.LoadConfig()
	logger := logger.SetupPrettySlog(os.Stdout)

	// Start gRPC server
	lis, err := net.Listen("tcp", ":5059")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	// Grpc clients
	grpcServer := grpc.NewServer()
	grpcCategoryClient, err := client.NewCategoryClient(ctx, fmt.Sprintf("0.0.0.0:%d", productGrpcPort))
	if err != nil {
		log.Fatalf("Failed to create category client: %v", err)
	}
	grpcCityClient, err := client.NewCityClient(ctx, fmt.Sprintf("0.0.0.0:%d", productGrpcPort))
	if err != nil {
		log.Fatalf("Failed to create city client: %v", err)
	}
	grpcPriceClient, err := client.NewPriceClient(ctx, fmt.Sprintf("0.0.0.0:%d", priceGrpcPort))
	if err != nil {
		log.Fatalf("Failed to create price client: %v", err)
	}
	grpcStockClient, err := client.NewStockClient(ctx, fmt.Sprintf("0.0.0.0:%d", storeGrpcPort))
	if err != nil {
		log.Fatalf("Failed to create stock client: %v", err)
	}
	grpcProductClient, err := client.NewProductClient(ctx, fmt.Sprintf("0.0.0.0:%d", productGrpcPort))
	if err != nil {
		log.Fatalf("Failed to create product client: %v", err)
	}

	// HTTP clients
	productURL, _ := url.Parse(fmt.Sprintf("http://0.0.0.0:%d", productHttpPort))
	storeURL, _ := url.Parse(fmt.Sprintf("http://0.0.0.0:%d", storeHttpPort))
	priceURL, _ := url.Parse(fmt.Sprintf("http://0.0.0.0:%d", priceHttpPort))

	categoryHttpClient := clientHttp.NewCategoryClient(productURL)
	cityHttpClient := clientHttp.NewCityClient(productURL)
	priceHttpClient := clientHttp.NewPriceClient(priceURL)
	stockHttpClient := clientHttp.NewStockClient(storeURL)
	productHttpClient := clientHttp.NewProductClient(productURL)

	cache, err := redis.New(ctx, cfg)
	if err != nil {
		log.Fatalf("Failed to initialize Redis: %v", err)
	}

	// Construct service layers from clients
	categorySvc := service.NewCategoryService(grpcCategoryClient, categoryHttpClient, cache, logger)
	citySvc := service.NewCityService(grpcCityClient, cityHttpClient, cache, logger)
	priceSvc := service.NewPriceService(grpcPriceClient, priceHttpClient, cache, logger)
	stockSvc := service.NewStockService(grpcStockClient, stockHttpClient, cache, logger)
	productSvc := service.NewProductService(grpcProductClient, productHttpClient, cache, logger)
	collectorSvc := service.NewCollectorService(cache, logger)

	// Create ETL handler
	h := handler.NewEtlHandler(categorySvc, citySvc, priceSvc, stockSvc, productSvc, collectorSvc, cache, logger)
	etlPb.RegisterETLServiceServer(grpcServer, h)

	// Register ETL handler to the gRPC server (only ONCE, and pass the handler itself)
	mux := runtime.NewServeMux()
	opts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}

	// Only ONE registration for the gateway
	if err := etlPb.RegisterETLServiceHandlerFromEndpoint(ctx, mux, "localhost:5059", opts); err != nil {
		log.Fatalf("failed to register gateway: %v", err)
	}

	go func() {
		log.Println("gRPC server listening on :5059")
		if err := grpcServer.Serve(lis); err != nil {
			log.Fatalf("failed to serve gRPC: %v", err)
		}
	}()

	log.Println("HTTP Gateway listening on :8099")
	if err := http.ListenAndServe(":8099", mux); err != nil {
		log.Fatalf("failed to start HTTP server: %v", err)
	}

}
