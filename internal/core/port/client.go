package port

import (
	"context"

	"github.com/rauan06/etl-grpc-service-go/internal/core/domain"
)

type CLientI interface {
	ListCategories(ctx context.Context, params domain.ListParamsSt, ids []string) (*domain.CategoryListRep, error)
	GetCategory(ctx context.Context, id string) (*domain.CategoryMain, error)
	ListCities(ctx context.Context, params domain.ListParamsSt, ids []string) (*domain.CityListRep, error)
	GetCity(ctx context.Context, id string) (*domain.CityMain, error)
	ListProducts(ctx context.Context, params domain.ListParamsSt, ids, categoryIDs []string, withCategory bool) (*domain.ProductListRep, error)
	GetProduct(ctx context.Context, id string) (*domain.ProductMain, error)

	ListPrices(ctx context.Context, params domain.ListParamsSt, productIds, cityIds []string) (*domain.PriceListRep, error)
	GetPrice(ctx context.Context, productId, cityId string) (*domain.PriceMain, error)
	ListStocks(ctx context.Context, params domain.ListParamsSt, productIds, cityIds []string) (*domain.StockListRep, error)
	GetStock(ctx context.Context, productId, cityId string) (*domain.StockMain, error)
}
