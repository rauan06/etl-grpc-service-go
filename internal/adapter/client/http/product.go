package http

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

func NewProductClient(URL *url.URL) *productClient {
	return &productClient{
		URL.JoinPath(pathProduct),
	}
}

func (r *productClient) ListProducts(ctx context.Context, params domain.ListParamsSt, ids, categoryIDs []string, withCategory bool) (*domain.ProductListRep, error) {
	// Format params to url.Values format
	urlParams := r.URL.Query()

	if params.Page != "" {
		urlParams.Set("list_params.page", params.Page)
	}

	if params.PageSize != "" {
		urlParams.Set("list_params.page_size", params.PageSize)
	}

	if withCategory {
		urlParams.Set("with_category", "true")
	} else {
		urlParams.Set("with_category", "false")
	}

	for _, sortVal := range params.Sort {
		urlParams.Add("list_params.sort", sortVal)
	}

	for _, id := range ids {
		urlParams.Add("ids", id)
	}

	for _, categoryVal := range categoryIDs {
		urlParams.Add("category_ids", categoryVal)
	}

	r.URL.RawQuery = urlParams.Encode()

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

func (r *productClient) GetProduct(ctx context.Context, id string) (*domain.ProductMain, error) {
	resp, err := http.Get(r.URL.JoinPath(id).String())
	if err != nil {
		return nil, err
	}

	var products = &domain.ProductMain{}
	if err := json.NewDecoder(resp.Body).Decode(products); err != nil {
		return nil, err
	}

	return products, nil
}
