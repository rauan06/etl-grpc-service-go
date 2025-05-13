package service

import (
	"category/internal/core/domain"
	"category/internal/core/port"
	"category/internal/core/util"
	"context"
	"log/slog"
	"time"
)

type StockService struct {
	grpcClient port.StockClient
	httpClient port.StockClient
	cache      port.CacheRepository

	logger *slog.Logger
	ctx    context.Context
}

func NewStockService(grpcClient port.StockClient, httpClient port.StockClient, cache port.CacheRepository, logger *slog.Logger) StockService {
	return StockService{
		grpcClient: grpcClient,
		httpClient: httpClient,
		cache:      cache,
		logger:     logger,
		ctx:        context.Background(),
	}
}

func (s *StockService) Run() {
	stocks := make(chan domain.StockMain)
	defer close(stocks)

	go s.CollectStocks(stocks)

	s.SearchStocks(stocks)
}

func (s *StockService) Status() int {
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

func (s *StockService) Stop() {
	if s.ctx == nil {
		return
	}

	s.ctx.Done()
	s.logger.InfoContext(s.ctx, "stopped stock service gracefully")
}

func (s *StockService) SearchStocks(stocks chan<- domain.StockMain) {
	var page int64
	for {
		params := domain.ListParamsSt{
			Page: page,
		}

		resp, err := s.fetchStocks(s.ctx, params)
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
		for _, stock := range resp.Results {
			select {
			case stocks <- stock:
			case <-s.ctx.Done():
				return
			}
		}

		page++
	}
}

func (s *StockService) fetchStocks(ctx context.Context, params domain.ListParamsSt) (*domain.StockListRep, error) {
	resp, err := s.grpcClient.ListStocks(ctx, params, []string{}, []string{})
	if err == nil && len(resp.Results) > 0 {
		return resp, nil
	}

	s.logger.WarnContext(ctx, "gRPC failed, retrying with HTTP", "error", err.Error())
	resp, err = s.httpClient.ListStocks(ctx, params, []string{}, []string{})
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (s *StockService) CollectStocks(stocks <-chan domain.StockMain) {
	for stock := range stocks {
		if err := s.cacheStock(stock); err != nil {
			s.logger.ErrorContext(s.ctx, "error caching stock", "stockID", stock.CityId, "error", err.Error())
		}
	}
}

func (s *StockService) cacheStock(stock domain.StockMain) error {
	data, err := util.Serialize(stock)
	if err != nil {
		return err
	}

	key := util.GenerateCacheKey("stock", stock.CityId)
	if err := s.cache.Set(s.ctx, key, data); err != nil {
		return err
	}

	return nil
}
