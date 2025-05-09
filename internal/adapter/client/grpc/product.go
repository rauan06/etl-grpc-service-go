package client

import (
	"context"

	pb "category/protos/product/v1"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type ProductClient struct {
	conn    *grpc.ClientConn
	service pb.ProductClient
}

func NewProductClient(ctx context.Context, target string) (*ProductClient, error) {
	conn, err := grpc.NewClient(target, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}
	return &ProductClient{
		conn:    conn,
		service: pb.NewProductClient(conn),
	}, nil
}

func (c *ProductClient) Close() {
	c.conn.Close()
}

func (c *ProductClient) ListProducts(ctx context.Context, listParams *pb.ProductListReq) (*pb.ProductListRep, error) {
	resp, err := c.service.List(ctx, listParams)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (c *ProductClient) GetProduct(ctx context.Context, id string) (*pb.ProductMain, error) {
	resp, err := c.service.Get(ctx, &pb.ProductGetReq{Id: id})
	if err != nil {
		return nil, err
	}
	return resp, nil
}
