package grpc

import (
	"context"
	"strconv"

	"category/internal/adapter/handler"
	"category/internal/core/domain"
	pb "category/protos/product/v1/pb"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type CityClient struct {
	conn    *grpc.ClientConn
	service pb.CityClient
}

func NewCityClient(ctx context.Context, url string) (*CityClient, error) {
	conn, err := grpc.NewClient(url, grpc.WithTransportCredentials(insecure.NewCredentials()))
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
	page, err := strconv.ParseInt(params.Page, 10, 64)
	if err != nil {
		return nil, handler.ErrParseInt64
	}
	pageSize, err := strconv.ParseInt(params.PageSize, 10, 64)
	if err != nil {
		return nil, handler.ErrParseInt64
	}

	resp, err := c.service.List(ctx, &pb.CityListReq{
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

	var results []domain.CityMain
	for _, cat := range resp.Results {
		results = append(results, domain.CityMain{
			UpdatedAt: cat.UpdatedAt.String(),
			CreatedAt: cat.CreatedAt.String(),
			Deleted:   cat.Deleted,
			ID:        cat.Id,
			Name:      cat.Name,
		})
	}

	return &domain.CityListRep{
		PaginationInfo: domain.PaginationInfoSt{
			Page:     strconv.FormatInt(resp.PaginationInfo.Page, 10),
			PageSize: strconv.FormatInt(resp.PaginationInfo.PageSize, 10),
		},
		Results: results,
	}, nil
}

func (c *CityClient) GetCity(ctx context.Context, id string) (*domain.CityMain, error) {
	resp, err := c.service.Get(ctx, &pb.CityGetReq{Id: id})
	if err != nil {
		return nil, err
	}

	return &domain.CityMain{
		UpdatedAt: resp.UpdatedAt.String(),
		CreatedAt: resp.CreatedAt.String(),
		Deleted:   resp.Deleted,
		ID:        resp.Id,
		Name:      resp.Name,
	}, nil
}
