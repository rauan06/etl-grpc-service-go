package service

import (
	"category/internal/core/domain"
	"category/internal/core/port"
	"context"
	"strconv"
)

type CategoryService struct {
	repo port.CategoryRepository
}

func NewCategoryService(repo port.CategoryRepository) *CategoryService {
	return &CategoryService{
		repo,
	}
}

func (s *CategoryService) ListCategories(ctx context.Context, params domain.ProductListParamsSt, ids []string) (domain.ProductCategoryListRep, error) {
	newIds := make([]int64, len(ids))

	for _, rawId := range ids {
		id, err := strconv.ParseInt(rawId, 10, 64)
		if err != nil {
			return domain.ProductCategoryListRep{}, err
		}

		newIds = append(newIds, id)
	}

	_, err := s.repo.ListCategories(ctx, params, newIds)
	if err != nil {
		return domain.ProductCategoryListRep{}, err
	}

	panic("Implement me")
}
func (s *CategoryService) GetCategory(ctx context.Context, id int64) (domain.ProductCategoryMain, error) {
	panic("Implement me")
}
