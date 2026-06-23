package service

import (
	"context"
	"fmt"
	"time"

	"github.com/go-dev-frame/sponge/pkg/logger"
	"github.com/shopspring/decimal"

	polymarketSDK "github.com/0xNetuser/Polymarket-golang/polymarket"

	"be/internal/contract"
	"be/internal/dao"
	"be/internal/model"
	bePolymarket "be/internal/polymarket"
)

// TradeExecutor 交易执行器，负责通过 Polymarket CLOB API 执行交易策略
type TradeExecutor struct {
	polymarketClient *bePolymarket.Client
	tradeDao         dao.TradesDao
	strategyDao      dao.StrategiesDao
	vaultContract    *contract.VaultContractClient
}

// NewTradeExecutor 创建交易执行器
func NewTradeExecutor(polyClient *bePolymarket.Client, tradeDao dao.TradesDao, strategyDao dao.StrategiesDao, vaultContract *contract.VaultContractClient) *TradeExecutor {
	return &TradeExecutor{
		polymarketClient: polyClient,
		tradeDao:         tradeDao,
		strategyDao:      strategyDao,
		vaultContract:    vaultContract,
	}
}

// Execute 执行交易策略
//   - strategy: 生成的策略结果
//   - market: 目标市场
//   - 返回创建的 Trade 记录
func (e *TradeExecutor) Execute(ctx context.Context, strategy *StrategyResult, market *model.Markets) (*model.Trades, error) {
	// 1. 检查 strategy.Action 是否为 skip → 直接返回 nil
	if strategy.Action == "skip" {
		logger.Info("策略结果为跳过，不执行交易",
			logger.Int("strategy_id", strategy.ID),
			logger.String("skip_reason", strategy.SkipReason),
		)
		return nil, nil
	}

	logger.Info("开始执行交易策略",
		logger.Int("strategy_id", strategy.ID),
		logger.String("action", strategy.Action),
		logger.String("side", strategy.Side),
		logger.String("market_question", market.Question),
	)

	// 2. 更新策略状态为 "executing"
	strategyRecord, err := e.strategyDao.GetByID(ctx, uint64(strategy.ID))
	if err != nil {
		return nil, fmt.Errorf("获取策略记录失败: %w", err)
	}
	strategyRecord.Status = "executing"
	if err := e.strategyDao.UpdateByID(ctx, strategyRecord); err != nil {
		logger.Error("更新策略状态为 executing 失败", logger.Err(err), logger.Int("strategy_id", strategy.ID))
		return nil, fmt.Errorf("更新策略状态失败: %w", err)
	}

	// 3. 通过 polymarketClient 创建限价单（使用 GTC 类型）
	entryPrice := strategy.EntryPrice
	positionSize := strategy.PositionSize

	// 计算份额数量：shares = position_size / entry_price
	shares := positionSize / entryPrice

	// 映射 SDK side：buy -> BUY, sell -> SELL
	sdkSide := "BUY"

	result, err := e.polymarketClient.CreateOrder(market.PolymarketTokenID, entryPrice, shares, sdkSide)
	if err != nil {
		// 更新策略状态为 failed
		strategyRecord.Status = "failed"
		_ = e.strategyDao.UpdateByID(ctx, strategyRecord)
		logger.Error("创建 Polymarket 限价单失败",
			logger.Err(err),
			logger.String("token_id", market.PolymarketTokenID),
			logger.Float64("price", entryPrice),
			logger.Float64("size", shares),
		)
		return nil, fmt.Errorf("创建限价单失败: %w", err)
	}

	// 从响应中提取订单ID
	orderID := extractOrderID(result)

	logger.Info("Polymarket 限价单已创建",
		logger.String("order_id", orderID),
		logger.Float64("price", entryPrice),
		logger.Float64("size", shares),
	)

	// 4. 处理响应，创建 Trade 记录写入数据库
	entryPriceDec := decimal.NewFromFloat(entryPrice)
	sharesDec := decimal.NewFromFloat(shares)
	zeroDec := decimal.NewFromFloat(0)

	positionSizeDec := decimal.NewFromFloat(positionSize)

	trade := &model.Trades{
		StrategyID:        strategy.ID,
		MarketID:          int(market.ID),
		PolymarketOrderID: orderID,
		Side:              strategy.Side,
		Action:            strategy.Action,
		Amount:            &positionSizeDec,
		Price:             &entryPriceDec,
		Shares:            &sharesDec,
		Status:            "pending",
		Fee:               &zeroDec,
	}

	if err := e.tradeDao.Create(ctx, trade); err != nil {
		logger.Error("保存交易记录失败",
			logger.Err(err),
			logger.Int("strategy_id", strategy.ID),
		)
		return nil, fmt.Errorf("保存交易记录失败: %w", err)
	}

	logger.Info("交易记录已保存",
		logger.Uint64("trade_id", trade.ID),
		logger.String("order_id", orderID),
	)

	// 5. 更新策略状态为 "active"
	strategyRecord.Status = "active"
	nowTime := time.Now()
	strategyRecord.ExecutedAt = &nowTime
	if err := e.strategyDao.UpdateByID(ctx, strategyRecord); err != nil {
		logger.Error("更新策略状态为 active 失败", logger.Err(err), logger.Int("strategy_id", strategy.ID))
		return nil, fmt.Errorf("更新策略状态失败: %w", err)
	}

	// 6. 启动心跳（PostHeartbeat）
	if _, heartbeatErr := e.polymarketClient.ClobClient.PostHeartbeat(nil); heartbeatErr != nil {
		logger.Warn("发送 Polymarket 心跳失败", logger.Err(heartbeatErr))
	} else {
		logger.Debug("Polymarket 心跳发送成功")
	}

	// 7. 使用 logger 记录交易执行详情
	logger.Info("交易策略执行完成",
		logger.Int("strategy_id", strategy.ID),
		logger.String("action", strategy.Action),
		logger.String("side", strategy.Side),
		logger.String("order_id", orderID),
		logger.Float64("entry_price", entryPrice),
		logger.Float64("position_size", positionSize),
		logger.Uint64("trade_id", trade.ID),
		logger.String("market_question", market.Question),
	)

	return trade, nil
}

