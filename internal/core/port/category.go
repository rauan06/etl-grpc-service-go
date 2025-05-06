package port

import (
	"category/internal/core/domain"
	"context"
)

type CategoryService interface {
	ListCategories(ctx context.Context, params domain.ProductListParamsSt, ids []string) (domain.ProductCategoryListRep, error)
	GetCategory(ctx context.Context, id int64) (domain.ProductCategoryMain, error)
}

type CategoryRepository interface {
	ListCategories(ctx context.Context, params domain.ProductListParamsSt, ids []int64) ([]domain.ProductCategoryMain, error)
	GetCategory(ctx context.Context, id int64) (domain.ProductCategoryMain, error)
}
