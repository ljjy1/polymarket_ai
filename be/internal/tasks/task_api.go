package tasks

import (
	"context"
	"sync"

	"github.com/go-dev-frame/sponge/pkg/logger"
)

// 全局持有任务依赖，供 HTTP API 直接触发调用

var (
	taskDepsMu  sync.Mutex
	dailyDeps   *DailyTaskDeps
	monitorDeps *MonitorTaskDeps
)

// SetDailyDeps 设置全局日常任务依赖
func SetDailyDeps(deps *DailyTaskDeps) {
	taskDepsMu.Lock()
	defer taskDepsMu.Unlock()
	dailyDeps = deps
}

// SetMonitorDeps 设置全局监控任务依赖
func SetMonitorDeps(deps *MonitorTaskDeps) {
	taskDepsMu.Lock()
	defer taskDepsMu.Unlock()
	monitorDeps = deps
}

// ---------- 日常任务触发 ----------

// TriggerScan 手动触发市场扫描
func TriggerScan(ctx context.Context) error {
	taskDepsMu.Lock()
	deps := dailyDeps
	taskDepsMu.Unlock()
	if deps == nil {
		return nil
	}
	return HandleScanTask(ctx, &ScanPayload{}, deps)
}

// TriggerPredict 手动触发 AI 预测
func TriggerPredict(ctx context.Context, marketID int) error {
	taskDepsMu.Lock()
	deps := dailyDeps
	taskDepsMu.Unlock()
	if deps == nil {
		logger.Warn("[task-api] 日常任务依赖未初始化，跳过 AI 预测")
		return nil
	}
	return HandlePredictTask(ctx, &PredictPayload{MarketID: marketID}, deps)
}

// TriggerStrategy 手动触发策略生成
func TriggerStrategy(ctx context.Context, predictionID, marketID int) error {
	taskDepsMu.Lock()
	deps := dailyDeps
	taskDepsMu.Unlock()
	if deps == nil {
		logger.Warn("[task-api] 日常任务依赖未初始化，跳过策略生成")
		return nil
	}
	return HandleStrategyTask(ctx, &StrategyPayload{PredictionID: predictionID, MarketID: marketID}, deps)
}

// TriggerExecute 手动触发交易执行
func TriggerExecute(ctx context.Context, marketID, strategyID int) error {
	taskDepsMu.Lock()
	deps := dailyDeps
	taskDepsMu.Unlock()
	if deps == nil {
		logger.Warn("[task-api] 日常任务依赖未初始化，跳过交易执行")
		return nil
	}
	return HandleExecuteTask(ctx, &ExecutePayload{MarketID: marketID, StrategyID: strategyID}, deps)
}

// ---------- 监控任务触发 ----------

// TriggerMonitorPositions 手动触发持仓监控
func TriggerMonitorPositions(ctx context.Context) error {
	taskDepsMu.Lock()
	deps := monitorDeps
	taskDepsMu.Unlock()
	if deps == nil {
		logger.Warn("[task-api] 监控任务依赖未初始化，跳过持仓监控")
		return nil
	}
	return HandleMonitorPositions(ctx, &MonitorPositionsPayload{}, deps)
}

// TriggerVaultSnapshot 手动触发金库快照
func TriggerVaultSnapshot(ctx context.Context) error {
	taskDepsMu.Lock()
	deps := monitorDeps
	taskDepsMu.Unlock()
	if deps == nil {
		logger.Warn("[task-api] 监控任务依赖未初始化，跳过金库快照")
		return nil
	}
	return HandleVaultSnapshot(ctx, &VaultSnapshotPayload{}, deps)
}

// TriggerHealthCheck 手动触发健康检查
func TriggerHealthCheck(ctx context.Context) error {
	taskDepsMu.Lock()
	deps := monitorDeps
	taskDepsMu.Unlock()
	if deps == nil {
		logger.Warn("[task-api] 监控任务依赖未初始化，跳过健康检查")
		return nil
	}
	return HandleHealthCheck(ctx, &HealthCheckPayload{}, deps)
}

// TriggerSettlement 手动触发结算检查
func TriggerSettlement(ctx context.Context) error {
	taskDepsMu.Lock()
	deps := monitorDeps
	taskDepsMu.Unlock()
	if deps == nil {
		logger.Warn("[task-api] 监控任务依赖未初始化，跳过结算检查")
		return nil
	}
	return HandleSettlement(ctx, &SettlementPayload{}, deps)
}
