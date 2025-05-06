package port

import (
	"category/internal/core/domain"
	"context"
)

type CategoryService interface {
	GetCategory(ctx context.Context, params domain.ProductListParamsSt, ids []string) (domain.ProductCategoryListRep, error)
	GetCategoryByID(ctx context.Context, id string) (domain.ProductCategoryMain, error)
}

type CategoryRepository interface {
	GetCategory(ctx context.Context, params domain.ProductListParamsSt, ids []string) ([]domain.ProductCategoryMain, error)
	GetCategoryByID(ctx context.Context, id string) (domain.ProductCategoryMain, error)
}
