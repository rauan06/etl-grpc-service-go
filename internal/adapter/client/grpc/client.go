package grpc

type Client struct {
	Category *CategoryClient
	City     *CityClient
	Product  *ProductClient
	Price    *PriceClient
	Stock    *StockClient
}

func NewClient(category *CategoryClient, city *CityClient, product *ProductClient, price *PriceClient, stock *StockClient) *Client {
	return &Client{
		Category: category,
		City:     city,
		Product:  product,
		Price:    price,
		Stock:    stock,
	}
}
