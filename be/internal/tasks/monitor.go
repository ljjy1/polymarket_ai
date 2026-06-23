package tasks

import (
	"context"
	"errors"
	"math/big"
	"time"

	"github.com/go-dev-frame/sponge/pkg/logger"
	"github.com/go-dev-frame/sponge/pkg/sasynq"
	"github.com/go-dev-frame/sponge/pkg/sgorm/query"
	"github.com/shopspring/decimal"

	"be/internal/config"
	"be/internal/contract"
	"be/internal/dao"
	"be/internal/database"
	"be/internal/model"
	"be/internal/polymarket"
	"be/internal/service"
)

// 任务类型常量
const (
	TypeMonitorPositions     = "monitor:positions"
	TypeMonitorVaultSnapshot = "monitor:vault-snapshot"
	TypeMonitorHealth        = "monitor:health"
	TypeMonitorSettlement    = "monitor:settlement"
)

// ErrVaultContractNotConfigured 金库合约未配置时的错误
var ErrVaultContractNotConfigured = errors.New("vault contract not configured, skip vault snapshot")

// ---------- 任务 Payload ----------

// MonitorPositionsPayload 持仓监控任务载荷
type MonitorPositionsPayload struct{}

// VaultSnapshotPayload 金库快照任务载荷
type VaultSnapshotPayload struct{}

// HealthCheckPayload 健康检查任务载荷
type HealthCheckPayload struct{}

// SettlementPayload 结算检查任务载荷
type SettlementPayload struct{}

// ---------- 依赖注入 ----------

// MonitorTaskDeps 监控任务所需的依赖
type MonitorTaskDeps struct {
	PositionMonitor   *service.PositionMonitor
	StrategiesDao     dao.StrategiesDao
	TradesDao         dao.TradesDao
	MarketsDao        dao.MarketsDao
	VaultSnapshotsDao dao.VaultSnapshotsDao
	SystemLogsDao     dao.SystemLogsDao
	VaultContract     *contract.VaultContractClient
	PolymarketClient  *polymarket.Client
}

// NewMonitorTaskDeps 创建监控任务依赖
func NewMonitorTaskDeps() *MonitorTaskDeps {
	db := database.GetDB()

	strategiesDao := dao.NewStrategiesDao(db, nil)
	tradesDao := dao.NewTradesDao(db, nil)
	marketsDao := dao.NewMarketsDao(db, nil)
	vaultSnapshotsDao := dao.NewVaultSnapshotsDao(db, nil)
	systemLogsDao := dao.NewSystemLogsDao(db, nil)

	// 初始化 Polymarket 客户端（用于 CLOB 实时价格监控和结算检查）
	var polyClient *polymarket.Client
	cfg := config.Get()
	if cfg.Polymarket.ClobAPIURL != "" && cfg.Polymarket.PrivateKey != "" {
		var err error
		polyClient, err = polymarket.NewClient(
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
			logger.Warn("[tasks] 初始化 Polymarket 客户端失败，监控任务将使用 DB 价格作为 fallback",
				logger.Err(err),
			)
		}
	}

	positionMonitor := service.NewPositionMonitor(
		strategiesDao,
		tradesDao,
		marketsDao,
		polyClient,
		60,  // priceCheckInterval: 60秒
		120, // preResolutionMinutes: 提前120分钟平仓
		5,   // alertPriceChangePct: 价格波动5%告警
	)

	// 初始化金库合约客户端
	var vaultContract *contract.VaultContractClient
	if cfg.Vault.RPCURL != "" && cfg.Vault.ContractAddress != "" {
		var err error
		vaultContract, err = contract.NewVaultContractClient(cfg.Vault.RPCURL, cfg.Vault.ContractAddress, cfg.Vault.StrategistPrivateKey)
		if err != nil {
			logger.Warn("[tasks] 初始化金库合约客户端失败，金库快照任务将不可用",
				logger.Err(err),
				logger.String("rpc_url", cfg.Vault.RPCURL),
			)
		} else {
			logger.Info("[tasks] 金库合约客户端初始化成功",
				logger.String("contract", cfg.Vault.ContractAddress),
			)
		}
	} else {
		logger.Warn("[tasks] 金库合约配置为空，金库快照任务将跳过链上数据读取")
	}

	return &MonitorTaskDeps{
		PositionMonitor:   positionMonitor,
		StrategiesDao:     strategiesDao,
		TradesDao:         tradesDao,
		MarketsDao:        marketsDao,
		VaultSnapshotsDao: vaultSnapshotsDao,
		SystemLogsDao:     systemLogsDao,
		VaultContract:     vaultContract,
		PolymarketClient:  polyClient,
	}
}

