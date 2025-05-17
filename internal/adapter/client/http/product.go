package http

import (
	"context"
	"encoding/json"
	"net/http"
	"net/url"
	"strconv"

	"github.com/rauan06/etl-grpc-service-go/internal/core/domain"
)

const (
	pathProduct = "product"
)

type ProductClient struct {
	URL *url.URL
}

func NewProductClient(URL *url.URL) *ProductClient {
	return &ProductClient{
		URL.JoinPath(pathProduct),
	}
}

func (r *ProductClient) ListProducts(ctx context.Context, params domain.ListParamsSt, ids, categoryIDs []string, withCategory bool) (*domain.ProductListRep, error) {
	// Format params to url.Values format
	urlParams := r.URL.Query()

	if params.Page >= 0 {
		urlParams.Set("list_params.page", strconv.FormatInt(params.Page, 10))
	}

	if params.PageSize > 0 {
		urlParams.Set("list_params.page_size", strconv.FormatInt(params.PageSize, 10))
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

func (r *ProductClient) GetProduct(ctx context.Context, id string) (*domain.ProductMain, error) {
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
