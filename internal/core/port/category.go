package port

import (
	"category/internal/core/domain"
	"context"
)

type CategoryService interface {
	ListCategories(ctx context.Context, params domain.ProductListParamsSt, ids []string) (domain.ProductCategoryListRep, error)
	GetCategory(ctx context.Context, id uint64) (domain.ProductCategoryMain, error)
}

type CategoryRepository interface {
	GetCategory(ctx context.Context, params domain.ProductListParamsSt, ids []string) ([]domain.ProductCategoryMain, error)
	GetCategoryByID(ctx context.Context, id uint64) (domain.ProductCategoryMain, error)
}
