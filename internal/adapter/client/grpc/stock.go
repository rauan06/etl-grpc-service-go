package grpc

import (
	"context"
	"strconv"

	"category/internal/core/domain"
	pb "category/protos/store/v1/pb"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type StockClient struct {
	conn    *grpc.ClientConn
	service pb.ProductStockClient
}

func NewStoreClient(ctx context.Context, target string) (*StockClient, error) {
	conn, err := grpc.NewClient(target, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}
	return &StockClient{
		conn:    conn,
		service: pb.NewProductStockClient(conn),
	}, nil
}

func (c *StockClient) Close() {
	c.conn.Close()
}

func (c *StockClient) ListStores(ctx context.Context, params domain.ListParamsSt, productIds, cityIDs []string) (*domain.StockListRep, error) {
	page, err := strconv.ParseInt(params.Page, 10, 64)
	if err != nil {
		return nil, domain.ErrParseInt64
	}
	pageSize, err := strconv.ParseInt(params.PageSize, 10, 64)
	if err != nil {
		return nil, domain.ErrParseInt64
	}

	resp, err := c.service.List(ctx, &pb.ProductStockListReq{
		ListParams: &pb.ListParamsSt{
			Page:     page,
			PageSize: pageSize,
			Sort:     params.Sort,
		},
		ProductIds: productIds,
		CityIds:    cityIDs,
	})
	if err != nil {
		return nil, err
	}

	var results []domain.StockMain
	for _, prod := range resp.Results {
		results = append(results, domain.StockMain{
			ProductId: prod.ProductId,
			CityId:    prod.CityId,
			Value:     prod.Value,
		})
	}

	return &domain.StockListRep{
		PaginationInfo: domain.PaginationInfoSt{
			Page:     strconv.FormatInt(resp.PaginationInfo.Page, 10),
			PageSize: strconv.FormatInt(resp.PaginationInfo.PageSize, 10),
		},
		Results: results,
	}, nil
}

func (c *StockClient) GetStore(ctx context.Context, productId, cityId string) (*domain.StockMain, error) {
	resp, err := c.service.Get(ctx, &pb.ProductStockGetReq{
		CityId:    cityId,
		ProductId: productId,
	})
	if err != nil {
		return nil, err
	}

	return &domain.StockMain{
		ProductId: resp.ProductId,
		CityId:    resp.CityId,
		Value:     resp.Value,
	}, nil
}
