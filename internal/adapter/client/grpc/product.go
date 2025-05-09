package client

import (
	pb "category/internal/adapter/handler/grpc/product/v1"

	"google.golang.org/grpc"
)

type ProductGrpcClient struct {
	conn    *grpc.ClientConn
	service pb.ProductClient
}

func NewProductGrpcClient(url string) (*ProductGrpcClient, error) {
	conn, err := grpc.Dial(url, grpc.WithInsecure())
	if err != nil {
		return nil, err
	}
	c := pb.NewProductClient(conn)
	return &ProductGrpcClient{conn, c}, nil
}

func (c *ProductGrpcClient) Close() {
	c.conn.Close()
}

func (c *ProductGrpcClient) ListProducts() {}
func (c *ProductGrpcClient) GetProduct()   {}
