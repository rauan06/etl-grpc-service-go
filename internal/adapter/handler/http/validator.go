package handler

import (
	"errors"
	"net/url"
	"strconv"
)

func validatePaginationParams(params url.Values) error {
	pageStr := params.Get("list_params.page")
	var page int64
	if pageStr != "" {
		p, err := strconv.ParseInt(pageStr, 10, 64)
		if err != nil {
			return errors.New("")
		}
		page = p
	}

	pageSizeStr := params.Get("list_params.page_size")
	var pageSize int64
	if pageSizeStr != "" {
		ps, err := strconv.ParseInt(pageSizeStr, 10, 64)
		if err != nil {
			return errors.New("page number must be a number")
		}
		pageSize = ps
	}

	if pageSize < 0 || page < 0 {
		return errors.New("page numbers cannot be negative")
	}

	return nil
}
