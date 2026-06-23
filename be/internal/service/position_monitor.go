package service

import (
	"context"
	"fmt"
	"math"
	"strconv"
	"sync"
	"time"

	"github.com/go-dev-frame/sponge/pkg/logger"
	"github.com/go-dev-frame/sponge/pkg/sgorm/query"
	"github.com/shopspring/decimal"

	"be/internal/dao"
	"be/internal/model"
	"be/internal/polymarket"
)

// MonitorResult 持仓监控检查结果
type MonitorResult struct {
	StrategyID   uint64  `json:"strategyId"`
	MarketID     uint64  `json:"marketId"`
	CurrentPrice float64 `json:"currentPrice"`
	Action       string  `json:"action"` // "hold" | "take_profit" | "stop_loss" | "pre_resolution"
	Reason       string  `json:"reason"`
	Pnl          float64 `json:"pnl"` // 如果平仓，计算盈亏
}

type priceSnapshot struct {
	price     float64
	timestamp time.Time
}

// PositionMonitor 持仓监控器，负责监控活跃持仓并自动执行止盈止损
// 使用 CLOB 实时价格检查止盈止损（与 Python 参考版一致），而非依赖 DB 存储的静态价格。
type PositionMonitor struct {
	strategyDao      dao.StrategiesDao
	tradeDao         dao.TradesDao
	marketDao        dao.MarketsDao
	polymarketClient *polymarket.Client // 用于获取 CLOB 实时价格

	priceCheckInterval   int // 价格检查间隔（秒）
	preResolutionMinutes int // 提前平仓分钟数
	alertPriceChangePct  int // 价格波动告警百分比

	mu         sync.Mutex
	lastPrices map[uint64]*priceSnapshot // strategyID -> 最近一次价格快照，用于波动检测
}

// NewPositionMonitor 创建持仓监控器
// polyClient 用于获取 CLOB 实时价格（Python 参考版：clob.get_current_price）。
func NewPositionMonitor(strategyDao dao.StrategiesDao, tradeDao dao.TradesDao, marketDao dao.MarketsDao,
	polyClient *polymarket.Client,
	priceCheckInterval, preResolutionMinutes, alertPriceChangePct int) *PositionMonitor {
	return &PositionMonitor{
		strategyDao:          strategyDao,
		tradeDao:             tradeDao,
		marketDao:            marketDao,
		polymarketClient:     polyClient,
		priceCheckInterval:   priceCheckInterval,
		preResolutionMinutes: preResolutionMinutes,
		alertPriceChangePct:  alertPriceChangePct,
		lastPrices:           make(map[uint64]*priceSnapshot),
	}
}

// Check 检查所有活跃持仓，返回需要操作的持仓列表
func (m *PositionMonitor) Check(ctx context.Context) ([]*MonitorResult, error) {
	// 查询所有活跃策略
	params := &query.Params{
		Columns: []query.Column{
			{
				Name:  "status",
				Exp:   "=",
				Value: "active",
			},
		},
	}
	strategies, total, err := m.strategyDao.GetByColumns(ctx, params)
	if err != nil {
		return nil, fmt.Errorf("查询活跃策略失败: %w", err)
	}

	if total == 0 || len(strategies) == 0 {
		logger.Info("持仓监控：无活跃持仓")
		return nil, nil
	}

	logger.Info("持仓监控：开始检查活跃持仓", logger.Int("count", len(strategies)))

	var results []*MonitorResult
	now := time.Now()

	for _, strategy := range strategies {
		result := m.checkStrategy(ctx, strategy, now)
		if result != nil {
			results = append(results, result)
		}
	}

	logger.Info("持仓监控：检查完成",
		logger.Int("active_count", len(strategies)),
		logger.Int("action_required", len(results)))

	return results, nil
}

