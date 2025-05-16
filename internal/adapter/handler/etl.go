package handler

import (
	"context"
	"encoding/json"
	"log/slog"

	pb "github.com/rauan06/etl-grpc-service-go/protos/etl/v1/pb"

	"github.com/rauan06/etl-grpc-service-go/internal/core/domain"
	"github.com/rauan06/etl-grpc-service-go/internal/core/port"
)

type EtlHandler struct {
	svcs []port.Service

	logger *slog.Logger
	cache  port.CacheRepository
	status int

	pb.UnimplementedETLServiceServer
}

func NewEtlHandler(cache port.CacheRepository, logger *slog.Logger, services ...port.Service) *EtlHandler {
	h := &EtlHandler{
		svcs:   []port.Service{},
		status: domain.StatusNotStarted,
		cache:  cache,
		logger: logger,
	}

	h.svcs = append(h.svcs, services...)

	return h
}

func (h *EtlHandler) Start(ctx context.Context, req *pb.ETLRequest) (*pb.ETLResponse, error) {
	for _, svc := range h.svcs {
		svc.Run()
	}

	return &pb.ETLResponse{
		Code:    "200",
		Message: "ETL started",
		Fields:  map[string]string{"status": "started"},
	}, nil
}

func (h *EtlHandler) Stop(ctx context.Context, req *pb.ETLRequest) (*pb.ETLResponse, error) {
	for _, svc := range h.svcs {
		svc.Stop()
	}

	return &pb.ETLResponse{
		Code:    "200",
		Message: "ETL stopped",
		Fields:  map[string]string{"status": "stopped"},
	}, nil
}

func (h *EtlHandler) Status(ctx context.Context, req *pb.ETLRequest) (*pb.ETLResponse, error) {
	resp := &pb.ETLResponse{
		Code:    "200",
		Message: "ETL " + domain.StatusToString(h.status),
		Fields:  map[string]string{},
	}

	for _, svc := range h.svcs {
		resp.Fields[svc.GetServiceName()] = domain.StatusToString(svc.Status())
	}

	return resp, nil
}

func (h *EtlHandler) GetValidProducts(ctx context.Context, req *pb.ETLRequest) (*pb.FullProductListResponse, error) {
	keys, err := h.cache.Scan("stored")
	if err != nil {
		h.logger.ErrorContext(ctx, "error while scanning redis for 'stored'", "error", err.Error())
		return nil, err
	}

	h.logger.InfoContext(ctx, "redis keys found", "count", len(keys), "keys", keys)

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
			continue
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

		response = append(response, &item)
	}

	return &pb.FullProductListResponse{
		Results: response,
	}, nil
}
