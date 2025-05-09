package port

import (
	"context"
	"net/url"

	"category/internal/core/domain"
)

type CityService interface {
	ListCities(ctx context.Context, params url.Values) (*domain.CityListRep, error)
	GetCity(ctx context.Context, id int64) (*domain.CityMain, error)
}

type CityClient interface {
	ListCities(ctx context.Context, params url.Values) (*domain.CityListRep, error)
	GetCity(ctx context.Context, id int64) (*domain.CityMain, error)
}
