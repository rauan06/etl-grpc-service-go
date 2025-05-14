package port

import (
	"category/internal/core/domain"
	"context"
)

type PriceClient interface {
	ListPrices(ctx context.Context, params domain.ListParamsSt, productIds, cityIds []string) (*domain.PriceListRep, error)
	GetPrice(ctx context.Context, productId, cityId string) (*domain.PriceMain, error)
}
