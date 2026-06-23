package model

import (
	"github.com/go-dev-frame/sponge/pkg/sgorm"
	"github.com/shopspring/decimal"
	"time"
)

// Trades 交易执行记录表
type Trades struct {
	sgorm.Model `gorm:"embedded"` // embed id and time

	StrategyID        int              `gorm:"column:strategy_id;type:int(11);not null;comment:关联策略ID（逻辑外键→strategies.id, 无物理约束）" json:"strategyID"`
	MarketID          int              `gorm:"column:market_id;type:int(11);not null;comment:关联市场ID（逻辑外键→markets.id, 无物理约束）" json:"marketID"`
	PolymarketOrderID string           `gorm:"column:polymarket_order_id;type:varchar(128);not null;comment:Polymarket CLOB 订单ID（全局唯一）" json:"polymarketOrderID"`
	Side              string           `gorm:"column:side;type:varchar(8);not null;comment:交易方向: yes-Yes方, no-No方" json:"side"`
	Action            string           `gorm:"column:action;type:varchar(8);not null;comment:操作类型: buy-买入建仓, sell-卖出平仓" json:"action"`
	Amount            *decimal.Decimal `gorm:"column:amount;type:decimal(38,18);not null;comment:交易金额（买入时为USDC支出, 卖出时为USDC收入）" json:"amount"`
	Price             *decimal.Decimal `gorm:"column:price;type:decimal(38,18);not null;comment:成交单价" json:"price"`
	Shares            *decimal.Decimal `gorm:"column:shares;type:decimal(38,18);not null;comment:成交份额数量" json:"shares"`
	Status            string           `gorm:"column:status;type:varchar(16);not null;comment:订单状态: pending-待成交, filled-已成交, partial-部分成交, cancelled-已取消, failed-失败" json:"status"`
	Fee               *decimal.Decimal `gorm:"column:fee;type:decimal(38,18);default:0;not null;comment:交易手续费（USDC）" json:"fee"`
	Pnl               *decimal.Decimal `gorm:"column:pnl;type:decimal(38,18);comment:盈亏金额（平仓时有值, USDC）" json:"pnl"`
	CloseReason       string           `gorm:"column:close_reason;type:varchar(32);comment:平仓原因: take_profit-止盈, stop_loss-止损, pre_resolution-到期前, manual-手动" json:"closeReason"`
	FilledAt          *time.Time       `gorm:"column:filled_at;type:datetime;comment:订单成交时间" json:"filledAt"`
	ClosedAt          *time.Time       `gorm:"column:closed_at;type:datetime;comment:平仓时间" json:"closedAt"`
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
