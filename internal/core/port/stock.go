package port

import (
	"context"

	"github.com/rauan06/etl-grpc-service-go/internal/core/domain"
)

type StockClient interface {
	ListStocks(ctx context.Context, params domain.ListParamsSt, productIds, cityIds []string) (*domain.StockListRep, error)
	GetStock(ctx context.Context, productId, cityId string) (*domain.StockMain, error)
}
