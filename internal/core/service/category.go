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

	cache port.CacheRepository
}

func NewCategoryService(grpcClient port.CategoryClient, httpCLient port.CategoryClient, cache port.CacheRepository) *CategoryService {
	return &CategoryService{
		grpcClient,
		httpCLient,

		cache,
	}
}

func (s *CategoryService) Run(ctx context.Context, logger *slog.Logger) error {
	categories := make(chan domain.CategoryMain)

	go s.CollectCategories(ctx, categories)

	ctx, cancel := context.WithCancel(ctx)
	s.SearchCategories(ctx, categories, logger)

	close(categories)
	cancel()

	return nil
}

func (s *CategoryService) Stop(ctx context.Context, logger *slog.Logger) {
	// ctx.Done()

	err := s.cache.DeleteByPrefix(ctx, "category")
	if err != nil {
		logger.ErrorContext(ctx, "error while deleting category keys from redis")
		return
	}

	logger.InfoContext(ctx, "stopped category service gracefully")
}

func (s *CategoryService) SearchCategories(ctx context.Context, categories chan<- domain.CategoryMain, logger *slog.Logger) {
	var page int64 = 0

	for ; ; page++ {
		params := domain.ListParamsSt{
			Page: strconv.FormatInt(page, 10),
		}
		resp, err := s.grpcClient.ListCategories(ctx, params, []string{})
		if err != nil {
			logger.WarnContext(ctx, "gRPC failed, retrying with HTTP", "error", err.Error())
		}

		// No results. Try with http
		if len(resp.Results) == 0 {
			resp, err = s.httpClient.ListCategories(ctx, params, []string{})
			if err != nil {
				logger.ErrorContext(ctx, "both gRPC and HTTP clients failed", "error", err.Error())
				return
			}
		}

		if len(resp.Results) == 0 {
			select {
			case <-ctx.Done():
				return
			case <-time.After(3 * time.Second):
				page = 1
				continue
			}
		}

		// Stream each result
		for _, category := range resp.Results {
			select {
			case categories <- category:
			case <-ctx.Done():
				return
			}
		}

	}
}

func (s *CategoryService) CollectCategories(ctx context.Context, categories <-chan domain.CategoryMain) {
	for category := range categories {
		data, err := util.Serialize(category)
		if err != nil {
			continue
		}
		s.cache.Set(ctx, util.GenerateCacheKey("category", category.ID), data)
		// fmt.Println(category)
	}
}
