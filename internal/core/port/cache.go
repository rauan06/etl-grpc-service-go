package port

import (
	"context"
)

//go:generate mockgen -source=cache.go -destination=mock/cache.go -package=mock

type CacheRepository interface {
	Set(ctx context.Context, key string, value []byte) error
	Get(ctx context.Context, key string) ([]byte, error)

	// Scans for pattern, return matchde keys
	Scan(pattern string) ([]string, error)
	Delete(ctx context.Context, key string) error
	DeleteByPrefix(ctx context.Context, prefix string) error
	Close() error
}
