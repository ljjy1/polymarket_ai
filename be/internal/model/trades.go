package model

import (
	"github.com/go-dev-frame/sponge/pkg/sgorm"
	"github.com/shopspring/decimal"
	"time"
)

// Trades 交易执行记录表
type Trades struct {
	sgorm.Model `gorm:"embedded"` // embed id and time

	StrategyID        int              `gorm:"column:strategy_id;type:int(11);not null" json:"strategyID"`                      // 关联策略ID（逻辑外键→strategies.id, 无物理约束）
	MarketID          int              `gorm:"column:market_id;type:int(11);not null" json:"marketID"`                          // 关联市场ID（逻辑外键→markets.id, 无物理约束）
	PolymarketOrderID string           `gorm:"column:polymarket_order_id;type:varchar(128);not null" json:"polymarketOrderID"`  // Polymarket CLOB 订单ID（全局唯一）
	Side              string           `gorm:"column:side;type:varchar(8);not null" json:"side"`                                // 交易方向: yes-Yes方, no-No方
	Action            string           `gorm:"column:action;type:varchar(8);not null" json:"action"`                            // 操作类型: buy-买入建仓, sell-卖出平仓
	Amount            *decimal.Decimal `gorm:"column:amount;type:decimal(38,18);not null" json:"amount"`                        // 交易金额（买入时为USDC支出, 卖出时为USDC收入）
	Price             *decimal.Decimal `gorm:"column:price;type:decimal(38,18);not null" json:"price"`                          // 成交单价
	Shares            *decimal.Decimal `gorm:"column:shares;type:decimal(38,18);not null" json:"shares"`                        // 成交份额数量
	Status            string           `gorm:"column:status;type:varchar(16);not null" json:"status"`                           // 订单状态: pending-待成交, filled-已成交, partial-部分成交, cancelled-已取消, failed-失败
	Fee               *decimal.Decimal `gorm:"column:fee;type:decimal(38,18);default:0.000000000000000000;not null" json:"fee"` // 交易手续费（USDC）
	Pnl               *decimal.Decimal `gorm:"column:pnl;type:decimal(38,18)" json:"pnl"`                                       // 盈亏金额（平仓时有值, USDC）
	CloseReason       string           `gorm:"column:close_reason;type:varchar(32)" json:"closeReason"`                         // 平仓原因: take_profit-止盈, stop_loss-止损, pre_resolution-到期前, manual-手动
	FilledAt          *time.Time       `gorm:"column:filled_at;type:datetime" json:"filledAt"`                                  // 订单成交时间
	ClosedAt          *time.Time       `gorm:"column:closed_at;type:datetime" json:"closedAt"`                                  // 平仓时间
}

// TradesColumnNames Whitelist for custom query fields to prevent sql injection attacks
var TradesColumnNames = map[string]bool{
	"id":                  true,
	"created_at":          true,
	"updated_at":          true,
	"deleted_at":          true,
	"strategy_id":         true,
	"market_id":           true,
	"polymarket_order_id": true,
	"side":                true,
	"action":              true,
	"amount":              true,
	"price":               true,
	"shares":              true,
	"status":              true,
	"fee":                 true,
	"pnl":                 true,
	"close_reason":        true,
	"filled_at":           true,
	"closed_at":           true,
}