// checkStrategy 检查单个策略的持仓状态
func (m *PositionMonitor) checkStrategy(ctx context.Context, strategy *model.Strategies, now time.Time) *MonitorResult {
	// 获取关联市场
	market, err := m.marketDao.GetByID(ctx, uint64(strategy.MarketID))
	if err != nil {
		logger.Warn("持仓监控：获取市场信息失败",
			logger.Err(err),
			logger.Any("market_id", strategy.MarketID),
			logger.Any("strategy_id", strategy.ID))
		return nil
	}

	if market == nil {
		logger.Warn("持仓监控：市场不存在",
			logger.Any("market_id", strategy.MarketID),
			logger.Any("strategy_id", strategy.ID))
		return nil
	}

	// 通过 CLOB 实时订单簿获取当前价格（与 Python 参考版一致）
	tokenID := market.PolymarketTokenID
	yesPrice := m.getCurrentPrice(tokenID)
	if yesPrice <= 0 {
		logger.Warn("持仓监控：CLOB 实时价格无效",
			logger.Any("strategy_id", strategy.ID),
			logger.String("side", strategy.Side),
			logger.String("token_id", tokenID))
		return nil
	}
	// 根据交易方向计算对应侧的价格
	currentPrice := yesPrice
	if strategy.Side == "no" {
		currentPrice = 1.0 - yesPrice
	}

	// 检查止盈
	if m.shouldTakeProfit(strategy, currentPrice) {
		pnl := m.calculatePnl(strategy, currentPrice)
		logger.Info("持仓监控：触发止盈",
			logger.Any("strategy_id", strategy.ID),
			logger.Float64("entry_price", m.decimalToFloat64(strategy.EntryPrice)),
			logger.Float64("current_price", currentPrice),
			logger.Float64("take_profit", m.decimalToFloat64(strategy.TakeProfit)),
			logger.Float64("pnl", pnl))
		return &MonitorResult{
			StrategyID:   strategy.ID,
			MarketID:     uint64(strategy.MarketID),
			CurrentPrice: currentPrice,
			Action:       "take_profit",
			Reason:       fmt.Sprintf("当前价格 %.4f 达到止盈线 %.4f", currentPrice, m.decimalToFloat64(strategy.TakeProfit)),
			Pnl:          pnl,
		}
	}

	// 检查止损
	if m.shouldStopLoss(strategy, currentPrice) {
		pnl := m.calculatePnl(strategy, currentPrice)
		logger.Info("持仓监控：触发止损",
			logger.Any("strategy_id", strategy.ID),
			logger.Float64("entry_price", m.decimalToFloat64(strategy.EntryPrice)),
			logger.Float64("current_price", currentPrice),
			logger.Float64("stop_loss", m.decimalToFloat64(strategy.StopLoss)),
			logger.Float64("pnl", pnl))
		return &MonitorResult{
			StrategyID:   strategy.ID,
			MarketID:     uint64(strategy.MarketID),
			CurrentPrice: currentPrice,
			Action:       "stop_loss",
			Reason:       fmt.Sprintf("当前价格 %.4f 达到止损线 %.4f", currentPrice, m.decimalToFloat64(strategy.StopLoss)),
			Pnl:          pnl,
		}
	}

	// 检查是否接近结算时间
	if market.TargetDate != nil && !market.TargetDate.IsZero() {
		remaining := time.Until(*market.TargetDate)
		if remaining > 0 && remaining.Minutes() <= float64(m.preResolutionMinutes) {
			pnl := m.calculatePnl(strategy, currentPrice)
			logger.Info("持仓监控：市场即将结算，提前平仓",
				logger.Any("strategy_id", strategy.ID),
				logger.Any("market_id", strategy.MarketID),
				logger.Float64("remaining_minutes", remaining.Minutes()),
				logger.Float64("pnl", pnl))
			return &MonitorResult{
				StrategyID:   strategy.ID,
				MarketID:     uint64(strategy.MarketID),
				CurrentPrice: currentPrice,
				Action:       "pre_resolution",
				Reason:       fmt.Sprintf("市场将于 %.0f 分钟后结算，提前平仓", remaining.Minutes()),
				Pnl:          pnl,
			}
		}
	}

	// 检查价格波动（仅记录告警，不触发操作）
	m.checkPriceVolatility(strategy.ID, currentPrice, now)

	return nil
}

