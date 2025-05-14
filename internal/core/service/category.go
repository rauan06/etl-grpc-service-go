package service

import (
	"context"
	"log/slog"
	"time"

	"category/internal/core/domain"
	"category/internal/core/port"
	"category/internal/core/util"
)

type CategoryService struct {
	grpcClient port.CategoryClient
	httpClient port.CategoryClient
	cache      port.CacheRepository

	logger *slog.Logger
	ctx    context.Context
	cancel context.CancelFunc

	status int
}

func NewCategoryService(grpcClient port.CategoryClient, httpClient port.CategoryClient, cache port.CacheRepository, logger *slog.Logger) CategoryService {
	ctx, cancel := context.WithCancel(context.Background()) // <- create cancelable context

	return CategoryService{
		grpcClient: grpcClient,
		httpClient: httpClient,
		cache:      cache,
		logger:     logger,
		ctx:        ctx,
		cancel:     cancel,
		status:     domain.StatusNotStarted,
	}
}

func (s *CategoryService) Run() {
	if s.Status() == domain.StatusRunning {
		return
	}

	categories := make(chan domain.CategoryMain)

	go s.CollectCategories(categories)

	s.status = domain.StatusRunning

	s.logger.Info("category service has started")
	go s.SearchCategories(categories)
}

func (s *CategoryService) Status() int {
	return s.status
}

func (s *CategoryService) Stop() {
	s.cancel()

	s.status = domain.StatusShutdown

	s.logger.InfoContext(s.ctx, "stopped category service gracefully")
}

func (s *CategoryService) SearchCategories(categories chan<- domain.CategoryMain) {
	defer close(categories)

	var page int64
	for {
		params := domain.ListParamsSt{
			Page: page,
		}

		resp, err := s.fetchCategories(s.ctx, params)
		if err != nil {
			select {
			case <-s.ctx.Done():
				return
			case <-time.After(3 * time.Second):
				page = 0
				continue
			}
		}

		// Stream each result
		for _, category := range resp.Results {
			select {
			case categories <- category:
			case <-s.ctx.Done():
				return
			}
		}

		page++
	}
}

func (s *CategoryService) fetchCategories(ctx context.Context, params domain.ListParamsSt) (*domain.CategoryListRep, error) {
	// Attempt to fetch using gRPC
	resp, err := s.grpcClient.ListCategories(ctx, params, []string{})
	if err == nil && len(resp.Results) > 0 {
		return resp, nil
	}

	resp, err = s.httpClient.ListCategories(ctx, params, []string{})
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (s *CategoryService) CollectCategories(categories <-chan domain.CategoryMain) {
	for category := range categories {
		if err := s.cacheCategory(category); err != nil {
			s.logger.ErrorContext(s.ctx, "error caching category", "categoryID", category.ID, "error", err.Error())
		}
	}
}

func (s *CategoryService) cacheCategory(category domain.CategoryMain) error {
	data, err := util.Serialize(category)
	if err != nil {
		return err
	}

	key := util.GenerateCacheKey("category", category.ID)
	if err := s.cache.Set(s.ctx, key, data); err != nil {
		return err
	}

	return nil
}
