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
	pathStock = "product_stock"
)

type StockClient struct {
	URL *url.URL
}

func NewStockClient(URL *url.URL) *StockClient {
	// Initialize the StockClient with the base URL
	return &StockClient{
		URL: URL.JoinPath(pathStock),
	}
}

func (r *StockClient) ListStocks(ctx context.Context, params domain.ListParamsSt, productIds, cityIDs []string) (*domain.StockListRep, error) {
	// Format params to url.Values format
	urlParams := r.URL.Query()

	// Set pagination parameters
	if params.Page >= 0 {
		urlParams.Set("list_params.page", strconv.FormatInt(params.Page, 10))
	}

	if params.PageSize > 0 {
		urlParams.Set("list_params.page_size", strconv.FormatInt(params.PageSize, 10))
	}

	// Set sorting parameters
	for _, sortVal := range params.Sort {
		urlParams.Add("list_params.sort", sortVal)
	}

	// Add filters for product and city IDs
	for _, id := range productIds {
		urlParams.Add("product_ids", id)
	}

	for _, cityId := range cityIDs {
		urlParams.Add("city_ids", cityId)
	}

	// Encode the URL parameters
	r.URL.RawQuery = urlParams.Encode()

	// Perform the HTTP GET request
	resp, err := http.Get(r.URL.String())
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// Parse the response into the StockListRep struct
	var stockList *domain.StockListRep
	if err := json.NewDecoder(resp.Body).Decode(&stockList); err != nil {
		return nil, err
	}

	return stockList, nil
}

func (r *StockClient) GetStock(ctx context.Context, productId, cityId string) (*domain.StockMain, error) {
	// Construct the URL for the specific product and city stock
	reqURL := r.URL.JoinPath(productId, cityId).String()

	// Perform the HTTP GET request
	resp, err := http.Get(reqURL)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// Parse the response into the StockMain struct
	var stock *domain.StockMain
	if err := json.NewDecoder(resp.Body).Decode(&stock); err != nil {
		return nil, err
	}

	return stock, nil
}
