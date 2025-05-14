package service

import (
	"category/internal/core/domain"
	"category/internal/core/port"
	"category/internal/core/util"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"reflect"
	"sync"
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

	var wg sync.WaitGroup
	errChan := make(chan error, 3) // Buffer for scan errors

	for {
		select {
		case <-s.ctx.Done():
			wg.Wait() // Wait for any ongoing operations
			return
		default:
			// Concurrently scan all three key types
			var productKeys, priceKeys, stockKeys []string

			wg.Add(3)
			go func() {
				defer wg.Done()
				keys, err := s.cache.Scan("product")
				if err != nil {
					errChan <- fmt.Errorf("product scan error: %w", err)
					return
				}
				productKeys = keys
			}()

			go func() {
				defer wg.Done()
				keys, err := s.cache.Scan("price")
				if err != nil {
					errChan <- fmt.Errorf("price scan error: %w", err)
					return
				}
				priceKeys = keys
			}()

			go func() {
				defer wg.Done()
				keys, err := s.cache.Scan("stock")
				if err != nil {
					errChan <- fmt.Errorf("stock scan error: %w", err)
					return
				}
				stockKeys = keys
			}()

			wg.Wait()

			// Check for errors
			select {
			case err := <-errChan:
				s.logger.Error("error scanning prefix on redis", "error", err.Error())
				return
			default:
			}

			if len(productKeys) == 0 {
				time.Sleep(3 * time.Second)
				continue
			}

			// Concurrently fetch all data
			var productList []domain.ProductMain
			var priceList []domain.PriceMain
			var stockList []domain.StockMain

			wg.Add(3)
			go func() {
				defer wg.Done()
				s.fetchKeysToSlice(productKeys, &productList)
			}()
			go func() {
				defer wg.Done()
				s.fetchKeysToSlice(priceKeys, &priceList)
			}()
			go func() {
				defer wg.Done()
				s.fetchKeysToSlice(stockKeys, &stockList)
			}()
			wg.Wait()

			// Process data in parallel using worker pool
			type productTask struct {
				product domain.ProductMain
				prices  []domain.PriceMain
				stocks  []domain.StockMain
			}

			taskChan := make(chan productTask, len(productList))
			resultChan := make(chan domain.FullProduct, len(productList))

			// Worker function
			worker := func() {
				defer wg.Done()
				for task := range taskChan {
					fullProd := domain.FullProduct{
						ProductMain: task.product,
						Prices:      task.prices,
						Stocks:      task.stocks,
					}

					if fullProd.IsValid() {
						select {
						case resultChan <- fullProd:
						case <-s.ctx.Done():
							return
						}
					}
				}
			}

			// Create worker pool
			numWorkers := min(10, len(productList))
			wg.Add(numWorkers)
			for i := 0; i < numWorkers; i++ {
				go worker()
			}

			// Build lookup maps
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

			// Distribute tasks
			go func() {
				defer close(taskChan)
				for _, prod := range productList {
					task := productTask{
						product: prod,
						prices:  productIDToPrices[prod.ID],
						stocks:  productIDToStocks[prod.ID],
					}
					select {
					case taskChan <- task:
					case <-s.ctx.Done():
						return
					}
				}
			}()

			// Collect results
			go func() {
				wg.Wait()
				close(resultChan)
			}()

			// Send results to output channel
			for fullProd := range resultChan {
				select {
				case products <- fullProd:
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

		// s.logger.InfoContext(s.ctx, "stored full product")
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
