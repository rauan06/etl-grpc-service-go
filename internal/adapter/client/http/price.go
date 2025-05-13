package http

import (
	"category/internal/core/domain"
	"context"
	"encoding/json"
	"net/http"
	"net/url"
	"strconv"
)

const (
	pathPrice = "product_price"
)

type PriceClient struct {
	URL *url.URL
}

func NewPriceClient(URL *url.URL) *PriceClient {
	return &PriceClient{
		URL: URL.JoinPath(pathPrice),
	}
}

func (r *PriceClient) ListPrices(ctx context.Context, params domain.ListParamsSt, productIds, cityIDs []string) (*domain.PriceListRep, error) {
	// Format params to url.Values format
	urlParams := r.URL.Query()

	// Set pagination and sorting params
	if params.Page >= 0 {
		urlParams.Set("list_params.page", strconv.FormatInt(params.Page, 10))
	}

	if params.PageSize > 0 {
		urlParams.Set("list_params.page_size", strconv.FormatInt(params.PageSize, 10))
	}

	for _, sortVal := range params.Sort {
		urlParams.Add("list_params.sort", sortVal)
	}

	// Add product and city IDs to the query
	for _, productId := range productIds {
		urlParams.Add("product_ids", productId)
	}

	for _, cityId := range cityIDs {
		urlParams.Add("city_ids", cityId)
	}

	r.URL.RawQuery = urlParams.Encode()

	resp, err := http.Get(r.URL.String())
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var priceList *domain.PriceListRep
	if err := json.NewDecoder(resp.Body).Decode(&priceList); err != nil {
		return nil, err
	}

	return priceList, nil
}

func (r *PriceClient) GetPrice(ctx context.Context, productId, cityId string) (*domain.PriceMain, error) {
	// Construct URL for the specific price request
	url := r.URL.JoinPath(productId, cityId).String()

	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var price *domain.PriceMain
	if err := json.NewDecoder(resp.Body).Decode(&price); err != nil {
		return nil, err
	}

	return price, nil
}
