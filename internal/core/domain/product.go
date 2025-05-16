package domain

type FullProduct struct {
	ID          string      `json:"id"`
	ProductMain ProductMain `json:"product_main"`
	City        CityMain    `json:"city"`
	Price       PriceMain   `json:"price"`
	Stock       StockMain   `json:"stock"`
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
	if p.ID == "" {
		return false
	}

	if !p.ProductMain.IsValid() || !p.City.IsValid() || !p.Price.IsValid() || !p.Stock.IsValid() {
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
