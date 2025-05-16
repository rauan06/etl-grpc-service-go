package service

import (
	"context"
	"errors"
	"log/slog"
	"time"

	"github.com/rauan06/etl-grpc-service-go/internal/core/domain"
	"github.com/rauan06/etl-grpc-service-go/internal/core/port"
)

type StockService struct {
	client port.CLientI
	repo   port.Repository

	logger *slog.Logger
	ctx    context.Context
	cancel context.CancelFunc

	status int
}

func NewStockService(client port.CLientI, repo port.Repository, logger *slog.Logger) StockService {
	ctx, cancel := context.WithCancel(context.Background())

	return StockService{
		client: client,
		repo:   repo,
		logger: logger,
		ctx:    ctx,
		cancel: cancel,
		status: domain.StatusNotStarted,
	}
}

func (s *StockService) Run() {
	if s.Status() == domain.StatusRunning {
		return
	}

	s.ctx, s.cancel = context.WithCancel(context.Background())

	stocks := make(chan domain.StockMain)

	go s.CollectStocks(stocks)
	go s.SearchStocks(stocks)

	s.status = domain.StatusRunning

	s.logger.Info("stock service has started")
}

func (s *StockService) GetServiceName() string {
	return "Stock"
}

func (s *StockService) Status() int {
	return s.status
}

func (s *StockService) Stop() {
	if s.status == domain.StatusShutdown {
		return
	}

	s.cancel()

	s.status = domain.StatusShutdown

	s.logger.InfoContext(s.ctx, "stopped stock service gracefully")
}

func (s *StockService) SearchStocks(stocks chan<- domain.StockMain) {
	defer close(stocks)

	var page int64
	for {
		params := domain.ListParamsSt{
			Page: page,
		}

		resp, err := s.fetchStocks(s.ctx, params)
		if err != nil {
			select {
			case <-s.ctx.Done():
				return
			case <-time.After(3 * time.Second):
				page = 0
				continue
			}
		}

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
	resp, err := s.client.ListStocks(ctx, params, []string{}, []string{})
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (s *StockService) CollectStocks(stocks <-chan domain.StockMain) {
	for stock := range stocks {
		if err := s.saveStock(stock); err != nil {
			s.logger.InfoContext(s.ctx, "error caching stock", "stockID", stock.CityId, "error", err.Error())
		}
	}
}

func (s *StockService) saveStock(stock domain.StockMain) error {
	if !stock.IsValid() {
		return errors.New("recieved invalid stock")
	}

	s.repo.SavePair(domain.MarketPair{
		ProductId: stock.ProductId,
		CityId:    stock.CityId,
	})

	return nil
}
