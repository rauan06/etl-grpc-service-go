package grpc

import (
	"context"
	"strconv"

	"category/internal/core/domain"
	pb "category/protos/product/v1/pb"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type CategoryClient struct {
	conn    *grpc.ClientConn
	service pb.CategoryClient
}

func NewCategoryClient(ctx context.Context, url string) (*CategoryClient, error) {
	conn, err := grpc.NewClient(url, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}
	return &CategoryClient{
		conn:    conn,
		service: pb.NewCategoryClient(conn),
	}, nil
}

func (c *CategoryClient) Close() {
	c.conn.Close()
}

func (c *CategoryClient) ListCategories(ctx context.Context, params domain.ListParamsSt, ids []string) (*domain.CategoryListRep, error) {
	page, err := strconv.ParseInt(params.Page, 10, 64)
	if err != nil {
		return nil, domain.ErrParseInt64
	}
	pageSize, err := strconv.ParseInt(params.PageSize, 10, 64)
	if err != nil {
		return nil, domain.ErrParseInt64
	}

	resp, err := c.service.List(ctx, &pb.CategoryListReq{
		ListParams: &pb.ListParamsSt{
			Page:     page,
			PageSize: pageSize,
			Sort:     params.Sort,
		},
		Ids: ids,
	})
	if err != nil {
		return nil, err
	}

	var results []domain.CategoryMain
	for _, cat := range resp.Results {
		results = append(results, domain.CategoryMain{
			UpdatedAt: cat.UpdatedAt.String(),
			CreatedAt: cat.CreatedAt.String(),
			Deleted:   cat.Deleted,
			ID:        cat.Id,
			Name:      cat.Name,
		})
	}

	return &domain.CategoryListRep{
		PaginationInfo: domain.PaginationInfoSt{
			Page:     strconv.FormatInt(resp.PaginationInfo.Page, 10),
			PageSize: strconv.FormatInt(resp.PaginationInfo.PageSize, 10),
		},
		Results: results,
	}, nil
}

func (c *CategoryClient) GetCategory(ctx context.Context, id string) (*domain.CategoryMain, error) {
	resp, err := c.service.Get(ctx, &pb.CategoryGetReq{Id: id})
	if err != nil {
		return nil, err
	}

	return &domain.CategoryMain{
		UpdatedAt: resp.UpdatedAt.String(),
		CreatedAt: resp.CreatedAt.String(),
		Deleted:   resp.Deleted,
		ID:        resp.Id,
		Name:      resp.Name,
	}, nil
}
