package types

import (
	"time"

	"github.com/go-dev-frame/sponge/pkg/sgorm/query"
)

var _ time.Time

// Tip: suggested filling in the binding rules https://github.com/go-playground/validator in request struct fields tag.

// CreateStrategiesRequest request params
type CreateStrategiesRequest struct {
	PredictionID  int        `json:"predictionID" binding:""`  // 关联预测记录ID（逻辑外键→predictions.id, 无物理约束）
	MarketID      int        `json:"marketID" binding:""`      // 关联市场ID（逻辑外键→markets.id, 无物理约束）
	Action        string     `json:"action" binding:""`        // 策略动作: buy_yes-买入Yes, buy_no-买入No, skip-跳过
	Side          string     `json:"side" binding:""`          // 交易方向: yes-Yes方, no-No方（skip时为NULL）
	PositionSize  string     `json:"positionSize" binding:""`  // 仓位大小（USDC金额, skip时为0）
	EntryPrice    string     `json:"entryPrice" binding:""`    // 入场价格（skip时为0）
	TakeProfit    string     `json:"takeProfit" binding:""`    // 止盈价格（skip时为0）
	StopLoss      string     `json:"stopLoss" binding:""`      // 止损价格（skip时为0）
	KellyFraction string     `json:"kellyFraction" binding:""` // 凯利公式建议仓位比例（skip时为0）
	Edge          string     `json:"edge" binding:""`          // 策略采用的价差（正负值，有符号）
	SkipReason    string     `json:"skipReason" binding:""`    // 跳过交易的详细原因（未跳过时为空字符串）
	Status        string     `json:"status" binding:""`        // 策略状态: skipped-跳过, pending-待执行, executing-执行中, active-已开仓, closed-已平仓, failed-执行失败
	ExecutedAt    *time.Time `json:"executedAt" binding:""`    // 策略实际执行时间
}

// UpdateStrategiesByIDRequest request params
type UpdateStrategiesByIDRequest struct {
	ID uint64 `json:"id" binding:""` // uint64 id
	// 主键ID
	PredictionID  int        `json:"predictionID" binding:""`  // 关联预测记录ID（逻辑外键→predictions.id, 无物理约束）
	MarketID      int        `json:"marketID" binding:""`      // 关联市场ID（逻辑外键→markets.id, 无物理约束）
	Action        string     `json:"action" binding:""`        // 策略动作: buy_yes-买入Yes, buy_no-买入No, skip-跳过
	Side          string     `json:"side" binding:""`          // 交易方向: yes-Yes方, no-No方（skip时为NULL）
	PositionSize  string     `json:"positionSize" binding:""`  // 仓位大小（USDC金额, skip时为0）
	EntryPrice    string     `json:"entryPrice" binding:""`    // 入场价格（skip时为0）
	TakeProfit    string     `json:"takeProfit" binding:""`    // 止盈价格（skip时为0）
	StopLoss      string     `json:"stopLoss" binding:""`      // 止损价格（skip时为0）
	KellyFraction string     `json:"kellyFraction" binding:""` // 凯利公式建议仓位比例（skip时为0）
	Edge          string     `json:"edge" binding:""`          // 策略采用的价差（正负值，有符号）
	SkipReason    string     `json:"skipReason" binding:""`    // 跳过交易的详细原因（未跳过时为空字符串）
	Status        string     `json:"status" binding:""`        // 策略状态: skipped-跳过, pending-待执行, executing-执行中, active-已开仓, closed-已平仓, failed-执行失败
	ExecutedAt    *time.Time `json:"executedAt" binding:""`    // 策略实际执行时间
}

// StrategiesObjDetail detail
type StrategiesObjDetail struct {
	ID uint64 `json:"id"` // convert to uint64 id
	// 主键ID
	CreatedAt     *time.Time `json:"createdAt"`     // 创建时间
	UpdatedAt     *time.Time `json:"updatedAt"`     // 更新时间
	PredictionID  int        `json:"predictionID"`  // 关联预测记录ID（逻辑外键→predictions.id, 无物理约束）
	MarketID      int        `json:"marketID"`      // 关联市场ID（逻辑外键→markets.id, 无物理约束）
	Action        string     `json:"action"`        // 策略动作: buy_yes-买入Yes, buy_no-买入No, skip-跳过
	Side          string     `json:"side"`          // 交易方向: yes-Yes方, no-No方（skip时为NULL）
	PositionSize  string     `json:"positionSize"`  // 仓位大小（USDC金额, skip时为0）
	EntryPrice    string     `json:"entryPrice"`    // 入场价格（skip时为0）
	TakeProfit    string     `json:"takeProfit"`    // 止盈价格（skip时为0）
	StopLoss      string     `json:"stopLoss"`      // 止损价格（skip时为0）
	KellyFraction string     `json:"kellyFraction"` // 凯利公式建议仓位比例（skip时为0）
	Edge          string     `json:"edge"`          // 策略采用的价差（正负值，有符号）
	SkipReason    string     `json:"skipReason"`    // 跳过交易的详细原因（未跳过时为空字符串）
	Status        string     `json:"status"`        // 策略状态: skipped-跳过, pending-待执行, executing-执行中, active-已开仓, closed-已平仓, failed-执行失败
	ExecutedAt    *time.Time `json:"executedAt"`    // 策略实际执行时间
}

// CreateStrategiesReply only for api docs
type CreateStrategiesReply struct {
	Code int    `json:"code"` // return code
	Msg  string `json:"msg"`  // return information description
	Data struct {
		ID uint64 `json:"id"` // id
	} `json:"data"` // return data
}

// DeleteStrategiesByIDReply only for api docs
type DeleteStrategiesByIDReply struct {
	Code int      `json:"code"` // return code
	Msg  string   `json:"msg"`  // return information description
	Data struct{} `json:"data"` // return data
}

// UpdateStrategiesByIDReply only for api docs
type UpdateStrategiesByIDReply struct {
	Code int      `json:"code"` // return code
	Msg  string   `json:"msg"`  // return information description
	Data struct{} `json:"data"` // return data
}

// GetStrategiesByIDReply only for api docs
type GetStrategiesByIDReply struct {
	Code int    `json:"code"` // return code
	Msg  string `json:"msg"`  // return information description
	Data struct {
		Strategies StrategiesObjDetail `json:"strategies"`
	} `json:"data"` // return data
}

// ListStrategiessRequest request params
type ListStrategiessRequest struct {
	query.Params
}

// ListStrategiessReply only for api docs
type ListStrategiessReply struct {
	Code int    `json:"code"` // return code
	Msg  string `json:"msg"`  // return information description
	Data struct {
		Strategiess []StrategiesObjDetail `json:"strategiess"`
	} `json:"data"` // return data
}
