package grpc

import (
	"context"
	"strconv"

	"category/internal/adapter/handler"
	"category/internal/core/domain"
	pb "category/protos/price/v1/pb"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type PriceClient struct {
	conn    *grpc.ClientConn
	service pb.ProductPriceClient
}

func NewPriceClient(ctx context.Context, target string) (*PriceClient, error) {
	conn, err := grpc.NewClient(target, grpc.WithTransportCredentials(insecure.NewCredentials()))
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
	page, err := strconv.ParseInt(params.Page, 10, 64)
	if err != nil {
		return nil, handler.ErrParseInt64
	}
	pageSize, err := strconv.ParseInt(params.PageSize, 10, 64)
	if err != nil {
		return nil, handler.ErrParseInt64
	}

	resp, err := c.service.List(ctx, &pb.ProductPriceListReq{
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
			Page:     strconv.FormatInt(resp.PaginationInfo.Page, 10),
			PageSize: strconv.FormatInt(resp.PaginationInfo.PageSize, 10),
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
