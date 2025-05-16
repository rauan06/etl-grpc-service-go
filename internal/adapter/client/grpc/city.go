package grpc

import (
	"context"

	"github.com/rauan06/etl-grpc-service-go/internal/core/domain"
	pb "github.com/rauan06/etl-grpc-service-go/protos/product/v1/pb"

	"google.golang.org/grpc"
)

type CityClient struct {
	conn    *grpc.ClientConn
	service pb.CityClient
}

func NewCityClient(ctx context.Context, url string) (*CityClient, error) {
	conn, err := grpc.Dial(url, grpc.WithInsecure())
	if err != nil {
		return nil, err
	}
	return &CityClient{
		conn:    conn,
		service: pb.NewCityClient(conn),
	}, nil
}

func (c *CityClient) Close() {
	c.conn.Close()
}

func (c *CityClient) ListCities(ctx context.Context, params domain.ListParamsSt, ids []string) (*domain.CityListRep, error) {
	resp, err := c.service.List(ctx, &pb.CityListReq{
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

	if resp == nil {
		return nil, nil
	}

	var results []domain.CityMain
	for _, cat := range resp.Results {
		if cat == nil {
			continue
		}
		results = append(results, domain.CityMain{
			UpdatedAt: cat.UpdatedAt.String(),
			CreatedAt: cat.CreatedAt.String(),
			Deleted:   cat.Deleted,
			ID:        cat.Id,
			Name:      cat.Name,
		})
	}

	pagination := domain.PaginationInfoSt{}
	if resp.PaginationInfo != nil {
		pagination.Page = resp.PaginationInfo.Page
		pagination.PageSize = resp.PaginationInfo.PageSize
	}

	return &domain.CityListRep{
		PaginationInfo: pagination,
		Results:        results,
	}, nil
}

func (c *CityClient) GetCity(ctx context.Context, id string) (*domain.CityMain, error) {
	resp, err := c.service.Get(ctx, &pb.CityGetReq{Id: id})
	if err != nil {
		return nil, err
	}

	if resp == nil {
		return nil, nil
	}

	return &domain.CityMain{
		UpdatedAt: resp.UpdatedAt.String(),
		CreatedAt: resp.CreatedAt.String(),
		Deleted:   resp.Deleted,
		ID:        resp.Id,
		Name:      resp.Name,
	}, nil
}
