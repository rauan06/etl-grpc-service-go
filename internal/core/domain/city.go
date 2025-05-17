package domain

type CityMain struct {
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
	Deleted   bool   `json:"deleted"`
	ID        string `json:"id"`
	Name      string `json:"name"`
	Postcode  string `json:"postcode"`
}

type CityListRep struct {
	PaginationInfo PaginationInfoSt `json:"pagination_info"`
	Results        []CityMain       `json:"results"`
}

func (c *CityMain) IsValid() bool {
	if c.CreatedAt == "" || c.Deleted || c.UpdatedAt == "" || c.ID == "" || c.Name == "" {
		return false
	}

	return true
}
