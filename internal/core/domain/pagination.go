package domain

type ProductPaginationInfoSt struct {
	Page     int64 `json:"page"`
	PageSize int64 `json:"page_size"`
}

type ProductListParamsSt struct {
	Page     int64    `json:"page"`
	PageSize int64    `json:"page_size"`
	Sort     []string `json:"sort"`
}
