package service

import (
	"context"
	"errors"
	"net/url"

	"category/internal/core/domain"
	"category/internal/core/port"
)

type CityService struct {
	repo port.CityRepository
}

func NewCityService(repo port.CityRepository) *CityService {
	return &CityService{repo: repo}
}

func (s *CityService) ListCategories(ctx context.Context, params url.Values) (*domain.ProductCityListRep, error) {
	params["list_params.sort"] = filterValidSortParams(params["sort"])
	params["list_params.ids"] = filterValidIDs(params["ids"])
	params["list_params.page"] = params["page"]
	params["list_params.page_size"] = params["page_size"]

	clearParams(params)

	return s.repo.ListCities(ctx, params)
}

func (s *CityService) GetCity(ctx context.Context, id int64) (*domain.ProductCityMain, error) {
	if id < 0 {
		return nil, errors.New("id cannot be negative")
	}

	return s.repo.GetCity(ctx, id)
}
