package domain

type FullProduct struct {
	ProductMain ProductMain `json:"product_main"`
	Prices      []PriceMain `json:"price"`
	Stocks      []StockMain `json:"stock"`
}

type ProductMain struct {
	CreatedAt   string       `json:"created_at"`
	UpdatedAt   string       `json:"updated_at"`
	Deleted     bool         `json:"deleted"`
	ID          string       `json:"id"`
	Name        string       `json:"name"`
	Description string       `json:"description"`
	CategoryID  string       `json:"category_id"`
	Category    CategoryMain `json:"category"`
}

type ProductListRep struct {
	PaginationInfo PaginationInfoSt `json:"pagination_info"`
	Results        []ProductMain    `json:"results"`
}

func (p *FullProduct) IsValid() bool {
	if !p.ProductMain.IsValid() {
		return false
	}

	if p.Stocks == nil || len(p.Stocks) == 0 {
		return false
	}

	if p.Prices == nil || len(p.Prices) == 0 {
		return false
	}

	return true
}

func (p *ProductMain) IsValid() bool {
	if p.CreatedAt == "" || p.UpdatedAt == "" || p.ID == "" || p.Name == "" || p.CategoryID == "" || p.Description == "" {
		return false
	}

	if p.Deleted {
		return false
	}

	return p.Category.IsValid()
}
