package service

import (
	"category/internal/core/port"
)

type CategoryService struct {
	client *port.CategoryClient
}

func NewCategoryService(client *port.CategoryClient) *CategoryService {
	return &CategoryService{
		client,
	}
}
