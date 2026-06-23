package initial

import (
	"context"
	"time"

	"github.com/go-dev-frame/sponge/pkg/app"
	"github.com/go-dev-frame/sponge/pkg/logger"
	"github.com/go-dev-frame/sponge/pkg/tracer"

	"be/internal/config"
	"be/internal/database"
)

// Close 在服务退出后释放所有资源，返回一组关闭函数，供应用框架按顺序调用。
func Close(servers []app.IServer) []app.Close {
	var closes []app.Close

	// 关闭各个服务器（HTTP/gRPC等）
	for _, s := range servers {
		closes = append(closes, s.Stop)
	}

	// 关闭 Asynq 任务调度
	// Asynq server 的 Shutdown 由应用框架统一管理，此处仅记录日志
	closes = append(closes, func() error {
		logger.Info("[asynq] shutdown")
		return nil
	})

	// 关闭数据库连接
	closes = append(closes, func() error {
		return database.CloseDB()
	})

	// 关闭 Redis 连接（仅在缓存类型为 redis 时执行）
	if config.Get().App.CacheType == "redis" {
		closes = append(closes, func() error {
			return database.CloseRedis()
		})
	}

	// 关闭链路追踪（Tracing）
	if config.Get().App.EnableTrace {
		closes = append(closes, func() error {
			ctx, _ := context.WithTimeout(context.Background(), 2*time.Second) //nolint
			return tracer.Close(ctx)
		})
	}

	// 同步日志缓冲区，确保日志落盘
	closes = append(closes, func() error {
		return logger.Sync()
	})

	return closes
}
