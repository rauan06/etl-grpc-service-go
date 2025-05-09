package client

import (
	"context"

	pb "category/protos/product/v1"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type CityGrpcClient struct {
	conn    *grpc.ClientConn
	service pb.CityClient
}

func NewCityClient(ctx context.Context, url string) (*CityGrpcClient, error) {
	conn, err := grpc.NewClient(url, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}
	return &CityGrpcClient{
		conn:    conn,
		service: pb.NewCityClient(conn),
	}, nil
}

func (c *CityGrpcClient) Close() {
	c.conn.Close()
}

func (c *CityGrpcClient) ListCities(ctx context.Context, listParams *pb.CityListReq) (*pb.CityListRep, error) {
	resp, err := c.service.List(ctx, listParams)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (c *CityGrpcClient) GetCity(ctx context.Context, id string) (*pb.CityMain, error) {
	resp, err := c.service.Get(ctx, &pb.CityGetReq{Id: id})
	if err != nil {
		return nil, err
	}
	return resp, nil
}
