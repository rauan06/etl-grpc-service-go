package service

import (
	"category/internal/core/domain"
	"category/internal/core/port"
	"category/internal/core/util"
	"context"
	"log/slog"
	"strconv"
	"time"
)

type CityService struct {
	grpcClient port.CityClient
	httpClient port.CityClient
	cache      port.CacheRepository

	logger *slog.Logger
	ctx    context.Context
}

func NewCityService(grpcClient port.CityClient, httpClient port.CityClient, cache port.CacheRepository, logger *slog.Logger) CityService {
	return CityService{
		grpcClient: grpcClient,
		httpClient: httpClient,
		cache:      cache,
		logger:     logger,
		ctx:        context.Background(),
	}
}

func (s *CityService) Run(ctx context.Context) {
	s.ctx = ctx

	cities := make(chan domain.CityMain)
	defer close(cities)

	go s.CollectCities(cities)

	s.SearchCities(cities)
}

func (s *CityService) Status() int {
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

func (s *CityService) Stop() {
	if s.ctx == nil {
		return
	}

	s.ctx.Done()
	s.logger.InfoContext(s.ctx, "stopped city service gracefully")
}

func (s *CityService) SearchCities(cities chan<- domain.CityMain) {
	var page int64
	for {
		params := domain.ListParamsSt{
			Page: strconv.FormatInt(page, 10),
		}

		resp, err := s.fetchCities(s.ctx, params)
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

	// If gRPC fails or returns no results, try HTTP
	s.logger.WarnContext(ctx, "gRPC failed, retrying with HTTP", "error", err.Error())
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
