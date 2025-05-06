package port

import (
	"context"
	"net/url"

	"category/internal/core/domain"
)

type CityService interface {
	ListCities(ctx context.Context, params url.Values) (*domain.ProductCityListRep, error)
	GetCity(ctx context.Context, id int64) (*domain.ProductCityMain, error)
}

type CityRepository interface {
	ListCities(ctx context.Context, params url.Values) (*domain.ProductCityListRep, error)
	GetCity(ctx context.Context, id int64) (*domain.ProductCityMain, error)
}
