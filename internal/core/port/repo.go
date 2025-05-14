package port

import "category/internal/core/domain"

type Repository interface {
	GetProductById(id string) (domain.FullProduct, error)
	InsertProduct(domain.FullProduct) error
}
