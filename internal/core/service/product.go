package service

import (
	"context"
	"errors"
	"net/url"

	"category/internal/core/domain"
	"category/internal/core/port"
)

type ProductService struct {
	repo port.ProductRepository
}

func NewProductService(repo port.ProductRepository) *ProductService {
	return &ProductService{repo: repo}
}

func (s *ProductService) ListProducts(ctx context.Context, params url.Values) (*domain.ProductProductListRep, error) {
	params["list_params.sort"] = filterValidSortParams(params["list_params.sort"])
	params["list_params.ids"] = filterValidIDs(params["list_params.ids"])
	params["category_ids"] = filterValidIDs(params["category_ids"])

	return s.repo.ListProducts(ctx, params)
}

func (s *ProductService) GetProducts(ctx context.Context, id int64) (*domain.ProductProductMain, error) {
	if id < 0 {
		return nil, errors.New("id cannot be negative")
	}

	return s.repo.GetProduct(ctx, id)
}
