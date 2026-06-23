package model

import (
	"github.com/go-dev-frame/sponge/pkg/sgorm"
	"github.com/shopspring/decimal"
	"time"
)

// Markets Polymarket 市场信息表
type Markets struct {
	sgorm.Model `gorm:"embedded"` // embed id and time

	PolymarketConditionID string           `gorm:"column:polymarket_condition_id;type:varchar(128);not null;comment:Polymarket 条件ID（唯一标识一个预测市场）" json:"polymarketConditionID"`
	PolymarketTokenID     string           `gorm:"column:polymarket_token_id;type:varchar(128);not null;comment:Polymarket Yes Token ID（AI默认下注方向）" json:"polymarketTokenID"`
	EventSlug             string           `gorm:"column:event_slug;type:varchar(255);not null;comment:事件唯一标识符（Event Slug）" json:"eventSlug"`
	Question              string           `gorm:"column:question;type:varchar(500);not null;comment:预测市场问题标题（如\"BTC年底>$100K?\"）" json:"question"`
	PriceThreshold        int              `gorm:"column:price_threshold;type:int(11);not null;comment:价格阈值（用于筛选市场的价格门槛, 单位: 百分点）" json:"priceThreshold"`
	ScanDate              *time.Time       `gorm:"column:scan_date;type:date;not null;comment:扫描日期（幂等键，每日最多一条记录）" json:"scanDate"`
	TargetDate            *time.Time       `gorm:"column:target_date;type:date;not null;comment:市场预测目标日期（即 Polymarket 的到期日）" json:"targetDate"`
	CurrentYesPrice       *decimal.Decimal `gorm:"column:current_yes_price;type:decimal(38,18);not null;comment:当前 Yes 代币价格（即市场概率）" json:"currentYesPrice"`
	CurrentNoPrice        *decimal.Decimal `gorm:"column:current_no_price;type:decimal(38,18);not null;comment:当前 No 代币价格" json:"currentNoPrice"`
	SelectedAt            *time.Time       `gorm:"column:selected_at;type:datetime;not null;comment:被策略选中的时间戳" json:"selectedAt"`
	Status                string           `gorm:"column:status;type:varchar(16);default:active;not null;comment:市场状态: active-活跃, resolved-已结算, expired-已过期" json:"status"`
	Resolution            string           `gorm:"column:resolution;type:varchar(8);comment:结算结果: yes-是, no-否（仅在 resolved 后有值）" json:"resolution"`
}

// MarketsColumnNames Whitelist for custom query fields to prevent sql injection attacks
var MarketsColumnNames = map[string]bool{
	"id":                      true,
	"created_at":              true,
	"updated_at":              true,
	"deleted_at":              true,
	"polymarket_condition_id": true,
	"polymarket_token_id":     true,
	"event_slug":              true,
	"question":                true,
	"price_threshold":         true,
	"scan_date":               true,
	"target_date":             true,
	"current_yes_price":       true,
	"current_no_price":        true,
	"selected_at":             true,
	"status":                  true,
	"resolution":              true,
}
