package port

import (
	"context"
	"net/url"

	"category/internal/core/domain"
)

type ProductService interface {
	ListProducts(ctx context.Context, params url.Values) (*domain.ProductProductListRep, error)
	GetProduct(ctx context.Context, id int64) (*domain.ProductProductMain, error)
}

type ProductRepository interface {
	ListProducts(ctx context.Context, params url.Values) (*domain.ProductProductListRep, error)
	GetProduct(ctx context.Context, id int64) (*domain.ProductProductMain, error)
}
