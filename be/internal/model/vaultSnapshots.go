package model

import (
	"github.com/go-dev-frame/sponge/pkg/sgorm"
	"github.com/shopspring/decimal"
	"time"
)

// VaultSnapshots PolyVault金库快照表
type VaultSnapshots struct {
	sgorm.Model `gorm:"embedded"` // embed id and time

	TotalAssets    *decimal.Decimal `gorm:"column:total_assets;type:decimal(38,18);not null" json:"totalAssets"`                                    // 金库总资产（含链下策略债务, USDC）
	SharePrice     *decimal.Decimal `gorm:"column:share_price;type:decimal(38,18);not null" json:"sharePrice"`                                      // 当前份额价格（USDC/份额）
	Tvl            *decimal.Decimal `gorm:"column:tvl;type:decimal(38,18);not null" json:"tvl"`                                                     // 锁定总价值（Total Value Locked, USDC）
	DepositorCount int              `gorm:"column:depositor_count;type:int(11);default:0;not null" json:"depositorCount"`                           // 存款人数量
	DeployedAmount *decimal.Decimal `gorm:"column:deployed_amount;type:decimal(38,18);default:0.000000000000000000;not null" json:"deployedAmount"` // 已部署到链下策略的资金量（USDC）
	SnapshotAt     *time.Time       `gorm:"column:snapshot_at;type:datetime;not null" json:"snapshotAt"`                                            // 快照时间戳
}

// VaultSnapshotsColumnNames Whitelist for custom query fields to prevent sql injection attacks
var VaultSnapshotsColumnNames = map[string]bool{
	"id":              true,
	"created_at":      true,
	"updated_at":      true,
	"deleted_at":      true,
	"total_assets":    true,
	"share_price":     true,
	"tvl":             true,
	"depositor_count": true,
	"deployed_amount": true,
	"snapshot_at":     true,
}
