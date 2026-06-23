package service

import (
	"math"

	"github.com/go-dev-frame/sponge/pkg/logger"
)

// StrategyResult 策略生成结果
type StrategyResult struct {
	ID            int     `json:"id"`            // 数据库策略记录ID
	Action        string  `json:"action"`        // "buy_yes" | "buy_no" | "skip"
	Side          string  `json:"side"`          // "yes" | "no" | ""
	PositionSize  float64 `json:"positionSize"`  // USDC金额 (skip时为0)
	EntryPrice    float64 `json:"entryPrice"`    // 入场价格 (skip时为0)
	TakeProfit    float64 `json:"takeProfit"`    // 止盈价格 (skip时为0)
	StopLoss      float64 `json:"stopLoss"`      // 止损价格 (skip时为0)
	KellyFraction float64 `json:"kellyFraction"` // 凯利公式比例 (skip时为0)
	Edge          float64 `json:"edge"`          // 价差 (正负值)
	SkipReason    string  `json:"skipReason"`    // 跳过原因
}

// StrategyConfig 策略生成器配置
type StrategyConfig struct {
	MinEdge              float64 // 最小边阈值
	MinConfidence        float64 // 最小置信度阈值
	MaxPositionPct       float64 // 单笔最大仓位占比
	KellyMultiplier      float64 // 凯利公式乘数 (0~1)
	TakeProfitFactor     float64 // 止盈倍数 (edge的倍数)
	StopLossFactor       float64 // 止损倍数 (edge的倍数)
	PreResolutionMinutes int     // 结算前停止交易分钟数
}

// StrategyGenerator 策略生成器，将 AI 预测转化为可执行的交易策略
type StrategyGenerator struct {
	config StrategyConfig
}

// NewStrategyGenerator 创建策略生成器
func NewStrategyGenerator(cfg StrategyConfig) *StrategyGenerator {
	return &StrategyGenerator{config: cfg}
}

// DefaultStrategyConfig 返回默认策略配置
// 参数与 Python 参考版 (strategy_generator.py) 保持一致：
//   - MaxPositionPct:   0.10 (10%)  — 单笔不超过金库总资产的 10%
//   - TakeProfitFactor: 0.7         — 止盈 = entry_price + |edge| * 0.7
//   - StopLossFactor:   0.5         — 止损 = entry_price - |edge| * 0.5
//   - KellyMultiplier:  0.5         — 半凯利
func DefaultStrategyConfig() StrategyConfig {
	return StrategyConfig{
		MinEdge:              0.25,
		MinConfidence:        0.6,
		MaxPositionPct:       0.10, // Python 参考版: DEFAULT_MAX_POSITION_PCT = 0.10
		KellyMultiplier:      0.5,  // Python 参考版: DEFAULT_KELLY_MULTIPLIER = 0.5
		TakeProfitFactor:     0.7,  // Python 参考版: DEFAULT_TAKE_PROFIT_FACTOR = 0.7
		StopLossFactor:       0.5,  // Python 参考版: DEFAULT_STOP_LOSS_FACTOR = 0.5
		PreResolutionMinutes: 30,   // Python 参考版: 提前30分钟预解析平仓
	}
}

