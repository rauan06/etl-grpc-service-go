package main

import (
	"context"
	"fmt"
	"log"

	"category/internal/adapter/client/grpc"
	"category/internal/core/domain"
	"category/internal/core/service"
)

const (
	apiURL = "0.0.0.0:5050"
)

func main() {
	ctx := context.Background()

	client, err := grpc.NewCategoryClient(ctx, apiURL)
	if err != nil {
		log.Fatalf("Error initializing gRPC client: %v", err)
	}
	defer client.Close()

	categoryService := service.NewCategoryService(client)

	categories, err := categoryService.ListCategories(ctx, domain.ListParamsSt{}, []string{})
	if err != nil {
		log.Fatalf("Error extracting categories: %v", err)
	}

	fmt.Println("Successfully extracted and processed categories:", categories)
}
