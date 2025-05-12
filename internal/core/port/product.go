package port

import (
	"context"

	"category/internal/core/domain"
)

type ProductClient interface {
	ListProducts(ctx context.Context, params domain.ListParamsSt, ids, categoryIDs []string, withCategory bool) (*domain.ProductListRep, error)
	GetProduct(ctx context.Context, id string) (*domain.ProductMain, error)
}
