package handler

import "category/internal/core/port"

type ProductHandler struct {
	svc port.ProductService
}

func NewProducteHandler(svc port.ProductService) *ProductHandler {
	return &ProductHandler{
		svc,
	}
}
