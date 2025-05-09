package client

import (
	"context"
	"encoding/json"
	"net/http"
	"net/url"

	"category/internal/core/domain"
)

const (
	pathCategory = "category"
)

type CategoryClient struct {
	URL *url.URL
}

func NewCategoryClient(URL *url.URL) *CategoryClient {
	return &CategoryClient{
		URL.JoinPath(pathCategory),
	}
}

func (r *CategoryClient) ListCategories(ctx context.Context, params url.Values) (*domain.CategoryListRep, error) {
	resp, err := http.Get(r.URL.String())
	if err != nil {
		return nil, err
	}

	var categories *domain.CategoryListRep
	if err := json.NewDecoder(resp.Body).Decode(categories); err != nil {
		return nil, err
	}

	return categories, nil
}

func (r *CategoryClient) GetCategory(ctx context.Context, id int64) (*domain.CategoryMain, error) {
	resp, err := http.Get(r.URL.String())
	if err != nil {
		return nil, err
	}

	var categories *domain.CategoryMain
	if err := json.NewDecoder(resp.Body).Decode(categories); err != nil {
		return nil, err
	}

	return categories, nil
}
