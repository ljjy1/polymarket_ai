package types

import (
	"time"

	"github.com/go-dev-frame/sponge/pkg/sgorm/query"
)

var _ time.Time

// Tip: suggested filling in the binding rules https://github.com/go-playground/validator in request struct fields tag.

// CreateMarketsRequest request params
type CreateMarketsRequest struct {
	PolymarketConditionID string     `json:"polymarketConditionID" binding:""` // Polymarket 条件ID（唯一标识一个预测市场）
	PolymarketTokenID     string     `json:"polymarketTokenID" binding:""`     // Polymarket Yes Token ID（AI默认下注方向）
	EventSlug             string     `json:"eventSlug" binding:""`             // 事件唯一标识符（Event Slug）
	Question              string     `json:"question" binding:""`              // 预测市场问题标题（如"BTC年底>$100K?"）
	PriceThreshold        int        `json:"priceThreshold" binding:""`        // 价格阈值（用于筛选市场的价格门槛, 单位: 百分点）
	ScanDate              *time.Time `json:"scanDate" binding:""`              // 扫描日期（幂等键，每日最多一条记录）
	TargetDate            *time.Time `json:"targetDate" binding:""`            // 市场预测目标日期（即 Polymarket 的到期日）
	CurrentYesPrice       string     `json:"currentYesPrice" binding:""`       // 当前 Yes 代币价格（即市场概率）
	CurrentNoPrice        string     `json:"currentNoPrice" binding:""`        // 当前 No 代币价格
	SelectedAt            *time.Time `json:"selectedAt" binding:""`            // 被策略选中的时间戳
	Status                string     `json:"status" binding:""`                // 市场状态: active-活跃, resolved-已结算, expired-已过期
	Resolution            string     `json:"resolution" binding:""`            // 结算结果: yes-是, no-否（仅在 resolved 后有值）
}

// UpdateMarketsByIDRequest request params
type UpdateMarketsByIDRequest struct {
	ID uint64 `json:"id" binding:""` // uint64 id
	// 主键ID
	PolymarketConditionID string     `json:"polymarketConditionID" binding:""` // Polymarket 条件ID（唯一标识一个预测市场）
	PolymarketTokenID     string     `json:"polymarketTokenID" binding:""`     // Polymarket Yes Token ID（AI默认下注方向）
	EventSlug             string     `json:"eventSlug" binding:""`             // 事件唯一标识符（Event Slug）
	Question              string     `json:"question" binding:""`              // 预测市场问题标题（如"BTC年底>$100K?"）
	PriceThreshold        int        `json:"priceThreshold" binding:""`        // 价格阈值（用于筛选市场的价格门槛, 单位: 百分点）
	ScanDate              *time.Time `json:"scanDate" binding:""`              // 扫描日期（幂等键，每日最多一条记录）
	TargetDate            *time.Time `json:"targetDate" binding:""`            // 市场预测目标日期（即 Polymarket 的到期日）
	CurrentYesPrice       string     `json:"currentYesPrice" binding:""`       // 当前 Yes 代币价格（即市场概率）
	CurrentNoPrice        string     `json:"currentNoPrice" binding:""`        // 当前 No 代币价格
	SelectedAt            *time.Time `json:"selectedAt" binding:""`            // 被策略选中的时间戳
	Status                string     `json:"status" binding:""`                // 市场状态: active-活跃, resolved-已结算, expired-已过期
	Resolution            string     `json:"resolution" binding:""`            // 结算结果: yes-是, no-否（仅在 resolved 后有值）
}

// MarketsObjDetail detail
type MarketsObjDetail struct {
	ID uint64 `json:"id"` // convert to uint64 id
	// 主键ID
	CreatedAt             *time.Time `json:"createdAt"`             // 创建时间
	UpdatedAt             *time.Time `json:"updatedAt"`             // 更新时间
	PolymarketConditionID string     `json:"polymarketConditionID"` // Polymarket 条件ID（唯一标识一个预测市场）
	PolymarketTokenID     string     `json:"polymarketTokenID"`     // Polymarket Yes Token ID（AI默认下注方向）
	EventSlug             string     `json:"eventSlug"`             // 事件唯一标识符（Event Slug）
	Question              string     `json:"question"`              // 预测市场问题标题（如"BTC年底>$100K?"）
	PriceThreshold        int        `json:"priceThreshold"`        // 价格阈值（用于筛选市场的价格门槛, 单位: 百分点）
	ScanDate              *time.Time `json:"scanDate"`              // 扫描日期（幂等键，每日最多一条记录）
	TargetDate            *time.Time `json:"targetDate"`            // 市场预测目标日期（即 Polymarket 的到期日）
	CurrentYesPrice       string     `json:"currentYesPrice"`       // 当前 Yes 代币价格（即市场概率）
	CurrentNoPrice        string     `json:"currentNoPrice"`        // 当前 No 代币价格
	SelectedAt            *time.Time `json:"selectedAt"`            // 被策略选中的时间戳
	Status                string     `json:"status"`                // 市场状态: active-活跃, resolved-已结算, expired-已过期
	Resolution            string     `json:"resolution"`            // 结算结果: yes-是, no-否（仅在 resolved 后有值）
}

// CreateMarketsReply only for api docs
type CreateMarketsReply struct {
	Code int    `json:"code"` // return code
	Msg  string `json:"msg"`  // return information description
	Data struct {
		ID uint64 `json:"id"` // id
	} `json:"data"` // return data
}

// DeleteMarketsByIDReply only for api docs
type DeleteMarketsByIDReply struct {
	Code int      `json:"code"` // return code
	Msg  string   `json:"msg"`  // return information description
	Data struct{} `json:"data"` // return data
}

// UpdateMarketsByIDReply only for api docs
type UpdateMarketsByIDReply struct {
	Code int      `json:"code"` // return code
	Msg  string   `json:"msg"`  // return information description
	Data struct{} `json:"data"` // return data
}

// GetMarketsByIDReply only for api docs
type GetMarketsByIDReply struct {
	Code int    `json:"code"` // return code
	Msg  string `json:"msg"`  // return information description
	Data struct {
		Markets MarketsObjDetail `json:"markets"`
	} `json:"data"` // return data
}

// ListMarketssRequest request params
type ListMarketssRequest struct {
	query.Params
}

// ListMarketssReply only for api docs
type ListMarketssReply struct {
	Code int    `json:"code"` // return code
	Msg  string `json:"msg"`  // return information description
	Data struct {
		Marketss []MarketsObjDetail `json:"marketss"`
	} `json:"data"` // return data
}
