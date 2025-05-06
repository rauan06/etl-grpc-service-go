package repository

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

type CategoryRepository struct {
	URL *url.URL
}

func NewCategoryRepository(URL *url.URL) *CategoryRepository {
	return &CategoryRepository{
		URL.JoinPath(pathCategory),
	}
}

func (r *CategoryRepository) ListCategories(ctx context.Context, params url.Values) (*domain.ProductCategoryListRep, error) {
	resp, err := http.Get(r.URL.String())
	if err != nil {
		return nil, err
	}

	var categories *domain.ProductCategoryListRep
	if err := json.NewDecoder(resp.Body).Decode(categories); err != nil {
		return nil, err
	}

	return categories, nil
}

func (r *CategoryRepository) GetCategory(ctx context.Context, id int64) (*domain.ProductCategoryMain, error) {
	resp, err := http.Get(r.URL.String())
	if err != nil {
		return nil, err
	}

	var categories *domain.ProductCategoryMain
	if err := json.NewDecoder(resp.Body).Decode(categories); err != nil {
		return nil, err
	}

	return categories, nil
}
