package service

import (
	"context"
	"log/slog"
	"time"

	"github.com/google/uuid"
	"github.com/rauan06/etl-grpc-service-go/internal/core/domain"
	"github.com/rauan06/etl-grpc-service-go/internal/core/port"
)

type CollectorService struct {
	client port.CLientI
	repo   port.Repository

	logger *slog.Logger
	ctx    context.Context
	cancel context.CancelFunc

	status int
}

func NewCollectorService(client port.CLientI, repo port.Repository, logger *slog.Logger) CollectorService {
	ctx, cancel := context.WithCancel(context.Background())

	return CollectorService{
		client: client,
		repo:   repo,
		logger: logger,
		ctx:    ctx,
		cancel: cancel,
		status: domain.StatusNotStarted,
	}
}

func (s *CollectorService) Run() {
	if s.Status() == domain.StatusRunning {
		return
	}

	s.ctx, s.cancel = context.WithCancel(context.Background())

	s.status = domain.StatusRunning

	products := make(chan domain.FullProduct)

	go s.storeProductDetails(products)
	go s.collectProductDetails(products)
	s.logger.Info("collector service has started")
}

func (s *CollectorService) GetServiceName() string {
	return "Collector"
}

func (s *CollectorService) Status() int {
	return s.status
}

func (s *CollectorService) Stop() {
	if s.status == domain.StatusShutdown {
		return
	}

	s.cancel()

	s.status = domain.StatusShutdown

	s.logger.InfoContext(s.ctx, "stopped collector service gracefully")
}

func (s *CollectorService) collectProductDetails(products chan<- domain.FullProduct) {
	defer close(products)

	for {
		select {
		case <-s.ctx.Done():
			return
		default:
			marketPair, ok := s.repo.ReadPair()
			if !ok {
				// Retry if repo is empty or closed
				s.logger.Info("didn't recieve pir from storage, retrying...")
				time.Sleep(time.Second)
				continue
			}

			// Fetch models using our piared ids
			stock, err := s.transformPairToStock(marketPair)
			if err != nil {
				s.logger.Warn("failed to get stock", "error", err)

				// Continue means skipping this pair of ids
				continue
			}

			price, err := s.transformPairToPrice(marketPair)
			if err != nil {
				s.logger.Warn("failed to get price", "error", err)

				continue
			}

			city, err := s.transformPairToCity(marketPair)
			if err != nil {
				s.logger.Warn("failed to get city", "error", err)

				continue
			}

			product, err := s.transformPairToProduct(marketPair)
			if err != nil {
				s.logger.Warn("failed to get product", "error", err)

				continue
			}

			category, err := s.transformIdToCategory(product.CategoryID)
			if err != nil {
				s.logger.Warn("failed to get product", "error", err)

				continue
			}

			product.Category = *category

			// Group models and create new uuid
			fullProduct := domain.FullProduct{
				ID:          uuid.NewString(),
				ProductMain: *product,
				City:        *city,
				Price:       *price,
				Stock:       *stock,
			}

			if !fullProduct.IsValid() {
				s.logger.Warn("recieved invalid full product", "product", fullProduct)
				continue
			}

			products <- fullProduct
		}
	}
}

func (s *CollectorService) storeProductDetails(products <-chan domain.FullProduct) {
	for product := range products {
		s.repo.SaveResult(product.ID, product)
	}
}

func (s *CollectorService) transformPairToPrice(pair domain.MarketPair) (*domain.PriceMain, error) {
	return s.client.GetPrice(context.Background(), pair.ProductId, pair.CityId)
}

func (s *CollectorService) transformPairToStock(pair domain.MarketPair) (*domain.StockMain, error) {
	return s.client.GetStock(context.Background(), pair.ProductId, pair.CityId)
}

func (s *CollectorService) transformPairToProduct(pair domain.MarketPair) (*domain.ProductMain, error) {
	return s.client.GetProduct(context.Background(), pair.ProductId)
}

func (s *CollectorService) transformPairToCity(pair domain.MarketPair) (*domain.CityMain, error) {
	return s.client.GetCity(context.Background(), pair.CityId)
}

func (s *CollectorService) transformIdToCategory(categoryId string) (*domain.CategoryMain, error) {
	return s.client.GetCategory(context.Background(), categoryId)
}
