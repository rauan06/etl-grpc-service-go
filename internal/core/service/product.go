package service

import (
	"category/internal/core/domain"
	"category/internal/core/port"
	"context"
	"errors"
	"log/slog"
	"sync"
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

	ProductList struct {
		ListProducts map[string]domain.ProductMain
		Mu           *sync.Mutex
	}
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

		ProductList: struct {
			ListProducts map[string]domain.ProductMain
			Mu           *sync.Mutex
		}{
			map[string]domain.ProductMain{},
			&sync.Mutex{},
		},
	}
}

func (s *ProductService) Run() {
	s.ctx, s.cancel = context.WithCancel(s.ctx)

	products := make(chan domain.ProductMain)
	defer close(products)

	go s.CollectProducts(products)

	s.status = domain.StatusRunning

	s.logger.Info("product service has started")
	s.SearchProducts(products)

	s.status = domain.StatusShutdown
}

func (s *ProductService) Status() int {
	return s.status
}

func (s *ProductService) Stop() {
	s.cancel()
	s.status = domain.StatusShutdown
	s.logger.InfoContext(s.ctx, "stopped product service gracefully")
}

func (s *ProductService) SearchProducts(products chan<- domain.ProductMain) {
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
	s.ProductList.Mu.Lock()
	defer s.ProductList.Mu.Unlock()

	s.ProductList.ListProducts[product.ID] = product

	return nil
}
