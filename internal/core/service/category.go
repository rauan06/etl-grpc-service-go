package service

import (
	"context"
	"fmt"
	"log/slog"
	"strconv"
	"time"

	"category/internal/core/domain"
	"category/internal/core/port"
)

type CategoryService struct {
	grpcClient port.CategoryClient
	httpClient port.CategoryClient

	// cache port.CacheRepository
}

func NewCategoryService(grpcClient port.CategoryClient, httpCLient port.CategoryClient) *CategoryService {
	return &CategoryService{
		grpcClient,
		httpCLient,

		// cache,
	}
}

func (s *CategoryService) Run(ctx context.Context, logger *slog.Logger) error {
	categories := make(chan domain.CategoryMain)

	go s.CollectCategories(ctx, categories)

	s.SearchCategories(ctx, categories, logger)
	close(categories)

	return nil
}

func (s *CategoryService) Stop(ctx context.Context, logger *slog.Logger) {
	// ctx.Done()

	// err := s.cache.DeleteByPrefix(ctx, "category")
	// if err != nil {
	// 	logger.ErrorContext(ctx, "error while deleting category keys from redis")
	// 	return
	// }

	logger.InfoContext(ctx, "stopped category service gracefully")
}

func (s *CategoryService) SearchCategories(ctx context.Context, categories chan<- domain.CategoryMain, logger *slog.Logger) {
	client := s.grpcClient
	var page int64 = 1

	for {
		params := domain.ListParamsSt{
			Page: strconv.FormatInt(page, 10),
		}
		page++

		resp, err := client.ListCategories(ctx, params, []string{})
		if err != nil {
			logger.WarnContext(ctx, "gRPC failed, retrying with HTTP", "error", err.Error())
			client = s.httpClient
			resp, err = client.ListCategories(ctx, params, []string{})
			if err != nil {
				logger.ErrorContext(ctx, "both gRPC and HTTP clients failed", "error", err.Error())
				return
			}
		}

		// No results. Wait and try again
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
		// data, err := util.Serialize(category)
		// if err != nil {
		// 	continue
		// }
		// s.cache.Set(ctx, util.GenerateCacheKey("category", category.ID), data)
		fmt.Println(category)
	}
}
