package handler

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	pb "github.com/rauan06/etl-grpc-service-go/protos/etl/v1/pb"
	"google.golang.org/protobuf/types/known/timestamppb"

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

	h.status = domain.StatusRunning

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

	h.status = domain.StatusShutdown

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
	res := h.repo.GetResults() // []domain.FullProduct

	var response []*pb.FullProduct
	for _, prod := range res {
		t, err := parseProtobufTimestamp(prod.City.CreatedAt)
		if err != nil {
			return nil, err
		}
		createdAt := timestamppb.New(t)

		t, err = parseProtobufTimestamp(prod.City.UpdatedAt)
		if err != nil {
			return nil, err
		}
		updatedAt := timestamppb.New(t)

		item := &pb.FullProduct{
			Uuid: prod.ID,
			ProductMain: &pb.ProductMain{
				CreatedAt:   prod.ProductMain.CreatedAt, // if string
				UpdatedAt:   prod.ProductMain.UpdatedAt, // if string
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
			},
			City: &pb.CityMain{
				CreatedAt: createdAt,
				UpdatedAt: updatedAt,
				Deleted:   prod.City.Deleted,
				Id:        prod.City.ID,
				Name:      prod.City.Name,
				Postcode:  prod.City.Postcode,
			},
			Price: &pb.PriceMain{
				ProductId: prod.Price.ProductId,
				CityId:    prod.Price.CityId,
				Price:     prod.Price.Price,
			},
			Stock: &pb.StockMain{
				ProductId: prod.Stock.ProductId,
				CityId:    prod.Stock.CityId,
				Value:     prod.Stock.Value,
			},
		}

		response = append(response, item)
	}

	return &pb.FullProductListResponse{
		Results:        response,
		PaginationInfo: &pb.PaginationInfo{Page: 1, PageSize: int64(len(response))},
	}, nil
}

func parseProtobufTimestamp(ts string) (time.Time, error) {
	var seconds, nanos int64
	_, err := fmt.Sscanf(ts, "seconds:%d nanos:%d", &seconds, &nanos)
	if err != nil {
		return time.Time{}, err
	}
	return time.Unix(seconds, nanos).UTC(), nil
}
