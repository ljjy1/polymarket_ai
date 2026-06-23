package model

import (
	"github.com/go-dev-frame/sponge/pkg/sgorm"
	"github.com/shopspring/decimal"
	"time"
)

// Strategies 交易策略表
type Strategies struct {
	sgorm.Model `gorm:"embedded"` // embed id and time

	PredictionID  int              `gorm:"column:prediction_id;type:int(11);not null;comment:关联预测记录ID（逻辑外键→predictions.id, 无物理约束）" json:"predictionID"`
	MarketID      int              `gorm:"column:market_id;type:int(11);not null;comment:关联市场ID（逻辑外键→markets.id, 无物理约束）" json:"marketID"`
	Action        string           `gorm:"column:action;type:varchar(16);not null;comment:策略动作: buy_yes-买入Yes, buy_no-买入No, skip-跳过" json:"action"`
	Side          string           `gorm:"column:side;type:varchar(8);comment:交易方向: yes-Yes方, no-No方（skip时为NULL）" json:"side"`
	PositionSize  *decimal.Decimal `gorm:"column:position_size;type:decimal(38,18);not null;comment:仓位大小（USDC金额, skip时为0）" json:"positionSize"`
	EntryPrice    *decimal.Decimal `gorm:"column:entry_price;type:decimal(38,18);not null;comment:入场价格（skip时为0）" json:"entryPrice"`
	TakeProfit    *decimal.Decimal `gorm:"column:take_profit;type:decimal(38,18);not null;comment:止盈价格（skip时为0）" json:"takeProfit"`
	StopLoss      *decimal.Decimal `gorm:"column:stop_loss;type:decimal(38,18);not null;comment:止损价格（skip时为0）" json:"stopLoss"`
	KellyFraction *decimal.Decimal `gorm:"column:kelly_fraction;type:decimal(38,18);not null;comment:凯利公式建议仓位比例（skip时为0）" json:"kellyFraction"`
	Edge          *decimal.Decimal `gorm:"column:edge;type:decimal(38,18);not null;comment:策略采用的价差（正负值，有符号）" json:"edge"`
	SkipReason    string           `gorm:"column:skip_reason;type:text;not null;comment:跳过交易的详细原因（未跳过时为空字符串）" json:"skipReason"`
	Status        string           `gorm:"column:status;type:varchar(16);default:pending;not null;comment:策略状态: skipped-跳过, pending-待执行, executing-执行中, active-已开仓, closed-已平仓, failed-执行失败" json:"status"`
	ExecutedAt    *time.Time       `gorm:"column:executed_at;type:datetime;comment:策略实际执行时间" json:"executedAt"`
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
