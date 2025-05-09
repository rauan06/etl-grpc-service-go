package service

import (
	"category/internal/core/port"
)

type ProductService struct {
	client *port.ProductClient
}

func NewProductService(client *port.ProductClient) *ProductService {
	return &ProductService{
		client,
	}
}
