package domain

type StockMain struct {
	ProductId string `json:"product_id"`
	CityId    string `json:"city_id"`
	Value     int64  `json:"value"` // Format int64, but string-typed
}

type StockListRep struct {
	PaginationInfo PaginationInfoSt `json:"pagination_info"`
	Results        []StockMain      `json:"results"`
}
