package initial

import (
	"strconv"

	"github.com/go-dev-frame/sponge/pkg/app"
	"github.com/go-dev-frame/sponge/pkg/logger"

	"be/internal/config"
	"be/internal/dao"
	"be/internal/database"
	"be/internal/polymarket"
	"be/internal/server"
	"be/internal/service"
	"be/internal/tasks"
)

// CreateServices 创建并返回所有需要启动的服务列表
// 包括 HTTP 服务、WebSocket 实时盘口监控和异步任务（Asynq）相关的 Server 及定时调度器
func CreateServices() []app.IServer {
	// 获取全局配置
	var cfg = config.Get()
	var servers []app.IServer

	// 1. 创建 HTTP 服务
	httpAddr := ":" + strconv.Itoa(cfg.HTTP.Port) // 拼接监听地址
	httpServer := server.NewHTTPServer(httpAddr,
		server.WithHTTPIsProd(cfg.App.Env == "prod"), // 根据环境设置是否为生产模式
		server.WithHTTPTLS(cfg.HTTP.TLS),             // 配置 TLS 证书（如有）
	)
	servers = append(servers, httpServer)

	// 2. 创建 WebSocket 实时盘口监控服务
	db := database.GetDB()
	marketsDao := dao.NewMarketsDao(db, nil)

	polyClient, err := polymarket.NewClient(
		cfg.Polymarket.ClobAPIURL,
		cfg.Polymarket.GammaAPIURL,
		cfg.Polymarket.ChainID,
		cfg.Polymarket.PrivateKey,
		cfg.Polymarket.APIKey,
		cfg.Polymarket.APISecret,
		cfg.Polymarket.Passphrase,
		cfg.Proxy.Addr,
	)
	if err != nil {
		logger.Error("[ws] 创建 Polymarket 客户端失败", logger.Err(err))
		// WebSocket 服务不阻塞应用启动，记录错误后继续
	} else {
		orderBookMonitor := service.NewOrderBookMonitor(polyClient, marketsDao, cfg.Polymarket.WSMarketURL)
		servers = append(servers, orderBookMonitor)
		logger.Info("[ws] 实时盘口监控服务已创建", logger.String("ws_url", cfg.Polymarket.WSMarketURL))
	}

	// --- Asynq 服务启动 --- 文档查看 https://go-sponge.com/zh/component/job/asynq.html#%E5%AE%9A%E6%97%B6%E4%BB%BB%E5%8A%A1

	// 3. 创建 Asynq Scheduler，用于定时将任务入队
	scheduler := tasks.NewAsynqScheduler()
	// 注册定时任务规则（cron 表达式与任务类型的映射）
	tasks.RegisterDailyCron(scheduler)   // 注册日常定时任务
	tasks.RegisterMonitorCron(scheduler) // 注册监控定时任务

	// 4. 创建 Asynq 任务处理服务器
	asynqServer := tasks.NewAsynqServer()

	// 5. 初始化任务依赖（如数据库连接、外部客户端等）
	dailyDeps, err := tasks.NewDailyTaskDeps()
	if err != nil {
		// 如果依赖初始化失败，记录错误并返回空服务列表
		logger.Error("[tasks] 创建日常任务依赖失败", logger.Err(err))
		return nil
	}
	monitorDeps := tasks.NewMonitorTaskDeps()

	// 设置全局任务依赖，供 HTTP API 手动触发使用
	tasks.SetDailyDeps(dailyDeps)
	tasks.SetMonitorDeps(monitorDeps)

	// 6. 注册任务处理函数（将任务类型与具体处理逻辑绑定）
	tasks.RegisterDailyTasks(asynqServer, dailyDeps)     // 注册日常任务处理器
	tasks.RegisterMonitorTasks(asynqServer, monitorDeps) // 注册监控任务处理器

	// 7. 在 goroutine 中启动 Scheduler和Asynq Server，使其异步运行，不阻塞主流程
	go func() {
		logger.Info("[asynq] scheduler starting...")
		scheduler.Run()
		logger.Info("[asynq] server starting...")
		asynqServer.Run()
	}()

	// 返回 HTTP 服务、WebSocket 监控等服务（Asynq 相关组件通过 goroutine 独立运行，由框架统一管理生命周期）
	return servers
}
