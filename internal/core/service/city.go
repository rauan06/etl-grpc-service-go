package service

import (
	"category/internal/core/domain"
	"category/internal/core/port"
	"category/internal/core/util"
	"context"
	"log/slog"
	"time"
)

type CityService struct {
	grpcClient port.CityClient
	httpClient port.CityClient
	cache      port.CacheRepository

	logger *slog.Logger
	ctx    context.Context
	cancel context.CancelFunc

	status int
}

func NewCityService(grpcClient port.CityClient, httpClient port.CityClient, cache port.CacheRepository, logger *slog.Logger) CityService {
	ctx, cancel := context.WithCancel(context.Background())
	return CityService{
		grpcClient: grpcClient,
		httpClient: httpClient,
		cache:      cache,
		logger:     logger,
		ctx:        ctx,
		cancel:     cancel,
		status:     domain.StatusNotStarted,
	}
}

func (s *CityService) Run() {
	if s.Status() == domain.StatusRunning {
		return
	}

	cities := make(chan domain.CityMain)

	go s.CollectCities(cities)

	s.status = domain.StatusRunning

	s.logger.Info("city service has started")
	go s.SearchCities(cities)
}

func (s *CityService) Status() int {
	return s.status
}

func (s *CityService) Stop() {
	s.cancel()
	s.status = domain.StatusShutdown
	s.logger.InfoContext(s.ctx, "stopped city service gracefully")
}

func (s *CityService) SearchCities(cities chan<- domain.CityMain) {
	defer close(cities)

	var page int64
	for {
		params := domain.ListParamsSt{
			Page: page,
		}

		resp, err := s.fetchCities(s.ctx, params)
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
		for _, city := range resp.Results {
			select {
			case cities <- city:
			case <-s.ctx.Done():
				return
			}
		}

		page++
	}
}

func (s *CityService) fetchCities(ctx context.Context, params domain.ListParamsSt) (*domain.CityListRep, error) {
	// Attempt to fetch using gRPC
	resp, err := s.grpcClient.ListCities(ctx, params, []string{})
	if err == nil && len(resp.Results) > 0 {
		return resp, nil
	}

	resp, err = s.httpClient.ListCities(ctx, params, []string{})
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (s *CityService) CollectCities(cities <-chan domain.CityMain) {
	for city := range cities {
		if err := s.cacheCity(city); err != nil {
			s.logger.ErrorContext(s.ctx, "error caching city", "cityID", city.ID, "error", err.Error())
		}
	}
}

func (s *CityService) cacheCity(city domain.CityMain) error {
	data, err := util.Serialize(city)
	if err != nil {
		return err
	}

	key := util.GenerateCacheKey("city", city.ID)
	if err := s.cache.Set(s.ctx, key, data); err != nil {
		return err
	}

	return nil
}
