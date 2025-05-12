package domain

type PriceMain struct {
	ProductId string  `json:"product_id"`
	CityId    string  `json:"city_id"`
	Price     float64 `json:"price"` // format: double
}

type PriceListRep struct {
	PaginationInfo PaginationInfoSt `json:"pagination_info"`
	Results        []PriceMain      `json:"results"`
}
