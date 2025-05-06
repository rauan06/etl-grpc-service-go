package port

import (
	"category/internal/core/domain"
	"context"
)

type ProductService interface {
	GetProductList(ctx context.Context, params domain.ProductListParamsSt, ids, categoryIDs []string, withCategory bool) (domain.ProductProductListRep, error)
	GetProductByID(ctx context.Context, id string, withCategory bool) (domain.ProductProductMain, error)
}

type ProductRepository interface {
	GetProductList(ctx context.Context, params domain.ProductListParamsSt, ids, categoryIDs []string, withCategory bool) ([]domain.ProductProductMain, error)
	GetProductByID(ctx context.Context, id string, withCategory bool) (domain.ProductProductMain, error)
}
