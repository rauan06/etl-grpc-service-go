package handler

import (
	pb "category/internal/adapter/handler/grpc/product/v1"
	"category/internal/core/port"
	"context"
	"log/slog"
	"net/url"
	"strconv"
)

type ProductHandler struct {
	pb.UnimplementedProductServer
	svc    port.ProductService
	logger *slog.Logger
}

func NewGrpcProductHandler(svc port.ProductService, logger *slog.Logger) *ProductHandler {
	return &ProductHandler{
		svc:    svc,
		logger: logger,
	}
}

func (h *ProductHandler) List(ctx context.Context, req *pb.ProductListReq) (*pb.ProductListRep, error) {
	h.logger.InfoContext(ctx, "received ListProducts() gRPC request", "req", req)

	params := url.Values{}
	params.Set("list_params.page", strconv.FormatInt(req.GetListParams().GetPage(), 10))
	params.Set("list_params.page_size", strconv.FormatInt(req.GetListParams().GetPageSize(), 10))
	for _, s := range req.GetListParams().GetSort() {
		params.Add("list_params.sort", s)
	}

	if req.GetWithCategory() {
		params.Set("with_category", "true")
	} else {
		params.Set("with_category", "false")
	}

	products, err := h.svc.ListProducts(ctx, params)
	if err != nil {
		h.logger.ErrorContext(ctx, "failed to list products", "error", err)
		return nil, err
	}

	var results []*pb.ProductMain
	for _, p := range products.Results {
		results = append(results, &pb.ProductMain{
			Id:          p.ID,
			Name:        p.Name,
			Description: p.Description,
			Deleted:     p.Deleted,
			CreatedAt:   parseTime(p.CreatedAt),
			UpdatedAt:   parseTime(p.UpdatedAt),
			CategoryId:  p.CategoryID,
			Category:    nil, // Fill if `with_category == true`
		})
	}

	return &pb.ProductListRep{
		PaginationInfo: &pb.PaginationInfoSt{
			Page:     products.PaginationInfo.Page,
			PageSize: products.PaginationInfo.PageSize,
		},
		Results: results,
	}, nil
}

func (h *ProductHandler) Get(ctx context.Context, req *pb.ProductGetReq) (*pb.ProductMain, error) {
	h.logger.InfoContext(ctx, "received GetProduct() gRPC request", "id", req.GetId())

	id, err := strconv.ParseInt(req.GetId(), 10, 64)
	if err != nil {
		h.logger.ErrorContext(ctx, "invalid product ID", "error", err)
		return nil, err
	}

	product, err := h.svc.GetProduct(ctx, id)
	if err != nil {
		h.logger.ErrorContext(ctx, "failed to get product", "error", err)
		return nil, err
	}

	return &pb.ProductMain{
		Id:          product.ID,
		Name:        product.Name,
		Description: product.Description,
		Deleted:     product.Deleted,
		CreatedAt:   parseTime(product.CreatedAt),
		UpdatedAt:   parseTime(product.UpdatedAt),
		CategoryId:  product.CategoryID,
	}, nil
}
