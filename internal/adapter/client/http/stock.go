package http

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"category/internal/core/domain"
)

type StockClient struct {
	client  *http.Client
	baseURL string
}

func NewStockClient(ctx context.Context, baseURL string) (*StockClient, error) {
	client := &http.Client{}
	return &StockClient{
		client:  client,
		baseURL: baseURL,
	}, nil
}

func (c *StockClient) Close() {
	// No connection to close in HTTP, but you can manage client lifecycle if needed.
}

// ListStores fetches the list of stores.
func (c *StockClient) ListStores(ctx context.Context, params domain.ListParamsSt, productIds, cityIDs []string) (*domain.StockListRep, error) {
	page, err := strconv.ParseInt(params.Page, 10, 64)
	if err != nil {
		return nil, domain.ErrParseInt64
	}
	pageSize, err := strconv.ParseInt(params.PageSize, 10, 64)
	if err != nil {
		return nil, domain.ErrParseInt64
	}

	// Create the request body
	reqBody := map[string]interface{}{
		"listParams": map[string]interface{}{
			"page":     page,
			"pageSize": pageSize,
			"sort":     params.Sort,
		},
		"productIds": productIds,
		"cityIds":    cityIDs,
	}
	body, err := json.Marshal(reqBody)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request body: %w", err)
	}

	// Send the HTTP request
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.baseURL+"/product-stocks/list", bytes.NewBuffer(body))
	if err != nil {
		return nil, fmt.Errorf("failed to create HTTP request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send HTTP request: %w", err)
	}
	defer resp.Body.Close()

	// Check the response status code
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("received non-200 response: %s", resp.Status)
	}

	// Parse the response body
	var respBody struct {
		PaginationInfo struct {
			Page     int64 `json:"page"`
			PageSize int64 `json:"page_size"`
		} `json:"paginationInfo"`
		Results []struct {
			ProductId string `json:"productId"`
			CityId    string `json:"cityId"`
			Value     int64  `json:"value"`
		} `json:"results"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&respBody); err != nil {
		return nil, fmt.Errorf("failed to decode response body: %w", err)
	}

	// Build the response
	var results []domain.StockMain
	for _, prod := range respBody.Results {
		results = append(results, domain.StockMain{
			ProductId: prod.ProductId,
			CityId:    prod.CityId,
			Value:     prod.Value,
		})
	}

	return &domain.StockListRep{
		PaginationInfo: domain.PaginationInfoSt{
			Page:     strconv.FormatInt(respBody.PaginationInfo.Page, 10),
			PageSize: strconv.FormatInt(respBody.PaginationInfo.PageSize, 10),
		},
		Results: results,
	}, nil
}

// GetStore fetches a single store by productId and cityId.
func (c *StockClient) GetStore(ctx context.Context, productId, cityId string) (*domain.StockMain, error) {
	// Create the request body
	reqBody := map[string]string{
		"productId": productId,
		"cityId":    cityId,
	}
	body, err := json.Marshal(reqBody)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request body: %w", err)
	}

	// Send the HTTP request
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.baseURL+"/product-stocks/get", bytes.NewBuffer(body))
	if err != nil {
		return nil, fmt.Errorf("failed to create HTTP request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send HTTP request: %w", err)
	}
	defer resp.Body.Close()

	// Check the response status code
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("received non-200 response: %s", resp.Status)
	}

	// Parse the response body
	var respBody struct {
		ProductId string `json:"productId"`
		CityId    string `json:"cityId"`
		Value     int64  `json:"value"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&respBody); err != nil {
		return nil, fmt.Errorf("failed to decode response body: %w", err)
	}

	// Build and return the result
	return &domain.StockMain{
		ProductId: respBody.ProductId,
		CityId:    respBody.CityId,
		Value:     respBody.Value,
	}, nil
}
