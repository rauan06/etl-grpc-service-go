package port

import (
	"context"
	"sync"

	"category/internal/core/domain"
)

type CategoryClient interface {
	ListCategories(ctx context.Context, params domain.ListParamsSt, ids []string) (*domain.CategoryListRep, error)
	GetCategory(ctx context.Context, id string) (*domain.CategoryMain, error)
}

type CategoryService interface {
	Run(ctx context.Context, wg *sync.WaitGroup)
	Status() int
	Stop()
}
