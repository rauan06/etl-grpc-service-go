package main

import (
	"context"
	"fmt"
	"log"

	grpcClient "github.com/rauan06/etl-grpc-service-go/internal/adapter/client/grpc"
	"github.com/rauan06/etl-grpc-service-go/internal/core/domain"
)

func main() {
	data, err := grpcClient.NewStockClient(context.Background(), "0.0.0.0:5051")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(data.ListStocks(context.Background(), domain.ListParamsSt{Page: 2134295561239812903}, []string{}, []string{}))
}
