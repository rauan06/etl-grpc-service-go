package service

import (
	"category/internal/core/domain"
	"category/internal/core/port"
	"category/internal/core/util"
	"context"
	"log/slog"
	"maps"
	"sync"
	"sync/atomic"
	"time"
)

type CollectorService struct {
	svc    ProductService
	cache  port.CacheRepository
	logger *slog.Logger
	ctx    context.Context
	cancel context.CancelFunc

	status      int
	ProductList struct {
		listProducts map[string]domain.ProductMain
		loaded       atomic.Bool
		mu           *sync.Mutex
	}
}

func NewCollectorService(cache port.CacheRepository, logger *slog.Logger) CollectorService {
	ctx, cancel := context.WithCancel(context.Background())

	return CollectorService{
		cache:  cache,
		logger: logger,
		ctx:    ctx,
		cancel: cancel,
		status: domain.StatusNotStarted,

		ProductList: struct {
			listProducts map[string]domain.ProductMain
			loaded       atomic.Bool
			mu           *sync.Mutex
		}{
			listProducts: map[string]domain.ProductMain{},
			mu:           &sync.Mutex{},
		},
	}
}

func (s *CollectorService) Run() {
	go s.startLoadingProducts()

	s.status = domain.StatusRunning

	s.collectProductDetails()

	s.status = domain.StatusShutdown
}

func (s *CollectorService) Status() int {
	return s.status
}

func (s *CollectorService) Stop() {
	s.cancel()
	s.status = domain.StatusShutdown
	s.logger.InfoContext(s.ctx, "stopped collector service gracefully")
}

func (s *CollectorService) collectProductDetails() {
	for {
		select {
		case <-s.ctx.Done():
			return
		default:
			if s.ProductList.loaded.Load() {
				// Need to collect 5 data
				// 1. Validate product data
				// 2. Collect categoories
				// 3. Collect cities
				// 4. Collect stocks
				// 5. Collect prices
			}
			time.Sleep(3 * time.Second)
		}
	}
}

func (s *CollectorService) getCache(serviceName, id string) ([]byte, error) {
	key := util.GenerateCacheKey(serviceName, id)
	return s.cache.Get(s.ctx, key)
}

func (s *CollectorService) startLoadingProducts() {
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			s.loadProducts()
		case <-s.ctx.Done():
			return
		}
	}
}

func (s *CollectorService) loadProducts() {
	s.ProductList.mu.Lock()
	defer s.ProductList.mu.Unlock()

	s.svc.ProductList.Mu.Lock()
	defer s.svc.ProductList.Mu.Unlock()

	maps.Copy(s.ProductList.listProducts, s.svc.ProductList.ListProducts)
	s.ProductList.loaded.Store(true)
}