// ---------- 任务处理函数 ----------

// HandleMonitorPositions 监控活跃持仓，执行止盈止损逻辑
func HandleMonitorPositions(ctx context.Context, p *MonitorPositionsPayload, deps *MonitorTaskDeps) error {
	logger.Info("[tasks] 开始执行持仓监控任务", logger.String("type", TypeMonitorPositions))

	results, err := deps.PositionMonitor.Check(ctx)
	if err != nil {
		logger.Error("[tasks] 持仓监控检查失败", logger.Err(err))
		return err
	}

	for _, r := range results {
		logger.Info("[tasks] 持仓监控结果",
			logger.Uint64("strategy_id", r.StrategyID),
			logger.Uint64("market_id", r.MarketID),
			logger.String("action", r.Action),
			logger.String("reason", r.Reason),
			logger.Float64("pnl", r.Pnl),
		)
	}

	logger.Info("[tasks] 持仓监控任务完成", logger.Int("checked_count", len(results)))
	return nil
}

// HandleVaultSnapshot 从链上读取金库数据并记录快照
func HandleVaultSnapshot(ctx context.Context, p *VaultSnapshotPayload, deps *MonitorTaskDeps) error {
	logger.Info("[tasks] 开始执行金库快照任务", logger.String("type", TypeMonitorVaultSnapshot))

	if deps.VaultContract == nil {
		return ErrVaultContractNotConfigured
	}

	// 从链上读取金库数据
	totalAssets, err := deps.VaultContract.TotalAssets(ctx)
	if err != nil {
		logger.Error("[tasks] 读取 totalAssets 失败", logger.Err(err))
		return err
	}

	sharePrice, err := deps.VaultContract.SharePrice(ctx)
	if err != nil {
		logger.Error("[tasks] 读取 sharePrice 失败", logger.Err(err))
		return err
	}

	tvl, err := deps.VaultContract.AvailableBalance(ctx)
	if err != nil {
		logger.Error("[tasks] 读取 availableBalance 失败", logger.Err(err))
		return err
	}

	deployed, err := deps.VaultContract.StrategyDebt(ctx)
	if err != nil {
		logger.Error("[tasks] 读取 strategyDebt 失败", logger.Err(err))
		return err
	}

	now := time.Now()
	snapshot := &model.VaultSnapshots{
		TotalAssets:    decimalFromBigInt(totalAssets, 6),
		SharePrice:     decimalFromBigInt(sharePrice, 6),
		Tvl:            decimalFromBigInt(tvl, 6),
		DepositorCount: 0, // 暂无链上存款人数统计
		DeployedAmount: decimalFromBigInt(deployed, 6),
		SnapshotAt:     &now,
	}
	if err := deps.VaultSnapshotsDao.Create(ctx, snapshot); err != nil {
		logger.Error("[tasks] 保存金库快照失败", logger.Err(err))
		return err
	}

	logger.Info("[tasks] 金库快照记录完成",
		logger.String("total_assets", totalAssets.String()),
	)
	return nil
}

// decimalFromBigInt 将 big.Int（USDC 6位精度）转换为 decimal.Decimal
func decimalFromBigInt(val *big.Int, decimals int) *decimal.Decimal {
	if val == nil {
		return nil
	}
	d := decimal.NewFromBigInt(val, int32(-decimals))
	return &d
}

// HandleHealthCheck 健康检查
func HandleHealthCheck(ctx context.Context, p *HealthCheckPayload, deps *MonitorTaskDeps) error {
	logger.Info("[tasks] 开始执行健康检查任务", logger.String("type", TypeMonitorHealth))

	// 检查数据库连接
	db := database.GetDB()
	sqlDB, err := db.DB()
	if err != nil {
		logger.Error("[tasks] 健康检查：获取数据库连接失败", logger.Err(err))
		return err
	}
	if err := sqlDB.PingContext(ctx); err != nil {
		logger.Error("[tasks] 健康检查：数据库 Ping 失败", logger.Err(err))
		return err
	}

	logger.Info("[tasks] 健康检查完成：数据库连接正常")
	return nil
}

