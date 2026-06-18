package cache

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/go-dev-frame/sponge/pkg/cache"
	"github.com/go-dev-frame/sponge/pkg/encoding"
	"github.com/go-dev-frame/sponge/pkg/utils"

	"be/internal/database"
	"be/internal/model"
)

const (
	// cache prefix key, must end with a colon
	marketsCachePrefixKey = "markets:"
	// MarketsExpireTime expire time
	MarketsExpireTime = 5 * time.Minute
)

var _ MarketsCache = (*marketsCache)(nil)

// MarketsCache cache interface
type MarketsCache interface {
	Set(ctx context.Context, id uint64, data *model.Markets, duration time.Duration) error
	Get(ctx context.Context, id uint64) (*model.Markets, error)
	MultiGet(ctx context.Context, ids []uint64) (map[uint64]*model.Markets, error)
	MultiSet(ctx context.Context, data []*model.Markets, duration time.Duration) error
	Del(ctx context.Context, id uint64) error
	SetPlaceholder(ctx context.Context, id uint64) error
	IsPlaceholderErr(err error) bool
}

// marketsCache define a cache struct
type marketsCache struct {
	cache cache.Cache
}

// NewMarketsCache new a cache
func NewMarketsCache(cacheType *database.CacheType) MarketsCache {
	jsonEncoding := encoding.JSONEncoding{}
	cachePrefix := ""

	cType := strings.ToLower(cacheType.CType)
	switch cType {
	case "redis":
		c := cache.NewRedisCache(cacheType.Rdb, cachePrefix, jsonEncoding, func() interface{} {
			return &model.Markets{}
		})
		return &marketsCache{cache: c}
	case "memory":
		c := cache.NewMemoryCache(cachePrefix, jsonEncoding, func() interface{} {
			return &model.Markets{}
		})
		return &marketsCache{cache: c}
	}

	return nil // no cache
}

// GetMarketsCacheKey cache key
func (c *marketsCache) GetMarketsCacheKey(id uint64) string {
	return marketsCachePrefixKey + utils.Uint64ToStr(id)
}

// Set write to cache
func (c *marketsCache) Set(ctx context.Context, id uint64, data *model.Markets, duration time.Duration) error {
	if data == nil || id == 0 {
		return nil
	}
	cacheKey := c.GetMarketsCacheKey(id)
	err := c.cache.Set(ctx, cacheKey, data, duration)
	if err != nil {
		return err
	}
	return nil
}

// Get cache value
func (c *marketsCache) Get(ctx context.Context, id uint64) (*model.Markets, error) {
	var data *model.Markets
	cacheKey := c.GetMarketsCacheKey(id)
	err := c.cache.Get(ctx, cacheKey, &data)
	if err != nil {
		return nil, err
	}
	return data, nil
}

// MultiSet multiple set cache
func (c *marketsCache) MultiSet(ctx context.Context, data []*model.Markets, duration time.Duration) error {
	valMap := make(map[string]interface{})
	for _, v := range data {
		cacheKey := c.GetMarketsCacheKey(v.ID)
		valMap[cacheKey] = v
	}

	err := c.cache.MultiSet(ctx, valMap, duration)
	if err != nil {
		return err
	}

	return nil
}

// MultiGet multiple get cache, return key in map is id value
func (c *marketsCache) MultiGet(ctx context.Context, ids []uint64) (map[uint64]*model.Markets, error) {
	var keys []string
	for _, v := range ids {
		cacheKey := c.GetMarketsCacheKey(v)
		keys = append(keys, cacheKey)
	}

	itemMap := make(map[string]*model.Markets)
	err := c.cache.MultiGet(ctx, keys, itemMap)
	if err != nil {
		return nil, err
	}

	retMap := make(map[uint64]*model.Markets)
	for _, id := range ids {
		val, ok := itemMap[c.GetMarketsCacheKey(id)]
		if ok {
			retMap[id] = val
		}
	}

	return retMap, nil
}

// Del delete cache
func (c *marketsCache) Del(ctx context.Context, id uint64) error {
	cacheKey := c.GetMarketsCacheKey(id)
	err := c.cache.Del(ctx, cacheKey)
	if err != nil {
		return err
	}
	return nil
}

// SetPlaceholder set placeholder value to cache
func (c *marketsCache) SetPlaceholder(ctx context.Context, id uint64) error {
	cacheKey := c.GetMarketsCacheKey(id)
	return c.cache.SetCacheWithNotFound(ctx, cacheKey)
}

// IsPlaceholderErr check if cache is placeholder error
func (c *marketsCache) IsPlaceholderErr(err error) bool {
	return errors.Is(err, cache.ErrPlaceholder)
}