// CancelOrder 取消订单
func (e *TradeExecutor) CancelOrder(ctx context.Context, orderID string) error {
	_, err := e.polymarketClient.CancelOrder(orderID)
	if err != nil {
		logger.Error("取消 Polymarket 订单失败",
			logger.Err(err),
			logger.String("order_id", orderID),
		)
		return fmt.Errorf("取消订单失败: %w", err)
	}

	logger.Info("Polymarket 订单已取消", logger.String("order_id", orderID))
	return nil
}

// GetOrderStatus 获取订单状态
func (e *TradeExecutor) GetOrderStatus(ctx context.Context, orderID string) (string, error) {
	order, err := e.polymarketClient.ClobClient.GetOrder(orderID)
	if err != nil {
		logger.Error("获取 Polymarket 订单状态失败",
			logger.Err(err),
			logger.String("order_id", orderID),
		)
		return "", fmt.Errorf("获取订单状态失败: %w", err)
	}

	status := fmt.Sprintf("%v", order)
	logger.Debug("获取订单状态",
		logger.String("order_id", orderID),
		logger.Any("status", status),
	)

	return status, nil
}

// extractOrderID 从 PostOrderResultV2 中提取订单ID
func extractOrderID(result *polymarketSDK.PostOrderResultV2) string {
	if result == nil {
		return ""
	}

	// 尝试从 Response 中提取 orderID
	if result.Response != nil {
		if respMap, ok := result.Response.(map[string]interface{}); ok {
			if id, ok := respMap["orderID"]; ok {
				if idStr, ok := id.(string); ok {
					return idStr
				}
			}
		}
	}

	return ""
}
