package domain

type ProductPaginationInfoSt struct {
	Page     string `json:"page"`
	PageSize string `json:"page_size"`
}

type ProductListParamsSt struct {
	Page     string   `json:"page"`
	PageSize string   `json:"page_size"`
	Sort     []string `json:"sort"`
}
