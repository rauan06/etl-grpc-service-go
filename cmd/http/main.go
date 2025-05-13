package main

import (
	"context"
	"fmt"
	"log"
	"net/url"

	"category/internal/adapter/client/http"
	"category/internal/core/domain"
)

const (
	apiURL = "http://0.0.0.0:8082"
)

func main() {
	ctx := context.Background()

	URL, err := url.Parse(apiURL)
	if err != nil {
		log.Fatal(err.Error())
	}

	// Initialize http client
	client := http.NewPriceClient(URL)

	// Initialize service with the client

	// Extract data
	// params := domain.ListParamsSt{
	// 	Page:     "1",
	// 	PageSize: "10",
	// 	Sort:     []string{"name"},
	// }
	// ids := []string{"123", "456"}

	categories, err := client.ListPrices(ctx, domain.ListParamsSt{}, []string{}, []string{})
	if err != nil {
		log.Fatalf("Error extracting categories: %v", err)
	}

	// Load data into PostgreSQL (you need to implement the Postgres insertion logic here)

	fmt.Println("Successfully extracted and processed categories:", categories)
}