// shouldTakeProfit 检查是否触发止盈
func (m *PositionMonitor) shouldTakeProfit(strategy *model.Strategies, currentPrice float64) bool {
	if strategy.TakeProfit == nil || strategy.TakeProfit.IsZero() {
		return false
	}
	tp, _ := strategy.TakeProfit.Float64()
	// 对于 yes 和 no 两侧，止盈都是当前价格达到或超过止盈线
	return currentPrice >= tp
}

// shouldStopLoss 检查是否触发止损
func (m *PositionMonitor) shouldStopLoss(strategy *model.Strategies, currentPrice float64) bool {
	if strategy.StopLoss == nil || strategy.StopLoss.IsZero() {
		return false
	}
	sl, _ := strategy.StopLoss.Float64()
	// 对于 yes 和 no 两侧，止损都是当前价格达到或低于止损线
	return currentPrice <= sl
}

// getCurrentPrice 通过 CLOB 订单簿获取实时价格（与 Python 参考版一致）。
// 使用订单簿买卖中间价作为当前市场价格，而非依赖 DB 中可能过时的静态价格。
// tokenID: CLOB 交易对 token ID（Yes 侧用 market.PolymarketTokenID）。
func (m *PositionMonitor) getCurrentPrice(tokenID string) float64 {
	if m.polymarketClient == nil || tokenID == "" {
		return 0
	}
	ob, err := m.polymarketClient.GetOrderBook(tokenID)
	if err != nil || ob == nil || len(ob.Bids) == 0 || len(ob.Asks) == 0 {
		logger.Warn("持仓监控：获取 CLOB 实时价格失败",
			logger.Err(err),
			logger.String("token_id", tokenID),
		)
		return 0
	}
	bestBid := parseMonitorPrice(ob.Bids[0].Price)
	bestAsk := parseMonitorPrice(ob.Asks[0].Price)
	if bestBid <= 0 || bestAsk <= 0 {
		return 0
	}
	return (bestBid + bestAsk) / 2
}

// parseMonitorPrice 辅助函数：将价格字符串解析为 float64。
func parseMonitorPrice(priceStr string) float64 {
	price, err := strconv.ParseFloat(priceStr, 64)
	if err != nil {
		return 0
	}
	return price
}

// calculatePnl 计算当前盈亏
// PnL = (PositionSize / EntryPrice) * CurrentPrice - PositionSize
func (m *PositionMonitor) calculatePnl(strategy *model.Strategies, currentPrice float64) float64 {
	if strategy.PositionSize == nil || strategy.PositionSize.IsZero() ||
		strategy.EntryPrice == nil || strategy.EntryPrice.IsZero() {
		return 0
	}

	entryPrice, _ := strategy.EntryPrice.Float64()
	if entryPrice == 0 {
		return 0
	}

	positionSize, _ := strategy.PositionSize.Float64()
	shares := positionSize / entryPrice
	currentValue := shares * currentPrice
	return currentValue - positionSize
}

// checkPriceVolatility 检查价格波动是否超过阈值，超过则记录告警
func (m *PositionMonitor) checkPriceVolatility(strategyID uint64, currentPrice float64, now time.Time) {
	m.mu.Lock()
	defer m.mu.Unlock()

	snapshot, exists := m.lastPrices[strategyID]
	if !exists {
		m.lastPrices[strategyID] = &priceSnapshot{
			price:     currentPrice,
			timestamp: now,
		}
		return
	}

	// 检查是否在5分钟窗口内
	elapsed := now.Sub(snapshot.timestamp)
	if elapsed > 5*time.Minute {
		// 超过5分钟，重置快照
		snapshot.price = currentPrice
		snapshot.timestamp = now
		return
	}

	// 计算价格变化百分比
	if snapshot.price > 0 {
		changePct := math.Abs(currentPrice-snapshot.price) / snapshot.price * 100
		if changePct >= float64(m.alertPriceChangePct) {
			logger.Warn("持仓监控：价格大幅波动",
				logger.Any("strategy_id", strategyID),
				logger.Float64("previous_price", snapshot.price),
				logger.Float64("current_price", currentPrice),
				logger.Float64("change_pct", changePct),
				logger.Float64("elapsed_seconds", elapsed.Seconds()))
		}
	}

	// 更新快照
	snapshot.price = currentPrice
	snapshot.timestamp = now
}

