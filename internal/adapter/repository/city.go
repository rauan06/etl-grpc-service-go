package repository

import (
	"context"
	"encoding/json"
	"net/http"
	"net/url"

	"category/internal/core/domain"
)

const (
	pathCity = "City"
)

type CityRepository struct {
	URL *url.URL
}

func NewCityRepository(URL *url.URL) *CityRepository {
	return &CityRepository{
		URL.JoinPath(pathCity),
	}
}

func (r *CityRepository) ListCategories(ctx context.Context, params url.Values) (*domain.ProductCityListRep, error) {
	resp, err := http.Get(r.URL.String())
	if err != nil {
		return nil, err
	}

	var categories *domain.ProductCityListRep
	if err := json.NewDecoder(resp.Body).Decode(categories); err != nil {
		return nil, err
	}

	return categories, nil
}

func (r *CityRepository) GetCity(ctx context.Context, id int64) (*domain.ProductCityMain, error) {
	resp, err := http.Get(r.URL.String())
	if err != nil {
		return nil, err
	}

	var categories *domain.ProductCityMain
	if err := json.NewDecoder(resp.Body).Decode(categories); err != nil {
		return nil, err
	}

	return categories, nil
}
