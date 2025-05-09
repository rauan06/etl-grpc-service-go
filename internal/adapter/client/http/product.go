package client

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

type productClient struct {
	URL *url.URL
}

func NewproductClient(URL *url.URL) *productClient {
	return &productClient{
		URL.JoinPath(pathProduct),
	}
}

func (r *productClient) ListProducts(ctx context.Context, params url.Values) (*domain.ProductListRep, error) {
	resp, err := http.Get(r.URL.String())
	if err != nil {
		return nil, err
	}

	var products *domain.ProductListRep
	if err := json.NewDecoder(resp.Body).Decode(products); err != nil {
		return nil, err
	}

	return products, nil
}

func (r *productClient) Getproduct(ctx context.Context, id int64) (*domain.ProductMain, error) {
	resp, err := http.Get(r.URL.String())
	if err != nil {
		return nil, err
	}

	var products *domain.ProductMain
	if err := json.NewDecoder(resp.Body).Decode(products); err != nil {
		return nil, err
	}

	return products, nil
}
