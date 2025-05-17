package http

type Client struct {
	Category *CategoryClient
	City     *CityClient
	Product  *ProductClient
	Price    *PriceClient
	Stock    *StockClient
}

func NewClient(
	caategory *CategoryClient,
	city *CityClient,
	product *ProductClient,
	price *PriceClient,
	stock *StockClient,
) *Client {
	return &Client{
		Category: caategory,
		City:     city,
		Product:  product,
		Price:    price,
		Stock:    stock,
	}
}
