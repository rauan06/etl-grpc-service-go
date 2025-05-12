package domain

type PaginationInfoSt struct {
	Page     string `json:"page"`
	PageSize string `json:"page_size"`
}

type ListParamsSt struct {
	Page     string   `json:"page"`
	PageSize string   `json:"page_size"`
	Sort     []string `json:"sort"`
}