// ClosePosition 平仓指定持仓
func (m *PositionMonitor) ClosePosition(ctx context.Context, strategyID uint64, reason string) error {
	// 查询策略
	strategy, err := m.strategyDao.GetByID(ctx, strategyID)
	if err != nil {
		return fmt.Errorf("查询策略失败: %w", err)
	}
	if strategy == nil {
		return fmt.Errorf("策略不存在: %d", strategyID)
	}

	// 如果已平仓，直接返回
	if strategy.Status == "closed" {
		logger.Info("持仓监控：策略已平仓，跳过", logger.Any("strategy_id", strategyID))
		return nil
	}

	// 获取市场和 CLOB 实时价格
	market, err := m.marketDao.GetByID(ctx, uint64(strategy.MarketID))
	if err != nil {
		return fmt.Errorf("获取市场信息失败: %w", err)
	}
	yesPrice := m.getCurrentPrice(market.PolymarketTokenID)
	if yesPrice <= 0 {
		return fmt.Errorf("获取 CLOB 实时价格失败: token=%s", market.PolymarketTokenID)
	}
	currentPrice := yesPrice
	if strategy.Side == "no" {
		currentPrice = 1.0 - yesPrice
	}
	pnl := m.calculatePnl(strategy, currentPrice)

	// 更新策略状态为 closed
	strategy.Status = "closed"
	if err := m.strategyDao.UpdateByID(ctx, strategy); err != nil {
		return fmt.Errorf("更新策略状态失败: %w", err)
	}

	// 创建平仓 Trade 记录
	positionSize, _ := strategy.PositionSize.Float64()
	entryPrice, _ := strategy.EntryPrice.Float64()

	shares := decimal.NewFromFloat(0)
	if entryPrice > 0 {
		shares = decimal.NewFromFloat(positionSize / entryPrice)
	}

	amount := decimal.NewFromFloat(positionSize)
	price := decimal.NewFromFloat(currentPrice)
	pnlDecimal := decimal.NewFromFloat(pnl)

	now := time.Now()
	trade := &model.Trades{
		StrategyID:  int(strategyID),
		MarketID:    strategy.MarketID,
		Side:        strategy.Side,
		Action:      "sell",
		Amount:      &amount,
		Price:       &price,
		Shares:      &shares,
		Status:      "filled",
		Pnl:         &pnlDecimal,
		CloseReason: reason,
		FilledAt:    &now,
		ClosedAt:    &now,
	}

	if err := m.tradeDao.Create(ctx, trade); err != nil {
		// 创建交易记录失败，但策略状态已更新，记录错误但不回滚
		logger.Error("持仓监控：创建平仓交易记录失败",
			logger.Err(err),
			logger.Any("strategy_id", strategyID))
		return fmt.Errorf("创建平仓交易记录失败: %w", err)
	}

	logger.Info("持仓监控：平仓完成",
		logger.Any("strategy_id", strategyID),
		logger.String("reason", reason),
		logger.Float64("pnl", pnl))

	return nil
}

// decimalToFloat64 辅助方法：将 *decimal.Decimal 转为 float64
func (m *PositionMonitor) decimalToFloat64(d *decimal.Decimal) float64 {
	if d == nil {
		return 0
	}
	v, _ := d.Float64()
	return v
}
