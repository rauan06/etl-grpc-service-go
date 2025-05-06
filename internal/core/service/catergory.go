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
	return &CategoryService{
		repo,
	}
}

func (s *CategoryService) ListCategories(ctx context.Context, params url.Values) (*domain.ProductCategoryListRep, error) {
	sortParams := params["sort"]
	validSortParams := []string{}

	for _, param := range sortParams {
		param = strings.ToLower(param)

		if IsValidCategoryEntity(param) {
			validSortParams = append(validSortParams, param)
		}
	}

	// Parse ids (multi)
	idsStr := params["ids"]
	validIds := []string{}

	// Dropping invalid ids
	for _, id := range idsStr {
		idInt, err := strconv.ParseInt(id, 10, 64)
		if err != nil && idInt >= 0 {
			validIds = append(validIds, id)
		}
	}

	params["ids"] = validIds
	params["sort"] = validSortParams

	return s.repo.ListCategories(ctx, params)
}

func (s *CategoryService) GetCategory(ctx context.Context, id int64) (*domain.ProductCategoryMain, error) {
	if id < 0 {
		return nil, errors.New("id cannot be negative")
	}

	return s.repo.GetCategory(ctx, id)
}

//	type ProductCategoryMain struct {
//		CreatedAt string `json:"created_at"`
//		UpdatedAt string `json:"updated_at"`
//		Deleted   bool   `json:"deleted"`
//		ID        string `json:"id"`
//		Name      string `json:"name"`
//	}
//
// Checks if item in sort parametr is a part of category
func IsValidCategoryEntity(entity string) bool {
	switch entity {
	case "created_at":
		return true
	case "updated_at":
		return true
	// case "Deleted":
	// 	return true
	case "id":
		return true
	case "name":
		return true
	default:
		return false
	}
}
