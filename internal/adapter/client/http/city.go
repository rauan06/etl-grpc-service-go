package repository

import (
	"context"
	"encoding/json"
	"net/http"
	"net/url"

	"category/internal/core/domain"
)

const (
	pathCity = "city"
)

type CityClient struct {
	URL *url.URL
}

func NewCityClient(URL *url.URL) *CityClient {
	return &CityClient{
		URL.JoinPath(pathCity),
	}
}

func (r *CityClient) ListCategories(ctx context.Context, params url.Values) (*domain.CityListRep, error) {
	resp, err := http.Get(r.URL.String())
	if err != nil {
		return nil, err
	}

	var categories *domain.CityListRep
	if err := json.NewDecoder(resp.Body).Decode(categories); err != nil {
		return nil, err
	}

	return categories, nil
}

func (r *CityClient) GetCity(ctx context.Context, id int64) (*domain.CityMain, error) {
	resp, err := http.Get(r.URL.String())
	if err != nil {
		return nil, err
	}

	var categories *domain.CityMain
	if err := json.NewDecoder(resp.Body).Decode(categories); err != nil {
		return nil, err
	}

	return categories, nil
}
