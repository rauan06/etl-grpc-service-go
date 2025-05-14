package handler

import (
	"category/internal/core/domain"
	"category/internal/core/port"
	"category/internal/core/service"
	pb "category/protos/etl/v1/pb"
	"context"
	"encoding/json"
	"log/slog"
)

type EtlHandler struct {
	categorySvc  service.CategoryService
	citySvc      service.CityService
	priceSvc     service.PriceService
	stockSvc     service.StockService
	productSvc   service.ProductService
	collectorSVc service.CollectorService

	logger *slog.Logger
	cache  port.CacheRepository
	status int

	pb.UnimplementedETLServiceServer
}

func NewEtlHandler(categorySvc service.CategoryService, citySvc service.CityService, priceSvc service.PriceService, stockSvc service.StockService, productSvc service.ProductService, collectorSvc service.CollectorService, cache port.CacheRepository, logger *slog.Logger) *EtlHandler {
	return &EtlHandler{
		categorySvc:  categorySvc,
		citySvc:      citySvc,
		priceSvc:     priceSvc,
		stockSvc:     stockSvc,
		productSvc:   productSvc,
		collectorSVc: collectorSvc,
		status:       domain.StatusNotStarted,
		cache:        cache,
		logger:       logger,
	}
}

func (h *EtlHandler) Start(ctx context.Context, req *pb.ETLRequest) (*pb.ETLResponse, error) {
	h.categorySvc.Run()
	h.citySvc.Run()
	h.priceSvc.Run()
	h.stockSvc.Run()
	h.productSvc.Run()
	h.collectorSVc.Run()

	return &pb.ETLResponse{
		Code:    "200",
		Message: "ETL started",
		Fields:  map[string]string{"status": "started"},
	}, nil
}

func (h *EtlHandler) Stop(ctx context.Context, req *pb.ETLRequest) (*pb.ETLResponse, error) {
	h.categorySvc.Stop()
	h.citySvc.Stop()
	h.priceSvc.Stop()
	h.stockSvc.Stop()
	h.productSvc.Stop()
	h.collectorSVc.Stop()

	return &pb.ETLResponse{
		Code:    "200",
		Message: "ETL stopped",
		Fields:  map[string]string{"status": "stopped"},
	}, nil
}

func (h *EtlHandler) Status(ctx context.Context, req *pb.ETLRequest) (*pb.ETLResponse, error) {
	return &pb.ETLResponse{
		Code:    "200",
		Message: "ETL " + domain.StatusToString(h.status),
		Fields: map[string]string{
			"category_status":  domain.StatusToString(h.categorySvc.Status()),
			"city_status":      domain.StatusToString(h.citySvc.Status()),
			"price_status":     domain.StatusToString(h.priceSvc.Status()),
			"stock_status":     domain.StatusToString(h.stockSvc.Status()),
			"product_status":   domain.StatusToString(h.productSvc.Status()),
			"collector_status": domain.StatusToString(h.collectorSVc.Status()),
		},
	}, nil
}

func (h *EtlHandler) GetValidProducts(ctx context.Context, req *pb.ETLRequest) (*pb.FullProductListResponse, error) {
	keys, err := h.cache.Scan("stored")
	if err != nil {
		h.logger.ErrorContext(ctx, "error while scanning redis for 'stored'", "error", err.Error())
		return nil, err
	}

	// h.logger.InfoContext(ctx, "redis keys found", "count", len(keys))

	var res = []*domain.FullProduct{}
	for _, key := range keys {
		// Get the data for each key from Redis
		data, err := h.cache.Get(ctx, key)
		if err != nil {
			h.logger.ErrorContext(ctx, "error while getting key from redis", "error", err.Error())
			return nil, err
		}

		// Unmarshal the data into a FullProduct
		var item *domain.FullProduct
		err = json.Unmarshal(data, &item) // Pass a pointer to the item
		if err != nil {
			h.logger.ErrorContext(ctx, "error unmarshalling data", "error", err.Error())
			continue // Skip to the next key if unmarshalling fails
		}

		res = append(res, item)
	}

	var response []*pb.FullProduct
	for _, prod := range res {
		var item pb.FullProduct

		// Map product data to pb.FullProduct
		item.ProductMain = &pb.ProductMain{
			CreatedAt:   prod.ProductMain.CreatedAt,
			UpdatedAt:   prod.ProductMain.UpdatedAt,
			Deleted:     prod.ProductMain.Deleted,
			Id:          prod.ProductMain.ID,
			Name:        prod.ProductMain.Name,
			Description: prod.ProductMain.Description,
			CategoryId:  prod.ProductMain.CategoryID,
			Category: &pb.CategoryMain{
				CreatedAt: prod.ProductMain.Category.CreatedAt,
				UpdatedAt: prod.ProductMain.Category.UpdatedAt,
				Deleted:   prod.ProductMain.Category.Deleted,
				Id:        prod.ProductMain.Category.ID,
				Name:      prod.ProductMain.Category.Name,
			},
		}

		// Map prices to pb.PriceMain
		for _, price := range prod.Prices {
			item.Price = append(item.Price, &pb.PriceMain{
				ProductId: price.ProductId,
				CityId:    price.CityId,
				Price:     price.Price,
			})
		}

		// Map stock to pb.StockMain
		for _, stock := range prod.Stocks {
			item.Stock = append(item.Stock, &pb.StockMain{
				ProductId: stock.ProductId,
				CityId:    stock.CityId,
				Value:     stock.Value,
			})
		}

		// Add the populated FullProduct to the response slice
		response = append(response, &item)
	}

	// Return the FullProductListResponse with pagination information
	return &pb.FullProductListResponse{
		Results: response,
	}, nil
}
