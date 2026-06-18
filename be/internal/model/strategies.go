package model

import (
	"github.com/go-dev-frame/sponge/pkg/sgorm"
	"github.com/shopspring/decimal"
	"time"
)

// Strategies 交易策略表
type Strategies struct {
	sgorm.Model `gorm:"embedded"` // embed id and time

	PredictionID  int              `gorm:"column:prediction_id;type:int(11);not null" json:"predictionID"`          // 关联预测记录ID（逻辑外键→predictions.id, 无物理约束）
	MarketID      int              `gorm:"column:market_id;type:int(11);not null" json:"marketID"`                  // 关联市场ID（逻辑外键→markets.id, 无物理约束）
	Action        string           `gorm:"column:action;type:varchar(16);not null" json:"action"`                   // 策略动作: buy_yes-买入Yes, buy_no-买入No, skip-跳过
	Side          string           `gorm:"column:side;type:varchar(8)" json:"side"`                                 // 交易方向: yes-Yes方, no-No方（skip时为NULL）
	PositionSize  *decimal.Decimal `gorm:"column:position_size;type:decimal(38,18);not null" json:"positionSize"`   // 仓位大小（USDC金额, skip时为0）
	EntryPrice    *decimal.Decimal `gorm:"column:entry_price;type:decimal(38,18);not null" json:"entryPrice"`       // 入场价格（skip时为0）
	TakeProfit    *decimal.Decimal `gorm:"column:take_profit;type:decimal(38,18);not null" json:"takeProfit"`       // 止盈价格（skip时为0）
	StopLoss      *decimal.Decimal `gorm:"column:stop_loss;type:decimal(38,18);not null" json:"stopLoss"`           // 止损价格（skip时为0）
	KellyFraction *decimal.Decimal `gorm:"column:kelly_fraction;type:decimal(38,18);not null" json:"kellyFraction"` // 凯利公式建议仓位比例（skip时为0）
	Edge          *decimal.Decimal `gorm:"column:edge;type:decimal(38,18);not null" json:"edge"`                    // 策略采用的价差（正负值，有符号）
	SkipReason    string           `gorm:"column:skip_reason;type:text;not null" json:"skipReason"`                 // 跳过交易的详细原因（未跳过时为空字符串）
	Status        string           `gorm:"column:status;type:varchar(16);default:pending;not null" json:"status"`   // 策略状态: skipped-跳过, pending-待执行, executing-执行中, active-已开仓, closed-已平仓, failed-执行失败
	ExecutedAt    *time.Time       `gorm:"column:executed_at;type:datetime" json:"executedAt"`                      // 策略实际执行时间
}

// StrategiesColumnNames Whitelist for custom query fields to prevent sql injection attacks
var StrategiesColumnNames = map[string]bool{
	"id":             true,
	"created_at":     true,
	"updated_at":     true,
	"deleted_at":     true,
	"prediction_id":  true,
	"market_id":      true,
	"action":         true,
	"side":           true,
	"position_size":  true,
	"entry_price":    true,
	"take_profit":    true,
	"stop_loss":      true,
	"kelly_fraction": true,
	"edge":           true,
	"skip_reason":    true,
	"status":         true,
	"executed_at":    true,
}
