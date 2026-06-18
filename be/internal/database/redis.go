package database

import (
	"context"
	"strings"
	"sync"
	"time"

	"github.com/go-dev-frame/sponge/pkg/goredis"
	"github.com/go-dev-frame/sponge/pkg/tracer"
	"github.com/redis/go-redis/v9"

	"be/internal/config"
)

var (
	// ErrCacheNotFound No hit cache
	ErrCacheNotFound = goredis.ErrRedisNotFound
)

var (
	redisCli     *goredis.Client
	redisCliOnce sync.Once

	cacheType     *CacheType
	cacheTypeOnce sync.Once
)

// CacheType cache type
type CacheType struct {
	CType string          // cache type  memory or redis
	Rdb   *goredis.Client // if CType=redis, Rdb cannot be empty
}

// InitCache initial cache
func InitCache(cType string) {
	cacheType = &CacheType{
		CType: cType,
	}

	if cType == "redis" {
		cacheType.Rdb = GetRedisCli()
	}
}

// GetCacheType get cacheType
func GetCacheType() *CacheType {
	if cacheType == nil {
		cacheTypeOnce.Do(func() {
			InitCache(config.Get().App.CacheType)
		})
	}

	return cacheType
}

// KeyPrefixHook is a Redis hook that adds a prefix to all keys.
type KeyPrefixHook struct {
	Prefix string
}

var _ redis.Hook = (*KeyPrefixHook)(nil)

// DialHook implements redis.Hook.
func (h *KeyPrefixHook) DialHook(next redis.DialHook) redis.DialHook {
	return next
}

// ProcessHook implements redis.Hook.
func (h *KeyPrefixHook) ProcessHook(next redis.ProcessHook) redis.ProcessHook {
	return func(ctx context.Context, cmd redis.Cmder) error {
		h.addPrefix(cmd)
		return next(ctx, cmd)
	}
}

// ProcessPipelineHook implements redis.Hook.
func (h *KeyPrefixHook) ProcessPipelineHook(next redis.ProcessPipelineHook) redis.ProcessPipelineHook {
	return func(ctx context.Context, cmds []redis.Cmder) error {
		for _, cmd := range cmds {
			h.addPrefix(cmd)
		}
		return next(ctx, cmds)
	}
}

// addPrefix prepends the prefix to all key arguments of the command.
func (h *KeyPrefixHook) addPrefix(cmd redis.Cmder) {
	if h.Prefix == "" {
		return
	}
	args := cmd.Args()
	if len(args) < 2 {
		return
	}
	// Skip commands that don't operate on keys
	switch cmd.Name() {
	case "AUTH", "ECHO", "PING", "QUIT", "SELECT", "INFO", "COMMAND",
		"CLIENT", "CONFIG", "SHUTDOWN", "SLAVEOF", "REPLCONF",
		"SUBSCRIBE", "UNSUBSCRIBE", "PSUBSCRIBE", "PUNSUBSCRIBE":
		return
	}
	// Add prefix to only key arguments (position 1 for single-key commands like SET/GET,
	// and all positional keys for multi-key commands like DEL/MSET/MGET).
	// Command args structure (index 0 is the command name itself):
	//   SET key value [EX|PX] seconds → keys at position 1
	//   GET key                       → keys at position 1
	//   DEL key [key ...]             → keys at positions 1,2,...
	//   MGET key [key ...]            → keys at positions 1,2,...
	//   MSET key value key value ...  → keys at odd positions 1,3,5,...
	for i := 1; i < len(args); i++ {
		if key, ok := args[i].(string); ok {
			// Only add prefix to the first key arg, skip flag-like args (ex, px, nx, xx, keepttl, etc.)
			// This is a safe simplification: only the first arg after the command name is always a key.
			prefixKey := true
			if i > 1 {
				// For multi-key commands, treat all args after position 0 as potential keys
				// only if the previous arg was also a string that looks like a key.
				// Simple heuristic: skip short flag-like strings that are common Redis modifiers.
				switch strings.ToLower(key) {
				case "ex", "px", "exat", "pxat", "persist", "keepttl",
					"nx", "xx", "gt", "lt", "chars",
					"asc", "desc", "alpha", "nosort", "withscores",
					"limit", "get", "before", "after", "left", "right":
					prefixKey = false
				}
			}
			if prefixKey {
				args[i] = h.Prefix + key
			}
		}
	}
}

// InitRedis connect redis
func InitRedis() {
	redisCfg := config.Get().Redis
	opts := []goredis.Option{
		goredis.WithDialTimeout(time.Duration(redisCfg.DialTimeout) * time.Second),
		goredis.WithReadTimeout(time.Duration(redisCfg.ReadTimeout) * time.Second),
		goredis.WithWriteTimeout(time.Duration(redisCfg.WriteTimeout) * time.Second),
	}
	if config.Get().App.EnableTrace {
		opts = append(opts, goredis.WithTracing(tracer.GetProvider()))
	}

	var err error
	redisCli, err = goredis.Init(redisCfg.Dsn, opts...)
	if err != nil {
		panic("goredis.Init error: " + err.Error())
	}

	// register key prefix hook
	if redisCfg.KeyPrefix != "" {
		redisCli.AddHook(&KeyPrefixHook{Prefix: redisCfg.KeyPrefix})
	}
}

// GetRedisCli get redis client
func GetRedisCli() *goredis.Client {
	if redisCli == nil {
		redisCliOnce.Do(func() {
			InitRedis()
		})
	}

	return redisCli
}

// CloseRedis close redis
func CloseRedis() error {
	return goredis.Close(redisCli)
}
