package grpc

import "fmt"

type Client struct {
	Category *CategoryClient
	City     *CityClient
	Product  *ProductClient
	Price    *PriceClient
	Stock    *StockClient
}

func NewClient(category *CategoryClient, city *CityClient, product *ProductClient, price *PriceClient, stock *StockClient) *Client {
	fmt.Println(category, city, product, price, stock)

	return &Client{
		Category: category,
		City:     city,
		Product:  product,
		Price:    price,
		Stock:    stock,
	}
}
