package grpc

import (
	"context"

	pb "github.com/rauan06/etl-grpc-service-go/protos/store/v1/pb"

	"github.com/rauan06/etl-grpc-service-go/internal/core/domain"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type StockClient struct {
	conn    *grpc.ClientConn
	service pb.ProductStockClient
}

func NewStockClient(ctx context.Context, url string) (*StockClient, error) {
	conn, err := grpc.NewClient(url, grpc.WithTransportCredentials(insecure.NewCredentials()))
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

func (c *StockClient) ListStocks(ctx context.Context, params domain.ListParamsSt, productIds, cityIDs []string) (*domain.StockListRep, error) {
	resp, err := c.service.List(ctx, &pb.ProductStockListReq{
		ListParams: &pb.ListParamsSt{
			Page:     params.Page,
			PageSize: params.PageSize,
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
		if prod == nil {
			continue
		}
		results = append(results, domain.StockMain{
			ProductId: prod.ProductId,
			CityId:    prod.CityId,
			Value:     prod.Value,
		})
	}

	pagination := domain.PaginationInfoSt{}
	if resp.PaginationInfo != nil {
		pagination.Page = resp.PaginationInfo.Page
		pagination.PageSize = resp.PaginationInfo.PageSize
	}

	return &domain.StockListRep{
		PaginationInfo: pagination,
		Results:        results,
	}, nil
}

func (c *StockClient) GetStock(ctx context.Context, productId, cityId string) (*domain.StockMain, error) {
	resp, err := c.service.Get(ctx, &pb.ProductStockGetReq{
		CityId:    cityId,
		ProductId: productId,
	})
	if err != nil {
		return nil, err
	}

	if resp == nil {
		return nil, nil
	}

	return &domain.StockMain{
		ProductId: resp.ProductId,
		CityId:    resp.CityId,
		Value:     resp.Value,
	}, nil
}
