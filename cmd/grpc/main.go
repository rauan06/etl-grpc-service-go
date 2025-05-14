package main

import (
	"context"
	"fmt"
	"log"

	"category/internal/adapter/client/grpc"
	"category/internal/core/domain"
)

const (
	priceURL   = "0.0.0.0:5052"
	productURL = "0.0.0.0:5050"
	stockURL   = "0.0.0.0:5051"
)

func main() {
	ctx := context.Background()

	client, err := grpc.NewStockClient(ctx, stockURL)
	if err != nil {
		log.Fatalf("Error initializing gRPC client: %v", err)
	}
	defer client.Close()

	categories, err := client.ListStocks(ctx, domain.ListParamsSt{}, []string{}, []string{})
	if err != nil {
		log.Fatalf("Error extracting: %v", err)
	}

	fmt.Println("Successfully extracted and processed:", categories)
}
