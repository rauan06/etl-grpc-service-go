package service

import (
	"category/internal/core/domain"
	"category/internal/core/port"
	"category/internal/core/util"
	"context"
	"log/slog"
	"strconv"
	"time"
)

type ProductService struct {
	grpcClient port.ProductClient
	httpClient port.ProductClient
	cache      port.CacheRepository

	logger *slog.Logger
	ctx    context.Context
}

func NewProductService(grpcClient port.ProductClient, httpClient port.ProductClient, cache port.CacheRepository, logger *slog.Logger) ProductService {
	return ProductService{
		grpcClient: grpcClient,
		httpClient: httpClient,
		cache:      cache,
		logger:     logger,
		ctx:        context.Background(),
	}
}

func (s *ProductService) Run(ctx context.Context) {
	s.ctx = ctx

	products := make(chan domain.ProductMain)
	defer close(products)

	go s.CollectProducts(products)

	s.SearchProducts(products)
}

func (s *ProductService) Status() int {
	if s.ctx == nil {
		return domain.StatusNotStarted
	}

	select {
	case <-s.ctx.Done():
		return domain.StatusShutdown
	default:
		return domain.StatusRunning
	}
}

func (s *ProductService) Stop() {
	if s.ctx == nil {
		return
	}

	s.ctx.Done()
	s.logger.InfoContext(s.ctx, "stopped product service gracefully")
}

func (s *ProductService) SearchProducts(products chan<- domain.ProductMain) {
	var page int64
	for {
		params := domain.ListParamsSt{
			Page: strconv.FormatInt(page, 10),
		}

		resp, err := s.fetchProducts(s.ctx, params)
		if err != nil {
			s.logger.WarnContext(s.ctx, "both gRPC and HTTP clients failed", "error", err.Error())
			return
		}

		if len(resp.Results) == 0 {
			select {
			case <-s.ctx.Done():
				return
			case <-time.After(3 * time.Second):
				page = 0
				continue
			}
		}

		// Stream each result
		for _, product := range resp.Results {
			select {
			case products <- product:
			case <-s.ctx.Done():
				return
			}
		}

		page++
	}
}

func (s *ProductService) fetchProducts(ctx context.Context, params domain.ListParamsSt) (*domain.ProductListRep, error) {
	// Attempt to fetch using gRPC
	resp, err := s.grpcClient.ListProducts(ctx, params, []string{}, []string{}, true)
	if err == nil && len(resp.Results) > 0 {
		return resp, nil
	}

	// If gRPC fails or returns no results, try HTTP
	s.logger.WarnContext(ctx, "gRPC failed, retrying with HTTP", "error", err.Error())
	resp, err = s.httpClient.ListProducts(ctx, params, []string{}, []string{}, true)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (s *ProductService) CollectProducts(products <-chan domain.ProductMain) {
	for product := range products {
		if err := s.cacheProduct(product); err != nil {
			s.logger.ErrorContext(s.ctx, "error caching product", "productID", product.ID, "error", err.Error())
		}
	}
}

func (s *ProductService) cacheProduct(product domain.ProductMain) error {
	data, err := util.Serialize(product)
	if err != nil {
		return err
	}

	key := util.GenerateCacheKey("product", product.ID)
	if err := s.cache.Set(s.ctx, key, data); err != nil {
		return err
	}

	return nil
}
