package client

import (
	pb "category/internal/adapter/handler/grpc/product/v1"

	"google.golang.org/grpc"
)

type CategoryGrpcClient struct {
	conn    *grpc.ClientConn
	service pb.CategoryClient
}

func NewCategoryGrpcClient(url string) (*CategoryGrpcClient, error) {
	conn, err := grpc.Dial(url, grpc.WithInsecure())
	if err != nil {
		return nil, err
	}
	c := pb.NewCategoryClient(conn)
	return &CategoryGrpcClient{conn, c}, nil
}

func (c *CategoryGrpcClient) Close() {
	c.conn.Close()
}

func (c *CategoryGrpcClient) ListCategories() {}
func (c *CategoryGrpcClient) GetCategory()    {}
