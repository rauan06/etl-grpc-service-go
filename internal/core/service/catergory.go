package service

import (
	"context"
	"errors"
	"net/url"
	"strconv"
	"strings"

	"category/internal/core/domain"
	"category/internal/core/port"
)

type CategoryService struct {
	repo port.CategoryRepository
}

func NewCategoryService(repo port.CategoryRepository) *CategoryService {
	return &CategoryService{repo: repo}
}

func (s *CategoryService) ListCategories(ctx context.Context, params url.Values) (*domain.ProductCategoryListRep, error) {
	params["list_params.sort"] = filterValidSortParams(params["list_params.sort"])
	params["list_params.ids"] = filterValidIDs(params["list_params.ids"])

	return s.repo.ListCategories(ctx, params)
}

func (s *CategoryService) GetCategory(ctx context.Context, id int64) (*domain.ProductCategoryMain, error) {
	if id < 0 {
		return nil, errors.New("id cannot be negative")
	}

	return s.repo.GetCategory(ctx, id)
}

func filterValidSortParams(sortParams []string) []string {
	valid := make([]string, 0, len(sortParams))

	for _, param := range sortParams {
		if isValidCategoryField(strings.ToLower(param)) {
			valid = append(valid, param)
		}
	}

	return valid
}

func filterValidIDs(idParams []string) []string {
	valid := make([]string, 0, len(idParams))

	for _, idStr := range idParams {
		if idInt, err := strconv.ParseInt(idStr, 10, 64); err == nil && idInt >= 0 {
			valid = append(valid, idStr)
		}
	}

	return valid
}

func isValidCategoryField(field string) bool {
	switch field {
	case "created_at", "updated_at", "id", "name":
		return true
	default:
		return false
	}
}
