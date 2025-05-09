package main

import (
	http "category/internal/adapter/client/http"
	"category/internal/core/service"
	"log"
	"net/url"
)

const (
	apiURL = "api.com"
)

// This file is intended to launch project's services using http protocol client
func main() {
	URL, err := url.Parse(apiURL)
	if err != nil {
		log.Fatal(err)
	}

	categoryClient := http.NewCategoryClient(URL)
	cityClient := http.NewCityClient(URL)
	productClient := http.NewProductClient(URL)

	service.NewCategoryService(categoryClient)
	service.NewCityService(cityClient)
	service.NewProductService(productClient)
}
