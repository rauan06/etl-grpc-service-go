package service

import (
	"context"
	"log/slog"
	"strconv"
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
}

func NewCategoryService(grpcClient port.CategoryClient, httpClient port.CategoryClient, cache port.CacheRepository, logger *slog.Logger) CategoryService {
	return CategoryService{
		grpcClient: grpcClient,
		httpClient: httpClient,
		cache:      cache,
		logger:     logger,
		ctx:        context.Background(),
	}
}

func (s *CategoryService) Run(ctx context.Context) {
	s.ctx = ctx

	categories := make(chan domain.CategoryMain)
	defer close(categories)

	go s.CollectCategories(categories)

	s.SearchCategories(categories)
}

func (s *CategoryService) Status() int {
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

func (s *CategoryService) Stop() {
	if s.ctx == nil {
		return
	}

	s.ctx.Done()
	s.logger.InfoContext(s.ctx, "stopped category service gracefully")
}

func (s *CategoryService) SearchCategories(categories chan<- domain.CategoryMain) {
	var page int64
	for {
		params := domain.ListParamsSt{
			Page: strconv.FormatInt(page, 10),
		}

		resp, err := s.fetchCategories(s.ctx, params)
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

	// If gRPC fails or returns no results, try HTTP
	s.logger.WarnContext(ctx, "gRPC failed, retrying with HTTP", "error", err.Error())
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
