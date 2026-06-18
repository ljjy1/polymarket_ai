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
	systemLogsCachePrefixKey = "systemLogs:"
	// SystemLogsExpireTime expire time
	SystemLogsExpireTime = 5 * time.Minute
)

var _ SystemLogsCache = (*systemLogsCache)(nil)

// SystemLogsCache cache interface
type SystemLogsCache interface {
	Set(ctx context.Context, id uint64, data *model.SystemLogs, duration time.Duration) error
	Get(ctx context.Context, id uint64) (*model.SystemLogs, error)
	MultiGet(ctx context.Context, ids []uint64) (map[uint64]*model.SystemLogs, error)
	MultiSet(ctx context.Context, data []*model.SystemLogs, duration time.Duration) error
	Del(ctx context.Context, id uint64) error
	SetPlaceholder(ctx context.Context, id uint64) error
	IsPlaceholderErr(err error) bool
}

// systemLogsCache define a cache struct
type systemLogsCache struct {
	cache cache.Cache
}

// NewSystemLogsCache new a cache
func NewSystemLogsCache(cacheType *database.CacheType) SystemLogsCache {
	jsonEncoding := encoding.JSONEncoding{}
	cachePrefix := ""

	cType := strings.ToLower(cacheType.CType)
	switch cType {
	case "redis":
		c := cache.NewRedisCache(cacheType.Rdb, cachePrefix, jsonEncoding, func() interface{} {
			return &model.SystemLogs{}
		})
		return &systemLogsCache{cache: c}
	case "memory":
		c := cache.NewMemoryCache(cachePrefix, jsonEncoding, func() interface{} {
			return &model.SystemLogs{}
		})
		return &systemLogsCache{cache: c}
	}

	return nil // no cache
}

// GetSystemLogsCacheKey cache key
func (c *systemLogsCache) GetSystemLogsCacheKey(id uint64) string {
	return systemLogsCachePrefixKey + utils.Uint64ToStr(id)
}

// Set write to cache
func (c *systemLogsCache) Set(ctx context.Context, id uint64, data *model.SystemLogs, duration time.Duration) error {
	if data == nil || id == 0 {
		return nil
	}
	cacheKey := c.GetSystemLogsCacheKey(id)
	err := c.cache.Set(ctx, cacheKey, data, duration)
	if err != nil {
		return err
	}
	return nil
}

// Get cache value
func (c *systemLogsCache) Get(ctx context.Context, id uint64) (*model.SystemLogs, error) {
	var data *model.SystemLogs
	cacheKey := c.GetSystemLogsCacheKey(id)
	err := c.cache.Get(ctx, cacheKey, &data)
	if err != nil {
		return nil, err
	}
	return data, nil
}

// MultiSet multiple set cache
func (c *systemLogsCache) MultiSet(ctx context.Context, data []*model.SystemLogs, duration time.Duration) error {
	valMap := make(map[string]interface{})
	for _, v := range data {
		cacheKey := c.GetSystemLogsCacheKey(v.ID)
		valMap[cacheKey] = v
	}

	err := c.cache.MultiSet(ctx, valMap, duration)
	if err != nil {
		return err
	}

	return nil
}

// MultiGet multiple get cache, return key in map is id value
func (c *systemLogsCache) MultiGet(ctx context.Context, ids []uint64) (map[uint64]*model.SystemLogs, error) {
	var keys []string
	for _, v := range ids {
		cacheKey := c.GetSystemLogsCacheKey(v)
		keys = append(keys, cacheKey)
	}

	itemMap := make(map[string]*model.SystemLogs)
	err := c.cache.MultiGet(ctx, keys, itemMap)
	if err != nil {
		return nil, err
	}

	retMap := make(map[uint64]*model.SystemLogs)
	for _, id := range ids {
		val, ok := itemMap[c.GetSystemLogsCacheKey(id)]
		if ok {
			retMap[id] = val
		}
	}

	return retMap, nil
}

// Del delete cache
func (c *systemLogsCache) Del(ctx context.Context, id uint64) error {
	cacheKey := c.GetSystemLogsCacheKey(id)
	err := c.cache.Del(ctx, cacheKey)
	if err != nil {
		return err
	}
	return nil
}

// SetPlaceholder set placeholder value to cache
func (c *systemLogsCache) SetPlaceholder(ctx context.Context, id uint64) error {
	cacheKey := c.GetSystemLogsCacheKey(id)
	return c.cache.SetCacheWithNotFound(ctx, cacheKey)
}

// IsPlaceholderErr check if cache is placeholder error
func (c *systemLogsCache) IsPlaceholderErr(err error) bool {
	return errors.Is(err, cache.ErrPlaceholder)
}
