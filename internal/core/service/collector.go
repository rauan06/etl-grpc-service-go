package service

import (
	"category/internal/core/domain"
	"category/internal/core/port"
	"category/internal/core/util"
	"context"
	"encoding/json"
	"errors"
	"log/slog"
	"reflect"
	"time"
)

type CollectorService struct {
	cache  port.CacheRepository
	logger *slog.Logger
	ctx    context.Context
	cancel context.CancelFunc

	status int
}

func NewCollectorService(cache port.CacheRepository, logger *slog.Logger) CollectorService {
	ctx, cancel := context.WithCancel(context.Background())

	return CollectorService{
		cache:  cache,
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

	s.status = domain.StatusRunning

	products := make(chan domain.FullProduct)

	go s.storeProductDetails(products)
	go s.collectProductDetails(products)
	s.logger.Info("collector service has started")
}

func (s *CollectorService) Status() int {
	return s.status
}

func (s *CollectorService) Stop() {
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
			productKeys, err := s.cache.Scan("product")
			if err != nil {
				s.logger.Error("error scanning prefix on redis", "error", err.Error())
				return
			}

			priceKeys, err := s.cache.Scan("price")
			if err != nil {
				s.logger.Error("error scanning prefix on redis", "error", err.Error())
				return
			}

			stockKeys, err := s.cache.Scan("stock")
			if err != nil {
				s.logger.Error("error scanning prefix on redis", "error", err.Error())
				return
			}

			var productList []domain.ProductMain
			var priceList []domain.PriceMain
			var stockList []domain.StockMain
			s.fetchKeysToSlice(productKeys, &productList)
			s.fetchKeysToSlice(priceKeys, &priceList)
			s.fetchKeysToSlice(stockKeys, &stockList)

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
	}
}

func (s *CollectorService) getCache(serviceName, id string) ([]byte, error) {
	key := util.GenerateCacheKey(serviceName, id)
	return s.cache.Get(s.ctx, key)
}

func (s *CollectorService) fetchKeysToSlice(keys []string, dst interface{}) error {
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
