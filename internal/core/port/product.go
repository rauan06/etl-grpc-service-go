package port

import (
	"context"
	"net/url"

	"category/internal/core/domain"
)

type ProductService interface {
	ListProducts(ctx context.Context, params url.Values) (*domain.ProductListRep, error)
	GetProduct(ctx context.Context, id int64) (*domain.ProductMain, error)
}

type ProductClient interface {
	ListProducts(ctx context.Context, params url.Values) (*domain.ProductListRep, error)
	GetProduct(ctx context.Context, id int64) (*domain.ProductMain, error)
}
