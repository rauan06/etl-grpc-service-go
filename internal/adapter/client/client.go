package client

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	"github.com/rauan06/etl-grpc-service-go/internal/adapter/client/grpc"
	"github.com/rauan06/etl-grpc-service-go/internal/adapter/client/http"
	"github.com/rauan06/etl-grpc-service-go/internal/core/domain"
	"github.com/rauan06/etl-grpc-service-go/internal/core/port"
)

// Client is a structure that implements port.ClientI interface
type Client struct {
	GRPC   *grpc.Client
	HTTP   *http.Client
	logger *slog.Logger
}

var _ port.CLientI = (*Client)(nil)

// NewClient accepts http and grpc client instances, and a logger.
// Returns new instance of Client
func NewClient(
	grpc *grpc.Client,
	http *http.Client,
	logger *slog.Logger,
) *Client {
	return &Client{
		GRPC:   grpc,
		HTTP:   http,
		logger: logger,
	}
}

// ListProducts sends get request to the /product endpoint via grpc
// uses http if grpc fails
func (c *Client) ListProducts(ctx context.Context, params domain.ListParamsSt, ids, categoryIDs []string, withCategory bool) (*domain.ProductListRep, error) {
	// Validate context
	if ctx == nil {
		return nil, errors.New("nil context")
	}

	// Initialize empty slices if nil
	if ids == nil {
		ids = []string{}
	}
	if categoryIDs == nil {
		categoryIDs = []string{}
	}

	var resp *domain.ProductListRep
	var err error

	// Try GRPC if available
	if c.GRPC != nil && c.GRPC.Product != nil {
		resp, err = c.GRPC.Product.ListProducts(ctx, params, ids, categoryIDs, withCategory)
		if err == nil {
			return resp, nil
		}
		c.logger.Warn("GRPC ListProducts failed", "error", err)
	} else {
		c.logger.Warn("GRPC client not available")
	}

	// Fallback to HTTP if available
	if c.HTTP != nil && c.HTTP.Product != nil {
		resp, err = c.HTTP.Product.ListProducts(ctx, params, ids, categoryIDs, withCategory)
		if err != nil {
			return nil, fmt.Errorf("all clients failed: last error: %w", err)
		}
		return resp, nil
	}

	return nil, errors.New("no available client could handle the request")
}

func (c *Client) ListCategories(ctx context.Context, params domain.ListParamsSt, ids []string) (*domain.CategoryListRep, error) {
	resp, err := c.GRPC.Category.ListCategories(ctx, params, ids)
	if err != nil {
		c.logger.Warn("errors listing categories via grpc, trying http", "error", err)

		resp, err = c.HTTP.Category.ListCategories(ctx, params, ids)
		if err != nil {
			c.logger.Error("error listing categories both using http and grpc", "error", err)
			return nil, err
		}
	}
	return resp, nil
}

func (c *Client) GetCategory(ctx context.Context, id string) (*domain.CategoryMain, error) {
	resp, err := c.GRPC.Category.GetCategory(ctx, id)
	if err != nil {
		c.logger.Warn("errors getting category via grpc, trying http", "error", err)

		resp, err = c.HTTP.Category.GetCategory(ctx, id)
		if err != nil {
			c.logger.Error("error getting category both using http and grpc", "error", err)
			return nil, err
		}
	}
	return resp, nil
}

func (c *Client) ListCities(ctx context.Context, params domain.ListParamsSt, ids []string) (*domain.CityListRep, error) {
	resp, err := c.GRPC.City.ListCities(ctx, params, ids)
	if err != nil {
		c.logger.Warn("errors listing cities via grpc, trying http", "error", err)

		resp, err = c.HTTP.City.ListCities(ctx, params, ids)
		if err != nil {
			c.logger.Error("error listing cities both using http and grpc", "error", err)
			return nil, err
		}
	}
	return resp, nil
}

