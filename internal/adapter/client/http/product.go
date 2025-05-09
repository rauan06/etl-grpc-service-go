package repository

import (
	"category/internal/core/domain"
	"context"
	"encoding/json"
	"net/http"
	"net/url"
)

const (
	pathProduct = "product"
)

type productRepository struct {
	URL *url.URL
}

func NewProductRepository(URL *url.URL) *productRepository {
	return &productRepository{
		URL.JoinPath(pathProduct),
	}
}

func (r *productRepository) ListProducts(ctx context.Context, params url.Values) (*domain.ProductProductListRep, error) {
	resp, err := http.Get(r.URL.String())
	if err != nil {
		return nil, err
	}

	var products *domain.ProductProductListRep
	if err := json.NewDecoder(resp.Body).Decode(products); err != nil {
		return nil, err
	}

	return products, nil
}

func (r *productRepository) GetProduct(ctx context.Context, id int64) (*domain.ProductProductMain, error) {
	resp, err := http.Get(r.URL.String())
	if err != nil {
		return nil, err
	}

	var products *domain.ProductProductMain
	if err := json.NewDecoder(resp.Body).Decode(products); err != nil {
		return nil, err
	}

	return products, nil
}
