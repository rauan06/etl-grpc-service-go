package repository

import (
	"category/internal/core/domain"
	"context"
	"net/url"
)

type CategoryRepository struct {
	baseURL *url.URL
}

func NewCategoryRepository(baseURL *url.URL) *CategoryRepository {
	return &CategoryRepository{
		baseURL,
	}
}

func (r *CategoryRepository) GetCategory(ctx context.Context, params domain.ProductListParamsSt, ids []string) ([]domain.ProductCategoryMain, error) {
	panic("Imlement me")
}
func (r *CategoryRepository) GetCategoryByID(ctx context.Context, id uint64) (domain.ProductCategoryMain, error) {
	panic("Imlement me")
}
