package main

import (
	"context"
	"fmt"
	"log"

	"category/internal/adapter/client/grpc"
	"category/internal/core/domain"
)

const (
	apiURL = "0.0.0.0:5434"
)

func main() {
	ctx := context.Background()

	client, err := grpc.NewPriceClient(ctx, apiURL)
	if err != nil {
		log.Fatalf("Error initializing gRPC client: %v", err)
	}
	defer client.Close()

	categories, err := client.ListPrices(ctx, domain.ListParamsSt{}, []string{}, []string{})
	if err != nil {
		log.Fatalf("Error extracting categories: %v", err)
	}

	fmt.Println("Successfully extracted and processed categories:", categories)
}
