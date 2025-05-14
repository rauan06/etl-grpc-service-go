package redis

import (
	"context"
	"fmt"
	"time"

	"category/internal/core/port"
	"category/pkg/config"

	"github.com/go-redis/redis"
)

/**
 * Redis implements port.CacheRepository interface
 * and provides an access to the redis library
 */
type Redis struct {
	client *redis.Client
}

// New creates a new instance of Redis
func New(ctx context.Context, config *config.Config) (port.CacheRepository, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     config.RedisURI,
		Password: config.RedisPassword,
		DB:       config.RedisDB,
	})

	_, err := client.Ping().Result()
	if err != nil {
		return nil, err
	}

	return &Redis{client}, nil
}

// Set stores the value in the redis database
func (r *Redis) Set(ctx context.Context, key string, value []byte) error {
	ttl := 3 * time.Minute
	return r.client.Set(key, value, ttl).Err()
}

// Get retrieves the value from the redis database
func (r *Redis) Get(ctx context.Context, key string) ([]byte, error) {
	res, err := r.client.Get(key).Result()
	bytes := []byte(res)
	return bytes, err
}

// Delete removes the value from the redis database
func (r *Redis) Delete(ctx context.Context, key string) error {
	return r.client.Del(key).Err()
}

// DeleteByPrefix removes the value from the redis database with the given prefix
func (r *Redis) DeleteByPrefix(ctx context.Context, prefix string) error {
	var cursor uint64
	var keys []string

	for {
		var err error
		keys, cursor, err = r.client.Scan(cursor, prefix, 100).Result()
		if err != nil {
			return err
		}

		for _, key := range keys {
			err := r.client.Del(key).Err()
			if err != nil {
				return err
			}
		}

		if cursor == 0 {
			break
		}
	}

	return nil
}

func (r *Redis) Scan(pattern string) ([]string, error) {
	var cursor uint64
	var result = []string{}

	pattern = pattern + "*"

	for {
		keys, nextCursor, err := r.client.Scan(cursor, pattern, 30).Result()
		if err != nil {
			return nil, fmt.Errorf("scan failed: %w", err)
		}

		result = append(result, keys...)

		cursor = nextCursor
		if cursor == 0 {
			break
		}
	}

	return result, nil
}

// Close closes the connection to the redis database
func (r *Redis) Close() error {
	return r.client.Close()
}
