package domain

type ProductCityMain struct {
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
	Deleted   bool   `json:"deleted"`
	ID        string `json:"id"`
	Name      string `json:"name"`
	Postcode  string `json:"postcode"`
}

type ProductCityListRep struct {
	PaginationInfo ProductPaginationInfoSt `json:"pagination_info"`
	Results        []ProductCityMain       `json:"results"`
}
