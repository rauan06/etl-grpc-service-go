package service

import (
	"context"
	"encoding/json"
	"errors"
	"log/slog"
	"reflect"
	"time"

	"github.com/rauan06/etl-grpc-service-go/internal/core/domain"
	"github.com/rauan06/etl-grpc-service-go/internal/core/port"
	"github.com/rauan06/etl-grpc-service-go/internal/core/util"
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

			// Fetch stock and price first
			stock := s.transformPairToStock(marketPair)
			price := s.transformPairToPrice(marketPair)

			if len(productList) == 0 {
				time.Sleep(3 * time.Second)
				continue
			}

			// Group prices and stocks by product ID
			productIDToPrices := make(map[string][]domain.PriceMain)
			for _, price := range priceList {
				if price.IsValid() {
					productIDToPrices[price.ProductId] = append(productIDToPrices[price.ProductId], price)
				}
			}

			productIDToStocks := make(map[string][]domain.StockMain)
			for _, stock := range stockList {
				if stock.IsValid() {
					productIDToStocks[stock.ProductId] = append(productIDToStocks[stock.ProductId], stock)
				}
			}

			var fullProds []domain.FullProduct
			for _, prod := range productList {
				fullProd := domain.FullProduct{
					ProductMain: prod,
					Prices:      productIDToPrices[prod.ID],
					Stocks:      productIDToStocks[prod.ID],
				}

				if !fullProd.IsValid() {
					continue
				}

				fullProds = append(fullProds, fullProd)
			}

			for _, prod := range fullProds {
				select {
				case products <- prod:
				case <-s.ctx.Done():
					return
				}
			}
		}
	}
}

func (s *CollectorService) storeProductDetails(products <-chan domain.FullProduct) {
	for product := range products {
		data, err := util.Serialize(product)
		if err != nil {
			s.logger.ErrorContext(s.ctx, "error while serializing valid full product", "error", err.Error())
			continue
		}

		err = s.cache.Set(s.ctx, util.GenerateCacheKey("stored", product.ProductMain.ID), data)
		if err != nil {
			s.logger.ErrorContext(s.ctx, "error while setting to redis", "error", err)
		}

		// s.logger.InfoContext(s.ctx, "stored full product", "product", product, "key", "stored:"+product.ProductMain.ID)
	}
}

func (s *CollectorService) getCache(serviceName, id string) ([]byte, error) {
	key := util.GenerateCacheKey(serviceName, id)
	return s.cache.Get(s.ctx, key)
}

func (s *CollectorService) transformPairToPrice(pair domain.MarketPair) domain.PriceMain {
	dstVal := reflect.ValueOf(dst)
	if dstVal.Kind() != reflect.Ptr || dstVal.Elem().Kind() != reflect.Slice {
		return errors.New("dst must be a pointer to a slice")
	}

	sliceVal := dstVal.Elem()
	elemType := sliceVal.Type().Elem() // type of individual item

	for _, key := range keys {
		data, err := s.cache.Get(s.ctx, key)
		if err != nil {
			continue
		}

		itemPtr := reflect.New(elemType) // create pointer to T
		if err := json.Unmarshal(data, itemPtr.Interface()); err != nil {
			continue
		}

		sliceVal.Set(reflect.Append(sliceVal, itemPtr.Elem()))
	}

	return nil
}

func (s *CollectorService) transformPairToStock(pair domain.MarketPair) domain.StockMain {
}
