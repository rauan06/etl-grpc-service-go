package service

import (
	"category/internal/core/domain"
	"category/internal/core/port"
	"category/internal/core/util"
	"context"
	"errors"
	"log/slog"
	"time"
)

type ProductService struct {
	grpcClient port.ProductClient
	httpClient port.ProductClient
	cache      port.CacheRepository

	logger *slog.Logger
	ctx    context.Context
	cancel context.CancelFunc

	status int
}

func NewProductService(grpcClient port.ProductClient, httpClient port.ProductClient, cache port.CacheRepository, logger *slog.Logger) ProductService {
	ctx, cancel := context.WithCancel(context.Background())

	return ProductService{
		grpcClient: grpcClient,
		httpClient: httpClient,
		cache:      cache,
		logger:     logger,
		ctx:        ctx,
		cancel:     cancel,
		status:     domain.StatusNotStarted,
	}
}

func (s *ProductService) Run() {
	if s.Status() == domain.StatusRunning {
		return
	}

	s.ctx, s.cancel = context.WithCancel(context.Background())

	products := make(chan domain.ProductMain)

	go s.CollectProducts(products)
	go s.SearchProducts(products)

	s.status = domain.StatusRunning

	s.logger.Info("product service has started")

}

func (s *ProductService) GetServiceName() string {
	return "Product"
}

func (s *ProductService) Status() int {
	return s.status
}

func (s *ProductService) Stop() {
	if s.status == domain.StatusShutdown {
		return
	}

	s.cancel()

	s.status = domain.StatusShutdown

	s.logger.InfoContext(s.ctx, "stopped product service gracefully")
}

func (s *ProductService) SearchProducts(products chan<- domain.ProductMain) {
	defer close(products)

	var page int64 = 0
	for {
		params := domain.ListParamsSt{
			Page: page,
		}

		resp, err := s.fetchProducts(s.ctx, params)
		if err != nil {
			select {
			case <-s.ctx.Done():
				return
			case <-time.After(3 * time.Second):
				page = 0
				continue
			}
		}

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
	resp, err := s.grpcClient.ListProducts(ctx, params, []string{}, []string{}, true)
	if err == nil && len(resp.Results) > 0 {
		return resp, nil
	}

	resp, err = s.httpClient.ListProducts(ctx, params, []string{}, []string{}, true)
	if err != nil || resp.Results == nil {
		return nil, errors.New("null response results")
	}

	return resp, nil
}

func (s *ProductService) CollectProducts(products <-chan domain.ProductMain) {
	for product := range products {
		if product == (domain.ProductMain{}) {
			continue
		}

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
