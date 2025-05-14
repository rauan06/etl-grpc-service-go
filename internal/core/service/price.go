package service

import (
	"category/internal/core/domain"
	"category/internal/core/port"
	"category/internal/core/util"
	"context"
	"log/slog"
	"time"
)

type PriceService struct {
	grpcClient port.PriceClient
	httpClient port.PriceClient
	cache      port.CacheRepository

	logger *slog.Logger
	ctx    context.Context
	cancel context.CancelFunc

	status int
}

func NewPriceService(grpcClient port.PriceClient, httpClient port.PriceClient, cache port.CacheRepository, logger *slog.Logger) PriceService {
	ctx, cancel := context.WithCancel(context.Background())

	return PriceService{
		grpcClient: grpcClient,
		httpClient: httpClient,
		cache:      cache,
		logger:     logger,
		ctx:        ctx,
		cancel:     cancel,
		status:     domain.StatusNotStarted,
	}
}

func (s *PriceService) Run() {
	prices := make(chan domain.PriceMain)
	defer close(prices)

	go s.CollectPrices(prices)

	s.status = domain.StatusRunning

	s.logger.Info("price service has started")
	s.SearchPrices(prices)

	s.status = domain.StatusShutdown
}

func (s *PriceService) Status() int {
	return s.status
}

func (s *PriceService) Stop() {
	s.cancel()
	s.status = domain.StatusShutdown
	s.logger.InfoContext(s.ctx, "stopped price service gracefully")
}

func (s *PriceService) SearchPrices(prices chan<- domain.PriceMain) {
	var page int64
	for {
		params := domain.ListParamsSt{
			Page: page,
		}

		resp, err := s.fetchPrices(s.ctx, params)
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
		for _, price := range resp.Results {
			select {
			case prices <- price:
			case <-s.ctx.Done():
				return
			}
		}

		page++
	}
}

func (s *PriceService) fetchPrices(ctx context.Context, params domain.ListParamsSt) (*domain.PriceListRep, error) {
	// Attempt to fetch using gRPC
	resp, err := s.grpcClient.ListPrices(ctx, params, []string{}, []string{})
	if err == nil && len(resp.Results) > 0 {
		return resp, nil
	}

	// If gRPC fails or returns no results, try HTTP
	resp, err = s.httpClient.ListPrices(ctx, params, []string{}, []string{})
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (s *PriceService) CollectPrices(prices <-chan domain.PriceMain) {
	for price := range prices {
		if err := s.cachePrice(price); err != nil {
			s.logger.ErrorContext(s.ctx, "error caching price", "priceID", price.CityId, "error", err.Error())
		}
	}
}

func (s *PriceService) cachePrice(price domain.PriceMain) error {
	data, err := util.Serialize(price)
	if err != nil {
		return err
	}

	key := util.GenerateCacheKey("price", price.CityId)
	if err := s.cache.Set(s.ctx, key, data); err != nil {
		return err
	}

	return nil
}
