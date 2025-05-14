//go:generate protoc ./protos/etl/*.proto --go_out=plugins=grpc:./pb
package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"
	"net/url"
	"os"

	client "category/internal/adapter/client/grpc"
	clientHttp "category/internal/adapter/client/http"
	"category/internal/adapter/handler"
	"category/internal/adapter/repository/redis"
	"category/internal/core/port"
	"category/internal/core/service"
	"category/pkg/config"
	"category/pkg/lib/logger"
	etlPb "category/protos/etl/v1/pb"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

const (
	dockerProductAddress = "demo-product"
	dockerPriceAddress   = "demo-price"
	dockerStoreAddress   = "demo-store"

	dockerHttpPort  = 8080
	productHttpPort = 8080
	storeHttpPort   = 8081
	priceHttpPort   = 8082

	dockerGrpcPort         = 5050
	defaultProductGrpcPort = 5050
	defaultStoreGrpcPort   = 5051
	defaultPriceGrpcPort   = 5052
)

var (
	productGrpcHost string
	priceGrpcHost   string
	storeGrpcHost   string

	productGrpcPort int
	priceGrpcPort   int
	storeGrpcPort   int
)

func isRunningInDocker() bool {
	if _, err := os.Stat("/.dockerenv"); err == nil {
		return true
	}
	return false
}

func main() {
	if isRunningInDocker() {
		log.Println("Running inside Docker")
		productGrpcHost = dockerProductAddress
		priceGrpcHost = dockerPriceAddress
		storeGrpcHost = dockerStoreAddress
	} else {
		log.Println("Running outside Docker")
		productGrpcHost = "0.0.0.0"
		priceGrpcHost = "0.0.0.0"
		storeGrpcHost = "0.0.0.0"
	}

	if isRunningInDocker() {
		productGrpcPort = dockerGrpcPort
		priceGrpcPort = dockerGrpcPort
		storeGrpcPort = dockerGrpcPort
	} else {
		productGrpcPort = defaultProductGrpcPort
		priceGrpcPort = defaultPriceGrpcPort
		storeGrpcPort = defaultStoreGrpcPort
	}

	ctx := context.Background()
	cfg := config.LoadConfig()
	logger := logger.SetupPrettySlog(os.Stdout)

	// Start gRPC server
	lis, err := net.Listen("tcp", ":5059")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	grpcServer := grpc.NewServer()

	// gRPC Clients
	grpcCategoryClient, _ := client.NewCategoryClient(ctx, fmt.Sprintf("%s:%d", priceGrpcHost, productGrpcPort))
	grpcCityClient, _ := client.NewCityClient(ctx, fmt.Sprintf("%s:%d", productGrpcHost, productGrpcPort))
	grpcPriceClient, _ := client.NewPriceClient(ctx, fmt.Sprintf("%s:%d", priceGrpcHost, priceGrpcPort))
	grpcStockClient, _ := client.NewStockClient(ctx, fmt.Sprintf("%s:%d", storeGrpcHost, storeGrpcPort))
	grpcProductClient, _ := client.NewProductClient(ctx, fmt.Sprintf("%s:%d", productGrpcHost, productGrpcPort))

	// HTTP Clients
	productURL, _ := url.Parse(fmt.Sprintf("http://0.0.0.0:%d", productHttpPort))
	storeURL, _ := url.Parse(fmt.Sprintf("http://0.0.0.0:%d", storeHttpPort))
	priceURL, _ := url.Parse(fmt.Sprintf("http://0.0.0.0:%d", priceHttpPort))

	categoryHttpClient := clientHttp.NewCategoryClient(productURL)
	cityHttpClient := clientHttp.NewCityClient(productURL)
	priceHttpClient := clientHttp.NewPriceClient(priceURL)
	stockHttpClient := clientHttp.NewStockClient(storeURL)
	productHttpClient := clientHttp.NewProductClient(productURL)

	// Cache (Redis)
	cache, err := redis.New(ctx, cfg)
	if err != nil {
		log.Fatalf("Failed to initialize Redis: %v", err)
	}

	// Service Layer
	categorySvc := service.NewCategoryService(grpcCategoryClient, categoryHttpClient, cache, logger)
	citySvc := service.NewCityService(grpcCityClient, cityHttpClient, cache, logger)
	priceSvc := service.NewPriceService(grpcPriceClient, priceHttpClient, cache, logger)
	stockSvc := service.NewStockService(grpcStockClient, stockHttpClient, cache, logger)
	productSvc := service.NewProductService(grpcProductClient, productHttpClient, cache, logger)
	collectorSvc := service.NewCollectorService(cache, logger)

	svcs := []port.Service{&categorySvc, &citySvc, &priceSvc, &stockSvc, &productSvc, &collectorSvc}

	// Handler
	h := handler.NewEtlHandler(
		cache,
		logger,
		svcs...,
	)

	// Register gRPC handler
	etlPb.RegisterETLServiceServer(grpcServer, h)

	// Register HTTP Gateway
	mux := runtime.NewServeMux()
	opts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}

	if err := etlPb.RegisterETLServiceHandlerFromEndpoint(ctx, mux, "localhost:5059", opts); err != nil {
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
