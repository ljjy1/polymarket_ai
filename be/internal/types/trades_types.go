package types

import (
	"time"

	"github.com/go-dev-frame/sponge/pkg/sgorm/query"
)

var _ time.Time

// Tip: suggested filling in the binding rules https://github.com/go-playground/validator in request struct fields tag.

// CreateTradesRequest request params
type CreateTradesRequest struct {
	StrategyID        int        `json:"strategyID" binding:""`        // 关联策略ID（逻辑外键→strategies.id, 无物理约束）
	MarketID          int        `json:"marketID" binding:""`          // 关联市场ID（逻辑外键→markets.id, 无物理约束）
	PolymarketOrderID string     `json:"polymarketOrderID" binding:""` // Polymarket CLOB 订单ID（全局唯一）
	Side              string     `json:"side" binding:""`              // 交易方向: yes-Yes方, no-No方
	Action            string     `json:"action" binding:""`            // 操作类型: buy-买入建仓, sell-卖出平仓
	Amount            string     `json:"amount" binding:""`            // 交易金额（买入时为USDC支出, 卖出时为USDC收入）
	Price             string     `json:"price" binding:""`             // 成交单价
	Shares            string     `json:"shares" binding:""`            // 成交份额数量
	Status            string     `json:"status" binding:""`            // 订单状态: pending-待成交, filled-已成交, partial-部分成交, cancelled-已取消, failed-失败
	Fee               string     `json:"fee" binding:""`               // 交易手续费（USDC）
	Pnl               string     `json:"pnl" binding:""`               // 盈亏金额（平仓时有值, USDC）
	CloseReason       string     `json:"closeReason" binding:""`       // 平仓原因: take_profit-止盈, stop_loss-止损, pre_resolution-到期前, manual-手动
	FilledAt          *time.Time `json:"filledAt" binding:""`          // 订单成交时间
	ClosedAt          *time.Time `json:"closedAt" binding:""`          // 平仓时间
}

// UpdateTradesByIDRequest request params
type UpdateTradesByIDRequest struct {
	ID uint64 `json:"id" binding:""` // uint64 id
	// 主键ID
	StrategyID        int        `json:"strategyID" binding:""`        // 关联策略ID（逻辑外键→strategies.id, 无物理约束）
	MarketID          int        `json:"marketID" binding:""`          // 关联市场ID（逻辑外键→markets.id, 无物理约束）
	PolymarketOrderID string     `json:"polymarketOrderID" binding:""` // Polymarket CLOB 订单ID（全局唯一）
	Side              string     `json:"side" binding:""`              // 交易方向: yes-Yes方, no-No方
	Action            string     `json:"action" binding:""`            // 操作类型: buy-买入建仓, sell-卖出平仓
	Amount            string     `json:"amount" binding:""`            // 交易金额（买入时为USDC支出, 卖出时为USDC收入）
	Price             string     `json:"price" binding:""`             // 成交单价
	Shares            string     `json:"shares" binding:""`            // 成交份额数量
	Status            string     `json:"status" binding:""`            // 订单状态: pending-待成交, filled-已成交, partial-部分成交, cancelled-已取消, failed-失败
	Fee               string     `json:"fee" binding:""`               // 交易手续费（USDC）
	Pnl               string     `json:"pnl" binding:""`               // 盈亏金额（平仓时有值, USDC）
	CloseReason       string     `json:"closeReason" binding:""`       // 平仓原因: take_profit-止盈, stop_loss-止损, pre_resolution-到期前, manual-手动
	FilledAt          *time.Time `json:"filledAt" binding:""`          // 订单成交时间
	ClosedAt          *time.Time `json:"closedAt" binding:""`          // 平仓时间
}

// TradesObjDetail detail
type TradesObjDetail struct {
	ID uint64 `json:"id"` // convert to uint64 id
	// 主键ID
	CreatedAt         *time.Time `json:"createdAt"`         // 创建时间
	UpdatedAt         *time.Time `json:"updatedAt"`         // 更新时间
	StrategyID        int        `json:"strategyID"`        // 关联策略ID（逻辑外键→strategies.id, 无物理约束）
	MarketID          int        `json:"marketID"`          // 关联市场ID（逻辑外键→markets.id, 无物理约束）
	PolymarketOrderID string     `json:"polymarketOrderID"` // Polymarket CLOB 订单ID（全局唯一）
	Side              string     `json:"side"`              // 交易方向: yes-Yes方, no-No方
	Action            string     `json:"action"`            // 操作类型: buy-买入建仓, sell-卖出平仓
	Amount            string     `json:"amount"`            // 交易金额（买入时为USDC支出, 卖出时为USDC收入）
	Price             string     `json:"price"`             // 成交单价
	Shares            string     `json:"shares"`            // 成交份额数量
	Status            string     `json:"status"`            // 订单状态: pending-待成交, filled-已成交, partial-部分成交, cancelled-已取消, failed-失败
	Fee               string     `json:"fee"`               // 交易手续费（USDC）
	Pnl               string     `json:"pnl"`               // 盈亏金额（平仓时有值, USDC）
	CloseReason       string     `json:"closeReason"`       // 平仓原因: take_profit-止盈, stop_loss-止损, pre_resolution-到期前, manual-手动
	FilledAt          *time.Time `json:"filledAt"`          // 订单成交时间
	ClosedAt          *time.Time `json:"closedAt"`          // 平仓时间
}

// CreateTradesReply only for api docs
type CreateTradesReply struct {
	Code int    `json:"code"` // return code
	Msg  string `json:"msg"`  // return information description
	Data struct {
		ID uint64 `json:"id"` // id
	} `json:"data"` // return data
}

// DeleteTradesByIDReply only for api docs
type DeleteTradesByIDReply struct {
	Code int      `json:"code"` // return code
	Msg  string   `json:"msg"`  // return information description
	Data struct{} `json:"data"` // return data
}

// UpdateTradesByIDReply only for api docs
type UpdateTradesByIDReply struct {
	Code int      `json:"code"` // return code
	Msg  string   `json:"msg"`  // return information description
	Data struct{} `json:"data"` // return data
}

// GetTradesByIDReply only for api docs
type GetTradesByIDReply struct {
	Code int    `json:"code"` // return code
	Msg  string `json:"msg"`  // return information description
	Data struct {
		Trades TradesObjDetail `json:"trades"`
	} `json:"data"` // return data
}

// ListTradessRequest request params
type ListTradessRequest struct {
	query.Params
}

// ListTradessReply only for api docs
type ListTradessReply struct {
	Code int    `json:"code"` // return code
	Msg  string `json:"msg"`  // return information description
	Data struct {
		Tradess []TradesObjDetail `json:"tradess"`
	} `json:"data"` // return data
}