// HandleSettlement 检查已到期市场的结算结果，结算后归还资金到金库
func HandleSettlement(ctx context.Context, p *SettlementPayload, deps *MonitorTaskDeps) error {
	logger.Info("[tasks] 开始执行结算检查任务", logger.String("type", TypeMonitorSettlement))

	// 查询所有活跃策略，检查关联市场是否已结算
	strategies, _, err := deps.StrategiesDao.GetByColumns(ctx, &query.Params{
		Page: 1,
		Size: 9999,
	})
	if err != nil {
		logger.Error("[tasks] 查询策略列表失败", logger.Err(err))
		return err
	}

	for _, s := range strategies {
		if s.Status != "active" && s.Status != "executing" {
			continue
		}
		market, err := deps.MarketsDao.GetByID(ctx, uint64(s.MarketID))
		if err != nil {
			logger.Warn("[tasks] 结算检查：获取市场信息失败",
				logger.Err(err),
				logger.Int("market_id", s.MarketID),
			)
			continue
		}

		// 检查市场是否已到期
		isExpired := market.TargetDate != nil && time.Now().After(*market.TargetDate)

		// 检查市场是否已在 DB 中标记为 resolved
		isResolved := market.Status == "resolved" && market.Resolution != ""

		if !isExpired && !isResolved {
			continue
		}

		logger.Info("[tasks] 发现已到期/已结算市场",
			logger.Int("market_id", s.MarketID),
			logger.Uint64("strategy_id", uint64(s.ID)),
			logger.String("question", market.Question),
			logger.String("market_status", market.Status),
			logger.String("resolution", market.Resolution),
		)

		// 如果市场已到期但尚未检查 resolution，尝试通过 Polymarket API 查询
		resolution := market.Resolution
		if resolution == "" && deps.PolymarketClient != nil {
			resolution = checkMarketResolution(deps.PolymarketClient, market)
			if resolution != "" {
				market.Resolution = resolution
				market.Status = "resolved"
				_ = deps.MarketsDao.UpdateByID(ctx, market)
			}
		}

		if resolution == "" {
			logger.Info("[tasks] 市场尚未结算，跳过",
				logger.Int("market_id", s.MarketID),
			)
			continue
		}

		// 获取该策略的成交记录
		trades, _, err := deps.TradesDao.GetByColumns(ctx, &query.Params{
			Page: 1,
			Size: 9999,
			Columns: []query.Column{
				{Name: "strategy_id", Value: s.ID},
			},
		})
		if err != nil {
			logger.Error("[tasks] 查询交易记录失败", logger.Err(err), logger.Uint64("strategy_id", uint64(s.ID)))
			continue
		}

		if len(trades) == 0 {
			logger.Warn("[tasks] 策略无关联交易记录，直接标记为 closed",
				logger.Uint64("strategy_id", uint64(s.ID)),
			)
			s.Status = "closed"
			_ = deps.StrategiesDao.UpdateByID(ctx, s)
			continue
		}

		// 计算盈亏并归还资金
		for _, trade := range trades {
			if trade.Status != "filled" {
				continue
			}

			// 判断交易方向与结算结果是否一致
			// side=yes, resolution=yes -> 胜; side=yes, resolution=no -> 负
			// side=no, resolution=no -> 胜; side=no, resolution=yes -> 负
			isWin := (trade.Side == "yes" && resolution == "yes") || (trade.Side == "no" && resolution == "no")

			var pnl decimal.Decimal
			if isWin {
				// 获利: shares * (1 - entryPrice) 对于 YES
				// 获利: shares * entryPrice 对于 NO
				if trade.Shares != nil && trade.Price != nil {
					if trade.Side == "yes" {
						// 对于 NO 方向，结算为 YES 时: shares * price (NO shares 价值)
						one := decimal.NewFromFloat(1)
						pnl = trade.Shares.Mul(one.Sub(*trade.Price))
					} else {
						pnl = trade.Shares.Mul(*trade.Price)
					}
					trade.Pnl = &pnl
				}
			} else {
				// 亏损: 损失全部本金（positionSize）
				if trade.Amount != nil {
					loss := trade.Amount.Neg()
					trade.Pnl = &loss
				}
			}

			// 更新平仓信息
			nowTime := time.Now()
			trade.CloseReason = "settlement"
			trade.ClosedAt = &nowTime
			_ = deps.TradesDao.UpdateByID(ctx, trade)
		}

		// 计算总归还金额（已提取的 strategyDebt = positionSize）
		var totalReturnFees decimal.Decimal // 实际盈亏后应归还金额
		for _, trade := range trades {
			if trade.Pnl != nil {
				if trade.Amount != nil {
					totalReturnFees = totalReturnFees.Add(*trade.Amount).Add(*trade.Pnl)
				}
			}
		}

		// 调用 depositFromStrategy 归还资金到金库
		if deps.VaultContract != nil && totalReturnFees.IsPositive() {
			// USDC 6位精度
			returnFloat, _ := totalReturnFees.Float64()
			returnAmount := new(big.Int).SetUint64(uint64(returnFloat * 1e6))

			if err := deps.VaultContract.DepositFromStrategy(ctx, returnAmount); err != nil {
				logger.Error("[tasks] 归还资金到金库失败",
					logger.Err(err),
					logger.Uint64("strategy_id", uint64(s.ID)),
					logger.String("amount", returnAmount.String()),
				)
				s.Status = "failed"
				_ = deps.StrategiesDao.UpdateByID(ctx, s)
				continue
			}
			logger.Info("[tasks] 已归还资金到金库",
				logger.Uint64("strategy_id", uint64(s.ID)),
				logger.String("amount", returnAmount.String()),
			)
		}

		// 更新策略状态为 closed
		s.Status = "closed"
		nowTime := time.Now()
		s.ExecutedAt = &nowTime
		if err := deps.StrategiesDao.UpdateByID(ctx, s); err != nil {
			logger.Error("[tasks] 更新策略状态为 closed 失败",
				logger.Err(err),
				logger.Uint64("strategy_id", uint64(s.ID)),
			)
		}
	}

	logger.Info("[tasks] 结算检查任务完成")
	return nil
}

