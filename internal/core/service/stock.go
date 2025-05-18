package service

import (
	"context"
	"errors"
	"log/slog"
	"sync"
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

	action   bool
	wg       *sync.WaitGroup
	statusMu *sync.Mutex
	status   int
}

func NewStockService(client port.CLientI, repo port.Repository, logger *slog.Logger) StockService {
	ctx, cancel := context.WithCancel(context.Background())

	return StockService{
		client:   client,
		repo:     repo,
		logger:   logger,
		ctx:      ctx,
		cancel:   cancel,
		wg:       new(sync.WaitGroup),
		action:   false,
		statusMu: new(sync.Mutex),
		status:   domain.StatusNotStarted,
	}
}

func (s *StockService) Run() {
	s.statusMu.Lock()
	if s.status == domain.StatusRunning {
		s.statusMu.Unlock()
		return
	}

	s.status = domain.StatusRunning

	s.cancel()
	s.ctx, s.cancel = context.WithCancel(context.Background())
	s.statusMu.Unlock()

	stocks := make(chan domain.StockMain)

	s.wg.Add(2)
	go func() {
		defer s.wg.Done()
		s.CollectStocks(stocks)
	}()
	go func() {
		defer s.wg.Done()
		s.SearchStocks(stocks)
	}()

	s.logger.Info("stock service has started")
}

func (s *StockService) GetServiceName() string {
	return "Stock"
}

func (s *StockService) Status() int {
	s.statusMu.Lock()
	defer s.statusMu.Unlock()

	return s.status
}

func (s *StockService) Stop() {
	s.statusMu.Lock()
	if s.status == domain.StatusShutdown {
		s.statusMu.Unlock()
		return
	}

	s.status = domain.StatusShutdown
	s.statusMu.Unlock()

	s.cancel()
	s.wg.Wait()
	s.logger.InfoContext(s.ctx, "stopped stock service gracefully")
}

func (s *StockService) SearchStocks(stocks chan<- domain.StockMain) {
	defer close(stocks)

	// workers := runtime.GOMAXPROCS(0)

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
	resp, err := s.client.ListStocks(s.ctx, params, []string{}, []string{})
	if err != nil {
		return nil, err
	}

	if len(resp.Results) == 0 {
		return nil, errors.New("empty resutls")
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
