package handler

import (
	"category/internal/core/domain"
	"category/internal/core/service"
	pb "category/protos/etl/v1/pb"
	"context"
)

type EtlHandler struct {
	categorySvc  service.CategoryService
	citySvc      service.CityService
	priceSvc     service.PriceService
	stockSvc     service.StockService
	productSvc   service.ProductService
	collectorSVc service.CollectorService

	status int

	pb.UnimplementedETLServiceServer
}

func NewEtlHandler(categorySvc service.CategoryService, citySvc service.CityService, priceSvc service.PriceService, stockSvc service.StockService, productSvc service.ProductService, collectorSvc service.CollectorService) *EtlHandler {
	return &EtlHandler{
		categorySvc:  categorySvc,
		citySvc:      citySvc,
		priceSvc:     priceSvc,
		stockSvc:     stockSvc,
		productSvc:   productSvc,
		collectorSVc: collectorSvc,
		status:       domain.StatusNotStarted,
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