// checkMarketResolution 通过 Polymarket API 检查市场结算结果
func checkMarketResolution(client *polymarket.Client, market *model.Markets) string {
	if client == nil || client.ClobClient == nil {
		return ""
	}
	// 通过 CLOB API 获取市场信息
	order, err := client.ClobClient.GetOrder(market.PolymarketConditionID)
	if err != nil {
		logger.Warn("[tasks] 查询 Polymarket 市场结算状态失败",
			logger.Err(err),
			logger.Int("market_id", int(market.ID)),
		)
		return ""
	}
	_ = order // 实际解析逻辑依赖 SDK 返回结构，暂时简化
	// 注：Polymarket CLOB API 的具体返回值结构需要对接时确认
	// 简化处理：如果已到期且订单不存在，认为已结算，resolution 未知
	return ""
}

// ---------- 任务注册 ----------

// RegisterMonitorTasks 注册监控任务处理函数
func RegisterMonitorTasks(srv *sasynq.Server, deps *MonitorTaskDeps) {
	mux := srv.Mux()

	sasynq.RegisterTaskHandler(mux, TypeMonitorPositions, sasynq.HandleFunc(func(ctx context.Context, p MonitorPositionsPayload) error {
		return HandleMonitorPositions(ctx, &p, deps)
	}))
	sasynq.RegisterTaskHandler(mux, TypeMonitorVaultSnapshot, sasynq.HandleFunc(func(ctx context.Context, p VaultSnapshotPayload) error {
		return HandleVaultSnapshot(ctx, &p, deps)
	}))
	sasynq.RegisterTaskHandler(mux, TypeMonitorHealth, sasynq.HandleFunc(func(ctx context.Context, p HealthCheckPayload) error {
		return HandleHealthCheck(ctx, &p, deps)
	}))
	sasynq.RegisterTaskHandler(mux, TypeMonitorSettlement, sasynq.HandleFunc(func(ctx context.Context, p SettlementPayload) error {
		return HandleSettlement(ctx, &p, deps)
	}))
}

// RegisterMonitorCron 注册监控定时任务
func RegisterMonitorCron(scheduler *sasynq.Scheduler) {
	// 每5分钟执行持仓监控
	scheduler.RegisterTask("@every 5m", TypeMonitorPositions, MonitorPositionsPayload{})
	// 每30分钟记录金库快照
	scheduler.RegisterTask("@every 30m", TypeMonitorVaultSnapshot, VaultSnapshotPayload{})
	// 每10分钟执行健康检查
	scheduler.RegisterTask("@every 10m", TypeMonitorHealth, HealthCheckPayload{})
	// 每1小时检查结算
	scheduler.RegisterTask("@every 1h", TypeMonitorSettlement, SettlementPayload{})

	logger.Info("[tasks] 监控定时任务注册完成")
}
