package tasks

import (
	"github.com/go-dev-frame/sponge/pkg/sasynq"

	"be/internal/config"
)

// NewAsynqClient 创建 Asynq 客户端
func NewAsynqClient() *sasynq.Client {
	asynqCfg := config.Get().Asynq
	return sasynq.NewClient(sasynq.RedisConfig{
		Addr:     asynqCfg.Addr,
		Password: asynqCfg.Password,
		DB:       asynqCfg.DB,
	})
}
