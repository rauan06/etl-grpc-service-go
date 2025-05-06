package port

import (
	"category/internal/core/domain"
	"context"
)

type CityService interface {
	GetCityList(ctx context.Context, params domain.ProductListParamsSt, ids []string) (domain.ProductCityListRep, error)
	GetCityByID(ctx context.Context, id string) (domain.ProductCityMain, error)
}

type CityRepository interface {
	GetCityList(ctx context.Context, params domain.ProductListParamsSt, ids []string) ([]domain.ProductCityMain, error)
	GetCityByID(ctx context.Context, id string) (domain.ProductCityMain, error)
}