func (c *Client) GetCity(ctx context.Context, id string) (*domain.CityMain, error) {
	resp, err := c.GRPC.City.GetCity(ctx, id)
	if err != nil {
		c.logger.Warn("errors getting city via grpc, trying http", "error", err)

		resp, err = c.HTTP.City.GetCity(ctx, id)
		if err != nil {
			c.logger.Error("error getting city both using http and grpc", "error", err)
			return nil, err
		}
	}
	return resp, nil
}

// func (c *Client) ListProducts(ctx context.Context, params domain.ListParamsSt, ids, categoryIDs []string, withCategory bool) (*domain.ProductListRep, error) {
// 	resp, err := c.GRPC.Product.ListProducts(ctx, params, ids, categoryIDs, withCategory)
// 	if err != nil {
// 		c.logger.Warn("errors listing products via grpc, trying http", "error", err)

// 		resp, err = c.HTTP.Product.ListProducts(ctx, params, ids, categoryIDs, withCategory)
// 		if err != nil {
// 			c.logger.Error("error listing products both using http and grpc", "error", err)
// 			return nil, err
// 		}
// 	}
// 	return resp, nil
// }

func (c *Client) GetProduct(ctx context.Context, id string) (*domain.ProductMain, error) {
	resp, err := c.GRPC.Product.GetProduct(ctx, id)
	if err != nil {
		c.logger.Warn("errors getting product via grpc, trying http", "error", err)

		resp, err = c.HTTP.Product.GetProduct(ctx, id)
		if err != nil {
			c.logger.Error("error getting product both using http and grpc", "error", err)
			return nil, err
		}
	}
	return resp, nil
}

func (c *Client) ListPrices(ctx context.Context, params domain.ListParamsSt, productIDs, cityIDs []string) (*domain.PriceListRep, error) {
	resp, err := c.GRPC.Price.ListPrices(ctx, params, productIDs, cityIDs)
	if err != nil {
		c.logger.Warn("errors listing prices via grpc, trying http", "error", err)

		resp, err = c.HTTP.Price.ListPrices(ctx, params, productIDs, cityIDs)
		if err != nil {
			c.logger.Error("error listing prices both using http and grpc", "error", err)
			return nil, err
		}
	}
	return resp, nil
}

func (c *Client) GetPrice(ctx context.Context, productID, cityID string) (*domain.PriceMain, error) {
	resp, err := c.GRPC.Price.GetPrice(ctx, productID, cityID)
	if err != nil {
		c.logger.Warn("errors getting price via grpc, trying http", "error", err)

		resp, err = c.HTTP.Price.GetPrice(ctx, productID, cityID)
		if err != nil {
			c.logger.Error("error getting price both using http and grpc", "error", err)
			return nil, err
		}
	}
	return resp, nil
}

func (c *Client) ListStocks(ctx context.Context, params domain.ListParamsSt, productIDs, cityIDs []string) (*domain.StockListRep, error) {
	resp, err := c.GRPC.Stock.ListStocks(ctx, params, productIDs, cityIDs)
	if err != nil {
		c.logger.Warn("errors listing stocks via grpc, trying http", "error", err)

		resp, err = c.HTTP.Stock.ListStocks(ctx, params, productIDs, cityIDs)
		if err != nil {
			c.logger.Error("error listing stocks both using http and grpc", "error", err)
			return nil, err
		}
	}
	return resp, nil
}

func (c *Client) GetStock(ctx context.Context, productID, cityID string) (*domain.StockMain, error) {
	resp, err := c.GRPC.Stock.GetStock(ctx, productID, cityID)
	if err != nil {
		c.logger.Warn("errors getting stock via grpc, trying http", "error", err)

		resp, err = c.HTTP.Stock.GetStock(ctx, productID, cityID)
		if err != nil {
			c.logger.Error("error getting stock both using http and grpc", "error", err)
			return nil, err
		}
	}
	return resp, nil
}
