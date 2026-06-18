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
	vaultSnapshotsCachePrefixKey = "vaultSnapshots:"
	// VaultSnapshotsExpireTime expire time
	VaultSnapshotsExpireTime = 5 * time.Minute
)

var _ VaultSnapshotsCache = (*vaultSnapshotsCache)(nil)

// VaultSnapshotsCache cache interface
type VaultSnapshotsCache interface {
	Set(ctx context.Context, id uint64, data *model.VaultSnapshots, duration time.Duration) error
	Get(ctx context.Context, id uint64) (*model.VaultSnapshots, error)
	MultiGet(ctx context.Context, ids []uint64) (map[uint64]*model.VaultSnapshots, error)
	MultiSet(ctx context.Context, data []*model.VaultSnapshots, duration time.Duration) error
	Del(ctx context.Context, id uint64) error
	SetPlaceholder(ctx context.Context, id uint64) error
	IsPlaceholderErr(err error) bool
}

// vaultSnapshotsCache define a cache struct
type vaultSnapshotsCache struct {
	cache cache.Cache
}

// NewVaultSnapshotsCache new a cache
func NewVaultSnapshotsCache(cacheType *database.CacheType) VaultSnapshotsCache {
	jsonEncoding := encoding.JSONEncoding{}
	cachePrefix := ""

	cType := strings.ToLower(cacheType.CType)
	switch cType {
	case "redis":
		c := cache.NewRedisCache(cacheType.Rdb, cachePrefix, jsonEncoding, func() interface{} {
			return &model.VaultSnapshots{}
		})
		return &vaultSnapshotsCache{cache: c}
	case "memory":
		c := cache.NewMemoryCache(cachePrefix, jsonEncoding, func() interface{} {
			return &model.VaultSnapshots{}
		})
		return &vaultSnapshotsCache{cache: c}
	}

	return nil // no cache
}

// GetVaultSnapshotsCacheKey cache key
func (c *vaultSnapshotsCache) GetVaultSnapshotsCacheKey(id uint64) string {
	return vaultSnapshotsCachePrefixKey + utils.Uint64ToStr(id)
}

// Set write to cache
func (c *vaultSnapshotsCache) Set(ctx context.Context, id uint64, data *model.VaultSnapshots, duration time.Duration) error {
	if data == nil || id == 0 {
		return nil
	}
	cacheKey := c.GetVaultSnapshotsCacheKey(id)
	err := c.cache.Set(ctx, cacheKey, data, duration)
	if err != nil {
		return err
	}
	return nil
}

// Get cache value
func (c *vaultSnapshotsCache) Get(ctx context.Context, id uint64) (*model.VaultSnapshots, error) {
	var data *model.VaultSnapshots
	cacheKey := c.GetVaultSnapshotsCacheKey(id)
	err := c.cache.Get(ctx, cacheKey, &data)
	if err != nil {
		return nil, err
	}
	return data, nil
}

// MultiSet multiple set cache
func (c *vaultSnapshotsCache) MultiSet(ctx context.Context, data []*model.VaultSnapshots, duration time.Duration) error {
	valMap := make(map[string]interface{})
	for _, v := range data {
		cacheKey := c.GetVaultSnapshotsCacheKey(v.ID)
		valMap[cacheKey] = v
	}

	err := c.cache.MultiSet(ctx, valMap, duration)
	if err != nil {
		return err
	}

	return nil
}

// MultiGet multiple get cache, return key in map is id value
func (c *vaultSnapshotsCache) MultiGet(ctx context.Context, ids []uint64) (map[uint64]*model.VaultSnapshots, error) {
	var keys []string
	for _, v := range ids {
		cacheKey := c.GetVaultSnapshotsCacheKey(v)
		keys = append(keys, cacheKey)
	}

	itemMap := make(map[string]*model.VaultSnapshots)
	err := c.cache.MultiGet(ctx, keys, itemMap)
	if err != nil {
		return nil, err
	}

	retMap := make(map[uint64]*model.VaultSnapshots)
	for _, id := range ids {
		val, ok := itemMap[c.GetVaultSnapshotsCacheKey(id)]
		if ok {
			retMap[id] = val
		}
	}

	return retMap, nil
}

// Del delete cache
func (c *vaultSnapshotsCache) Del(ctx context.Context, id uint64) error {
	cacheKey := c.GetVaultSnapshotsCacheKey(id)
	err := c.cache.Del(ctx, cacheKey)
	if err != nil {
		return err
	}
	return nil
}

// SetPlaceholder set placeholder value to cache
func (c *vaultSnapshotsCache) SetPlaceholder(ctx context.Context, id uint64) error {
	cacheKey := c.GetVaultSnapshotsCacheKey(id)
	return c.cache.SetCacheWithNotFound(ctx, cacheKey)
}

// IsPlaceholderErr check if cache is placeholder error
func (c *vaultSnapshotsCache) IsPlaceholderErr(err error) bool {
	return errors.Is(err, cache.ErrPlaceholder)
}
