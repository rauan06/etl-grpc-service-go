package port

import (
	"context"
	"net/url"

	"category/internal/core/domain"
)

type CategoryService interface {
	ListCategories(ctx context.Context, params url.Values) (*domain.ProductCategoryListRep, error)
	GetCategory(ctx context.Context, id int64) (*domain.ProductCategoryMain, error)
}

type CategoryRepository interface {
	ListCategories(ctx context.Context, params url.Values) (*domain.ProductCategoryListRep, error)
	GetCategory(ctx context.Context, id int64) (*domain.ProductCategoryMain, error)
}
