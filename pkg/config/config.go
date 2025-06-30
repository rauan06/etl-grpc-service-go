package config

import (
	"os"
)

// Container contains environment variables for the application
type (
	Container struct {
		Product *Product
		Price   *Price
		Stock   *Stock
	}
	// Environment variables for the application
	Product struct {
		URL      string
		PortHttp string
		PortGrpc string
	}
	Price struct {
		URL      string
		PortHttp string
		PortGrpc string
	}
	Stock struct {
		URL      string
		PortHttp string
		PortGrpc string
	}
)

// New creates a new container instance
func New() (*Container, error) {
	if os.Getenv("APP_ENV") != "production" {
		// err := godotenv.Load()
		// if err != nil {
		// 	return nil, err
		// }
	}

	product := &Product{
		URL:      os.Getenv("PRODUCT_URL"),
		PortHttp: os.Getenv("PRODUCT_PORT_HTTP"),
		PortGrpc: os.Getenv("PRODUCT_PORT_GRPC"),
	}

	price := &Price{
		URL:      os.Getenv("PRICE_URL"),
		PortHttp: os.Getenv("PRICE_PORT_HTTP"),
		PortGrpc: os.Getenv("PRICE_PORT_GRPC"),
	}

	stock := &Stock{
		URL:      os.Getenv("STOCK_URL"),
		PortHttp: os.Getenv("STOCK_PORT_HTTP"),
		PortGrpc: os.Getenv("STOCK_PORT_GRPC"),
	}

	return &Container{
		product,
		price,
		stock,
	}, nil
}
