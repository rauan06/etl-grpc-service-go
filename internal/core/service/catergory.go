package service

import (
	"category/internal/core/domain"
	"category/internal/core/port"
	"context"
)

type CategoryService struct {
	repo port.CategoryRepository
}

func NewCategoryService(repo port.CategoryRepository) *CategoryService {
	return &CategoryService{
		repo,
	}
}

func (s *CategoryService) ListCategories(ctx context.Context, params domain.ProductListParamsSt, ids []uint64) (domain.ProductCategoryListRep, error) {
	panic("Implement me")
}
func (s *CategoryService) GetCategory(ctx context.Context, id uint64) (domain.ProductCategoryMain, error) {
	panic("Implement me")
}
