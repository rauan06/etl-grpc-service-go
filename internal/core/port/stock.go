package port

import (
	"context"

	"category/internal/core/domain"
)

type Stocklient interface {
	ListStocks(ctx context.Context, params domain.ListParamsSt, productIds, cityIds []string) (*domain.StockListRep, error)
	GetStock(ctx context.Context, productId, cityId string) (*domain.StockMain, error)
}
