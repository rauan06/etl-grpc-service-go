package domain

type ProductCategoryMain struct {
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
	Deleted   bool   `json:"deleted"`
	ID        string `json:"id"`
	Name      string `json:"name"`
}

type ProductCategoryListRep struct {
	PaginationInfo ProductPaginationInfoSt `json:"pagination_info"`
	Results        []ProductCategoryMain   `json:"results"`
}
