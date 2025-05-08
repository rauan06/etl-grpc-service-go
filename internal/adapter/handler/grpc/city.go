package handler

import (
	pb "category/internal/adapter/handler/grpc/product/v1"
	"category/internal/core/port"
	"context"
	"log/slog"
	"net/url"
	"strconv"
)

type CityHandler struct {
	pb.UnimplementedCityServer
	svc    port.CityService
	logger *slog.Logger
}

func NewGrpcCityHandler(svc port.CityService, logger *slog.Logger) *CityHandler {
	return &CityHandler{
		svc:    svc,
		logger: logger,
	}
}

func (h *CityHandler) List(ctx context.Context, req *pb.CityListReq) (*pb.CityListRep, error) {
	h.logger.InfoContext(ctx, "received ListCities() gRPC request", "req", req)

	params := url.Values{}
	params.Set("list_params.page", strconv.FormatInt(req.GetListParams().GetPage(), 10))
	params.Set("list_params.page_size", strconv.FormatInt(req.GetListParams().GetPageSize(), 10))
	for _, s := range req.GetListParams().GetSort() {
		params.Add("list_params.sort", s)
	}

	if len(req.GetIds()) > 0 {
		for _, id := range req.GetIds() {
			params.Add("list_params.ids", id)
		}
	}

	cities, err := h.svc.ListCities(ctx, params)
	if err != nil {
		h.logger.ErrorContext(ctx, "failed to list cities", "error", err)
		return nil, err
	}

	var results []*pb.CityMain
	for _, c := range cities.Results {
		results = append(results, &pb.CityMain{
			Id:        c.ID,
			Name:      c.Name,
			Deleted:   c.Deleted,
			CreatedAt: parseTime(c.CreatedAt),
			UpdatedAt: parseTime(c.UpdatedAt),
		})
	}

	return &pb.CityListRep{
		PaginationInfo: &pb.PaginationInfoSt{
			Page:     cities.PaginationInfo.Page,
			PageSize: cities.PaginationInfo.PageSize,
		},
		Results: results,
	}, nil
}

func (h *CityHandler) Get(ctx context.Context, req *pb.CityGetReq) (*pb.CityMain, error) {
	h.logger.InfoContext(ctx, "received GetCity() gRPC request", "id", req.GetId())

	id, err := strconv.ParseInt(req.GetId(), 10, 64)
	if err != nil {
		h.logger.ErrorContext(ctx, "invalid ID format", "error", err)
		return nil, err
	}

	city, err := h.svc.GetCity(ctx, id)
	if err != nil {
		h.logger.ErrorContext(ctx, "failed to get city", "error", err)
		return nil, err
	}

	return &pb.CityMain{
		Id:        city.ID,
		Name:      city.Name,
		Deleted:   city.Deleted,
		CreatedAt: parseTime(city.CreatedAt),
		UpdatedAt: parseTime(city.UpdatedAt),
	}, nil
}
