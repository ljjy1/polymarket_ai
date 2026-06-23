package tasks

import (
	"time"

	"github.com/go-dev-frame/sponge/pkg/sasynq"
	"github.com/hibiken/asynq"

	"be/internal/config"
)

// NewAsynqServer 创建 Asynq 服务端
func NewAsynqServer() *sasynq.Server {
	asynqCfg := config.Get().Asynq
	redisCfg := sasynq.RedisConfig{
		Addr:     asynqCfg.Addr,
		Password: asynqCfg.Password,
		DB:       asynqCfg.DB,
	}

	// 使用默认配置并覆盖并发数
	srvCfg := sasynq.DefaultServerConfig()
	srvCfg.Config.Concurrency = asynqCfg.Concurrency

	srv := sasynq.NewServer(redisCfg, srvCfg)
	srv.Use(sasynq.LoggingMiddleware())

	return srv
}

// NewAsynqScheduler 创建定时任务调度器
func NewAsynqScheduler() *sasynq.Scheduler {
	asynqCfg := config.Get().Asynq
	redisCfg := sasynq.RedisConfig{
		Addr:     asynqCfg.Addr,
		Password: asynqCfg.Password,
		DB:       asynqCfg.DB,
	}

	loc, err := time.LoadLocation("Asia/Shanghai")
	if err != nil {
		panic("load timezone location failed: " + err.Error())
	}

	return sasynq.NewScheduler(redisCfg,
		sasynq.WithSchedulerOptions(&asynq.SchedulerOpts{
			Location: loc,
		}),
	)
}
