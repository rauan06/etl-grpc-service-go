package grpc

import (
	"context"

	pb "github.com/rauan06/etl-grpc-service-go/protos/price/v1/pb"

	"github.com/rauan06/etl-grpc-service-go/internal/core/domain"
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
		if prod == nil {
			continue
		}
		results = append(results, domain.PriceMain{
			ProductId: prod.ProductId,
			CityId:    prod.CityId,
			Price:     prod.Price,
		})
	}

	pagination := domain.PaginationInfoSt{}
	if resp.PaginationInfo != nil {
		pagination.Page = resp.PaginationInfo.Page
		pagination.PageSize = resp.PaginationInfo.PageSize
	}

	return &domain.PriceListRep{
		PaginationInfo: pagination,
		Results:        results,
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

	if resp == nil {
		return nil, nil // or fmt.Errorf("price not found for product %s in city %s", productId, cityId)
	}

	return &domain.PriceMain{
		ProductId: resp.ProductId,
		CityId:    resp.CityId,
		Price:     resp.Price,
	}, nil
}
