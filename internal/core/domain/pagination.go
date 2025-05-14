package domain

type PaginationInfoSt struct {
	Page     int64 `json:"page,string"`
	PageSize int64 `json:"page_size,string"`
}

type ListParamsSt struct {
	Page     int64    `json:"page,string"`
	PageSize int64    `json:"page_size,string"`
	Sort     []string `json:"sort"`
}
