package port

import (
	"context"

	"category/internal/core/domain"
)

type CategoryClient interface {
	ListCategories(ctx context.Context, params domain.ListParamsSt, ids []string) (*domain.CategoryListRep, error)
	GetCategory(ctx context.Context, id string) (*domain.CategoryMain, error)
}
