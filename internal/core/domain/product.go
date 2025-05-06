package domain

type ProductProductMain struct {
	CreatedAt   string              `json:"created_at"`
	UpdatedAt   string              `json:"updated_at"`
	Deleted     bool                `json:"deleted"`
	ID          string              `json:"id"`
	Name        string              `json:"name"`
	Description string              `json:"description"`
	CategoryID  string              `json:"category_id"`
	Category    ProductCategoryMain `json:"category"`
}

type ProductProductListRep struct {
	PaginationInfo ProductPaginationInfoSt `json:"pagination_info"`
	Results        []ProductProductMain    `json:"results"`
}
