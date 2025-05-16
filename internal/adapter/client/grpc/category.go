package grpc

import (
	"context"

	pb "github.com/rauan06/etl-grpc-service-go/protos/product/v1/pb"

	"github.com/rauan06/etl-grpc-service-go/internal/core/domain"
	"google.golang.org/grpc"
)

type CategoryClient struct {
	conn    *grpc.ClientConn
	service pb.CategoryClient
}

func NewCategoryClient(ctx context.Context, url string) (*CategoryClient, error) {
	conn, err := grpc.Dial(url, grpc.WithInsecure())
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
	resp, err := c.service.List(ctx, &pb.CategoryListReq{
		ListParams: &pb.ListParamsSt{
			Page:     params.Page,
			PageSize: params.PageSize,
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
			Page:     resp.PaginationInfo.Page,
			PageSize: resp.PaginationInfo.PageSize,
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
