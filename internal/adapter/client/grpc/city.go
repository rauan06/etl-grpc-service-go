package client

import (
	pb "category/internal/adapter/handler/grpc/product/v1"

	"google.golang.org/grpc"
)

type CityGrpcClient struct {
	conn    *grpc.ClientConn
	service pb.CityClient
}

func NewCityGrpcClient(url string) (*CityGrpcClient, error) {
	conn, err := grpc.Dial(url, grpc.WithInsecure())
	if err != nil {
		return nil, err
	}
	c := pb.NewCityClient(conn)
	return &CityGrpcClient{conn, c}, nil
}

func (c *CityGrpcClient) Close() {
	c.conn.Close()
}

func (c *CityGrpcClient) ListCities() {}
func (c *CityGrpcClient) GetCity()    {}
