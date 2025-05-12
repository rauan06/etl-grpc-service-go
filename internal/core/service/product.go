package service

import (
	"category/internal/core/domain"
	"category/internal/core/port"
	"context"
)

type ProductService struct {
	client port.ProductClient
}

func NewProductService(client port.ProductClient) *ProductService {
	return &ProductService{
		client,
	}
}

func (s *ProductService) ListProducts(ctx context.Context, params domain.ListParamsSt, ids, categoryIDs []string, withCategory bool) (*domain.ProductListRep, error) {
	return s.client.ListProducts(context.Background(), params, ids, categoryIDs, withCategory)
}

func (s *ProductService) GetProduct(ctx context.Context, id string) (*domain.ProductMain, error) {
	// if id < 0 {
	// 	return nil, fmt.Errorf("id cannot be empty")
	// }
	// stringId := strconv.FormatInt(id, 10)

	return s.client.GetProduct(context.Background(), id)
}
