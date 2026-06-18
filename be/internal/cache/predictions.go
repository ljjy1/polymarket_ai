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
	predictionsCachePrefixKey = "predictions:"
	// PredictionsExpireTime expire time
	PredictionsExpireTime = 5 * time.Minute
)

var _ PredictionsCache = (*predictionsCache)(nil)

// PredictionsCache cache interface
type PredictionsCache interface {
	Set(ctx context.Context, id uint64, data *model.Predictions, duration time.Duration) error
	Get(ctx context.Context, id uint64) (*model.Predictions, error)
	MultiGet(ctx context.Context, ids []uint64) (map[uint64]*model.Predictions, error)
	MultiSet(ctx context.Context, data []*model.Predictions, duration time.Duration) error
	Del(ctx context.Context, id uint64) error
	SetPlaceholder(ctx context.Context, id uint64) error
	IsPlaceholderErr(err error) bool
}

// predictionsCache define a cache struct
type predictionsCache struct {
	cache cache.Cache
}

// NewPredictionsCache new a cache
func NewPredictionsCache(cacheType *database.CacheType) PredictionsCache {
	jsonEncoding := encoding.JSONEncoding{}
	cachePrefix := ""

	cType := strings.ToLower(cacheType.CType)
	switch cType {
	case "redis":
		c := cache.NewRedisCache(cacheType.Rdb, cachePrefix, jsonEncoding, func() interface{} {
			return &model.Predictions{}
		})
		return &predictionsCache{cache: c}
	case "memory":
		c := cache.NewMemoryCache(cachePrefix, jsonEncoding, func() interface{} {
			return &model.Predictions{}
		})
		return &predictionsCache{cache: c}
	}

	return nil // no cache
}

// GetPredictionsCacheKey cache key
func (c *predictionsCache) GetPredictionsCacheKey(id uint64) string {
	return predictionsCachePrefixKey + utils.Uint64ToStr(id)
}

// Set write to cache
func (c *predictionsCache) Set(ctx context.Context, id uint64, data *model.Predictions, duration time.Duration) error {
	if data == nil || id == 0 {
		return nil
	}
	cacheKey := c.GetPredictionsCacheKey(id)
	err := c.cache.Set(ctx, cacheKey, data, duration)
	if err != nil {
		return err
	}
	return nil
}

// Get cache value
func (c *predictionsCache) Get(ctx context.Context, id uint64) (*model.Predictions, error) {
	var data *model.Predictions
	cacheKey := c.GetPredictionsCacheKey(id)
	err := c.cache.Get(ctx, cacheKey, &data)
	if err != nil {
		return nil, err
	}
	return data, nil
}

// MultiSet multiple set cache
func (c *predictionsCache) MultiSet(ctx context.Context, data []*model.Predictions, duration time.Duration) error {
	valMap := make(map[string]interface{})
	for _, v := range data {
		cacheKey := c.GetPredictionsCacheKey(v.ID)
		valMap[cacheKey] = v
	}

	err := c.cache.MultiSet(ctx, valMap, duration)
	if err != nil {
		return err
	}

	return nil
}

// MultiGet multiple get cache, return key in map is id value
func (c *predictionsCache) MultiGet(ctx context.Context, ids []uint64) (map[uint64]*model.Predictions, error) {
	var keys []string
	for _, v := range ids {
		cacheKey := c.GetPredictionsCacheKey(v)
		keys = append(keys, cacheKey)
	}

	itemMap := make(map[string]*model.Predictions)
	err := c.cache.MultiGet(ctx, keys, itemMap)
	if err != nil {
		return nil, err
	}

	retMap := make(map[uint64]*model.Predictions)
	for _, id := range ids {
		val, ok := itemMap[c.GetPredictionsCacheKey(id)]
		if ok {
			retMap[id] = val
		}
	}

	return retMap, nil
}

// Del delete cache
func (c *predictionsCache) Del(ctx context.Context, id uint64) error {
	cacheKey := c.GetPredictionsCacheKey(id)
	err := c.cache.Del(ctx, cacheKey)
	if err != nil {
		return err
	}
	return nil
}

// SetPlaceholder set placeholder value to cache
func (c *predictionsCache) SetPlaceholder(ctx context.Context, id uint64) error {
	cacheKey := c.GetPredictionsCacheKey(id)
	return c.cache.SetCacheWithNotFound(ctx, cacheKey)
}

// IsPlaceholderErr check if cache is placeholder error
func (c *predictionsCache) IsPlaceholderErr(err error) bool {
	return errors.Is(err, cache.ErrPlaceholder)
}
