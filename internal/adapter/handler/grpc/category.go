package handler

import (
	pb "category/internal/adapter/handler/grpc/product/v1"
	"category/internal/core/port"
	"context"
	"log/slog"
	"net/url"
	"strconv"
)

type CategoryGrpcHandler struct {
	pb.UnimplementedCategoryServer
	svc    port.CategoryService
	logger *slog.Logger
}

func NewCategoryGrpcHandler(svc port.CategoryService, logger *slog.Logger) *CategoryGrpcHandler {
	return &CategoryGrpcHandler{
		svc:    svc,
		logger: logger,
	}
}

func (h *CategoryGrpcHandler) List(ctx context.Context, req *pb.CategoryListReq) (*pb.CategoryListRep, error) {
	h.logger.InfoContext(ctx, "received ListCategories() gRPC request", "req", req)

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

	categories, err := h.svc.ListCategories(ctx, params)
	if err != nil {
		h.logger.ErrorContext(ctx, "failed to list categories", "error", err)
		return nil, err
	}

	var results []*pb.CategoryMain
	for _, c := range categories.Results {
		results = append(results, &pb.CategoryMain{
			Id:        c.ID,
			Name:      c.Name,
			Deleted:   c.Deleted,
			CreatedAt: parseTime(c.CreatedAt),
			UpdatedAt: parseTime(c.UpdatedAt),
		})
	}

	return &pb.CategoryListRep{
		PaginationInfo: &pb.PaginationInfoSt{
			Page:     categories.PaginationInfo.Page,
			PageSize: categories.PaginationInfo.PageSize,
		},
		Results: results,
	}, nil
}

func (h *CategoryGrpcHandler) Get(ctx context.Context, req *pb.CategoryGetReq) (*pb.CategoryMain, error) {
	h.logger.InfoContext(ctx, "received GetCategory() gRPC request", "id", req.GetId())

	id, err := strconv.ParseInt(req.GetId(), 10, 64)
	if err != nil {
		h.logger.ErrorContext(ctx, "invalid ID format", "error", err)
		return nil, err
	}

	category, err := h.svc.GetCategory(ctx, id)
	if err != nil {
		h.logger.ErrorContext(ctx, "failed to get category", "error", err)
		return nil, err
	}

	return &pb.CategoryMain{
		Id:        category.ID,
		Name:      category.Name,
		Deleted:   category.Deleted,
		CreatedAt: parseTime(category.CreatedAt),
		UpdatedAt: parseTime(category.UpdatedAt),
	}, nil
}
