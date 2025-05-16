package handler

import (
	"context"
	"log/slog"

	pb "github.com/rauan06/etl-grpc-service-go/protos/etl/v1/pb"

	"github.com/rauan06/etl-grpc-service-go/internal/core/domain"
	"github.com/rauan06/etl-grpc-service-go/internal/core/port"
)

type EtlHandler struct {
	svcs []port.Service
	repo port.Repository

	logger *slog.Logger
	status int

	pb.UnimplementedETLServiceServer
}

func NewEtlHandler(repo port.Repository, logger *slog.Logger, services ...port.Service) *EtlHandler {
	h := &EtlHandler{
		svcs:   []port.Service{},
		repo:   repo,
		status: domain.StatusNotStarted,
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
	res := h.repo.GetResults()

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

		response = append(response, &item)
	}

	return &pb.FullProductListResponse{
		Results: response,
	}, nil
}
