package grpc

import (
	"context"

	"category/internal/core/domain"
	pb "category/protos/price/v1/pb"

	"google.golang.org/grpc"
)

type PriceClient struct {
	conn    *grpc.ClientConn
	service pb.ProductPriceClient
}

func NewPriceClient(ctx context.Context, url string) (*PriceClient, error) {
	conn, err := grpc.Dial(url, grpc.WithInsecure())
	if err != nil {
		return nil, err
	}
	return &PriceClient{
		conn:    conn,
		service: pb.NewProductPriceClient(conn),
	}, nil
}

func (c *PriceClient) Close() {
	c.conn.Close()
}

func (c *PriceClient) ListPrices(ctx context.Context, params domain.ListParamsSt, productIds, cityIDs []string) (*domain.PriceListRep, error) {
	resp, err := c.service.List(ctx, &pb.ProductPriceListReq{
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

	var results []domain.PriceMain
	for _, prod := range resp.Results {
		results = append(results, domain.PriceMain{
			ProductId: prod.ProductId,
			CityId:    prod.CityId,
			Price:     prod.Price,
		})
	}

	return &domain.PriceListRep{
		PaginationInfo: domain.PaginationInfoSt{
			Page:     resp.PaginationInfo.Page,
			PageSize: resp.PaginationInfo.PageSize,
		},
		Results: results,
	}, nil
}

func (c *PriceClient) GetPrice(ctx context.Context, productId, cityId string) (*domain.PriceMain, error) {
	resp, err := c.service.Get(ctx, &pb.ProductPriceGetReq{
		CityId:    cityId,
		ProductId: productId,
	})
	if err != nil {
		return nil, err
	}

	return &domain.PriceMain{
		ProductId: resp.ProductId,
		CityId:    resp.CityId,
		Price:     resp.Price,
	}, nil
}
