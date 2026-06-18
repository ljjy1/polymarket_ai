package cache

import (
	"context"
	"strings"
	"time"

	"github.com/go-dev-frame/sponge/pkg/cache"
	"github.com/go-dev-frame/sponge/pkg/encoding"

	"be/internal/database"
)

const (
	// nonce cache prefix key, must end with a colon
	nonceCachePrefixKey = "nonce:"
	// DefaultNonceExpireTime default expire time
	DefaultNonceExpireTime = 5 * time.Minute
)

var _ NonceCache = (*nonceCache)(nil)

// NonceCache cache interface
type NonceCache interface {
	Set(ctx context.Context, address string, nonce string, duration time.Duration) error
	Get(ctx context.Context, address string) (string, error)
	Del(ctx context.Context, address string) error
}

// nonceCache define a cache struct
type nonceCache struct {
	cache cache.Cache
}

// NewNonceCache new a cache
func NewNonceCache(cacheType *database.CacheType) NonceCache {
	jsonEncoding := encoding.JSONEncoding{}
	cachePrefix := ""

	cType := strings.ToLower(cacheType.CType)
	switch cType {
	case "redis":
		c := cache.NewRedisCache(cacheType.Rdb, cachePrefix, jsonEncoding, func() interface{} {
			return ""
		})
		return &nonceCache{cache: c}
	case "memory":
		c := cache.NewMemoryCache(cachePrefix, jsonEncoding, func() interface{} {
			return ""
		})
		return &nonceCache{cache: c}
	}

	return nil // no cache
}

// GetNonceCacheKey cache key
func (c *nonceCache) GetNonceCacheKey(address string) string {
	return nonceCachePrefixKey + strings.ToLower(address)
}

// Set write to cache
func (c *nonceCache) Set(ctx context.Context, address string, nonce string, duration time.Duration) error {
	cacheKey := c.GetNonceCacheKey(address)
	return c.cache.Set(ctx, cacheKey, &nonce, duration)
}

// Get cache value
func (c *nonceCache) Get(ctx context.Context, address string) (string, error) {
	var nonce string
	cacheKey := c.GetNonceCacheKey(address)
	err := c.cache.Get(ctx, cacheKey, &nonce)
	if err != nil {
		return "", err
	}
	return nonce, nil
}

// Del delete cache
func (c *nonceCache) Del(ctx context.Context, address string) error {
	cacheKey := c.GetNonceCacheKey(address)
	return c.cache.Del(ctx, cacheKey)
}
