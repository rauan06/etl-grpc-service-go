package port

import (
	"context"
	"net/url"

	"category/internal/core/domain"
)

type CategoryService interface {
	ListCategories(ctx context.Context, params url.Values) (*domain.CategoryListRep, error)
	GetCategory(ctx context.Context, id int64) (*domain.CategoryMain, error)
}

type CategoryClient interface {
	ListCategories(ctx context.Context, params url.Values) (*domain.CategoryListRep, error)
	GetCategory(ctx context.Context, id int64) (*domain.CategoryMain, error)
}
