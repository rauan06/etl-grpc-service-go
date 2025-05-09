package service

import (
	"category/internal/core/domain"
	"category/internal/core/port"
	"context"
)

type CategoryService struct {
	client *port.CategoryClient
}

func NewCategoryService(client *port.CategoryClient) *CategoryService {
	return &CategoryService{
		client,
	}
}

func (s *CategoryService) ListCategories(ctx context.Context, params domain.ListParamsSt, ids []int64) (*domain.CategoryListRep, error) {
	return nil, nil
}

func (s *CategoryService) GetCategory(ctx context.Context, id int64) (*domain.CategoryMain, error) {
	return nil, nil
}
