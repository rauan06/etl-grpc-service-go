package port

import (
	"context"
	"net/url"

	"category/internal/core/domain"
)

type CityService interface {
	GetCityList(ctx context.Context, params url.Values) (*domain.ProductCityListRep, error)
	GetCityByID(ctx context.Context, id int64) (*domain.ProductCityMain, error)
}

type CityRepository interface {
	GetCityList(ctx context.Context, params url.Values) (*domain.ProductCityListRep, error)
	GetCityByID(ctx context.Context, id int64) (*domain.ProductCityMain, error)
}
