package main

import (
	client "category/internal/adapter/client/grpc"
	clientHttp "category/internal/adapter/client/http"
	"category/internal/adapter/handler"
	redis "category/internal/adapter/repository"
	"category/internal/core/service"
	"category/pkg/config"
	"category/pkg/lib/logger"
	pb "category/protos/etl/v1/pb"
	"context"
	"log"
	"net"
	"net/http"
	"net/url"
	"os"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
)

func main() {
	ctx := context.Background()
	cfg := config.LoadConfig()
	logger := logger.SetupPrettySlog(os.Stdout)

	// Start gRPC server
	go func() {
		lis, err := net.Listen("tcp", ":5059")
		if err != nil {
			log.Fatalf("failed to listen: %v", err)
		}

		// Create gRPC server
		grpcServer := grpc.NewServer()

		// Create gRPC clients for each service
		grpcCategoryClient, err := client.NewCategoryClient(ctx, "0.0.0.0:5059")
		if err != nil {
			log.Fatalf("Failed to create category client: %v", err)
		}
		grpcCityClient, err := client.NewCityClient(ctx, "0.0.0.0:5059")
		if err != nil {
			log.Fatalf("Failed to create city client: %v", err)
		}
		grpcPriceClient, err := client.NewPriceClient(ctx, "0.0.0.0:5059")
		if err != nil {
			log.Fatalf("Failed to create price client: %v", err)
		}
		grpcStockClient, err := client.NewStockClient(ctx, "0.0.0.0:5059")
		if err != nil {
			log.Fatalf("Failed to create stock client: %v", err)
		}
		grpcProductClient, err := client.NewProductClient(ctx, "0.0.0.0:5059")
		if err != nil {
			log.Fatalf("Failed to create product client: %v", err)
		}

		URL, _ := url.Parse("http://0.0.0.0:8080")

		categoryHttpClient := clientHttp.NewCategoryClient(URL)
		cityHttpClient := clientHttp.NewCityClient(URL)
		priceHttpClient := clientHttp.NewPriceClient(URL)
		stockHttpClient := clientHttp.NewStockClient(URL)
		productHttpClient := clientHttp.NewProductClient(URL)

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

		// Create ETL handler
		h := handler.NewEtlHandler(categorySvc, citySvc, priceSvc, stockSvc, productSvc)

		// Register ETL handler to the gRPC server (only ONCE, and pass the handler itself)
		pb.RegisterCategoryServiceServer(grpcServer, h)

		log.Println("gRPC server listening on :5059")
		if err := grpcServer.Serve(lis); err != nil {
			log.Fatalf("failed to serve gRPC: %v", err)
		}
	}()

	// Start HTTP gateway
	mux := runtime.NewServeMux()
	opts := []grpc.DialOption{grpc.WithInsecure()}
	if err := pb.RegisterCategoryServiceHandlerFromEndpoint(ctx, mux, "localhost:5059", opts); err != nil {
		log.Fatalf("failed to start HTTP gateway: %v", err)
	}

	log.Println("HTTP Gateway listening on :8099")
	if err := http.ListenAndServe(":8099", mux); err != nil {
		log.Fatalf("failed to start HTTP server: %v", err)
	}
}
