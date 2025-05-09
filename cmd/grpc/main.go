package main

import (
	"category/internal/adapter/client/grpc"
	"category/internal/core/service"
	"context"
	"log"
)

const (
	apiURL = "api.com"
)

// This file is intended to launch project's services using grpc protocol client
func main() {
	categoryClient, err := grpc.NewCategoryClient(context.Background(), apiURL)
	if err != nil {
		log.Fatal(err)
	}

	cityClient, err := grpc.NewCityClient(context.Background(), apiURL)
	if err != nil {
		log.Fatal(err)
	}
	productClient, err := grpc.NewProductClient(context.Background(), apiURL)
	if err != nil {
		log.Fatal(err)
	}

	service.NewCategoryService(categoryClient)
	service.NewCityService(cityClient)
	service.NewProductService(productClient)
}