// Generate 根据 AI 预测和市场数据生成交易策略
func (g *StrategyGenerator) Generate(prediction *PredictionResult, marketYesPrice float64, vaultBalance float64) *StrategyResult {
	if prediction == nil {
		logger.Warn("策略生成：预测结果为 nil，跳过")
		return skipResult("预测结果为 nil")
	}

	logger.Info("开始生成策略",
		logger.Float64("predicted_probability", prediction.PredictedProbability),
		logger.Float64("confidence", prediction.Confidence),
		logger.String("recommended_action", prediction.RecommendedAction),
		logger.Float64("market_yes_price", marketYesPrice),
		logger.Float64("vault_balance", vaultBalance),
	)

	// 1. 计算 edge
	edge := prediction.PredictedProbability - marketYesPrice
	absEdge := math.Abs(edge)

	// 2. HARD GATE 1: edge 不足
	if absEdge < g.config.MinEdge {
		logger.Info("策略生成跳过：edge 不足",
			logger.Float64("abs_edge", absEdge),
			logger.Float64("min_edge", g.config.MinEdge),
		)
		return skipResult("edge 不足")
	}

	// 3. HARD GATE 2: 置信度不足
	if prediction.Confidence < g.config.MinConfidence {
		logger.Info("策略生成跳过：置信度不足",
			logger.Float64("confidence", prediction.Confidence),
			logger.Float64("min_confidence", g.config.MinConfidence),
		)
		return skipResult("置信度不足")
	}

	// 4. HARD GATE 3: recommended_action 与 edge 方向不一致
	if (prediction.RecommendedAction == "buy_yes" && edge <= 0) ||
		(prediction.RecommendedAction == "buy_no" && edge >= 0) {
		logger.Info("策略生成跳过：推荐操作与 edge 方向不一致",
			logger.String("recommended_action", prediction.RecommendedAction),
			logger.Float64("edge", edge),
		)
		return skipResult("推荐操作与 edge 方向不一致")
	}

	// 5. 确定 side
	side := "yes"
	entryPrice := marketYesPrice
	if edge < 0 {
		side = "no"
		entryPrice = 1 - marketYesPrice
	}

	// 6. Kelly 计算
	// 对于 YES: p = PredictedProbability, q = 1-p, b = (1-entryPrice)/entryPrice
	// 对于 NO:  p = 1-PredictedProbability, q = PredictedProbability, b = entryPrice/(1-entryPrice)
	var p, q, b float64
	if side == "yes" {
		p = prediction.PredictedProbability
		q = 1 - p
		b = (1 - entryPrice) / entryPrice
	} else {
		p = 1 - prediction.PredictedProbability
		q = prediction.PredictedProbability
		b = entryPrice / (1 - entryPrice)
	}

	// f = (b*p - q) / b
	var kellyFraction float64
	if b > 0 {
		kellyFraction = (b*p - q) / b
		// 截断负值（不押注负期望）
		if kellyFraction < 0 {
			kellyFraction = 0
		}
	} else {
		kellyFraction = 0
	}

	// 应用 KellyMultiplier（半凯利等）
	kellyFraction *= g.config.KellyMultiplier

	// 7. 计算最大仓位
	maxPosition := vaultBalance * g.config.MaxPositionPct

	// 8. 计算实际仓位
	positionSize := kellyFraction * vaultBalance
	if positionSize > maxPosition {
		positionSize = maxPosition
	}
	if positionSize < 0 {
		positionSize = 0
	}

	// 9. 计算止盈止损
	takeProfit := entryPrice + (absEdge * g.config.TakeProfitFactor)
	stopLoss := entryPrice - (absEdge * g.config.StopLossFactor)

	// 限制止盈止损在 [0, 1] 范围内
	if takeProfit > 1 {
		takeProfit = 1
	}
	if stopLoss < 0 {
		stopLoss = 0
	}

	action := "buy_" + side

	logger.Info("策略生成完成",
		logger.String("action", action),
		logger.String("side", side),
		logger.Float64("position_size", positionSize),
		logger.Float64("entry_price", entryPrice),
		logger.Float64("take_profit", takeProfit),
		logger.Float64("stop_loss", stopLoss),
		logger.Float64("kelly_fraction", kellyFraction),
		logger.Float64("edge", edge),
	)

	return &StrategyResult{
		Action:        action,
		Side:          side,
		PositionSize:  positionSize,
		EntryPrice:    entryPrice,
		TakeProfit:    takeProfit,
		StopLoss:      stopLoss,
		KellyFraction: kellyFraction,
		Edge:          edge,
		SkipReason:    "",
	}
}

// skipResult 返回跳过交易的结果
func skipResult(reason string) *StrategyResult {
	return &StrategyResult{
		Action:     "skip",
		Side:       "",
		SkipReason: reason,
	}
}
