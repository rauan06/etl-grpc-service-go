package main

import (
	"context"
	"fmt"
	"log"

	grpcClient "github.com/rauan06/etl-grpc-service-go/internal/adapter/client/grpc"
	"github.com/rauan06/etl-grpc-service-go/internal/core/domain"
)

func main() {
	data, err := grpcClient.NewProductClient(context.Background(), "0.0.0.0:5050")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(data.ListProducts(context.Background(), domain.ListParamsSt{Page: 1}, []string{}, []string{}, false))
}
