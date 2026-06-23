package model

import (
	"github.com/go-dev-frame/sponge/pkg/sgorm"
	"github.com/shopspring/decimal"
	"time"
)

// VaultSnapshots PolyVault金库快照表
type VaultSnapshots struct {
	sgorm.Model `gorm:"embedded"` // embed id and time

	TotalAssets    *decimal.Decimal `gorm:"column:total_assets;type:decimal(38,18);not null;comment:金库总资产（含链下策略债务, USDC）" json:"totalAssets"`
	SharePrice     *decimal.Decimal `gorm:"column:share_price;type:decimal(38,18);not null;comment:当前份额价格（USDC/份额）" json:"sharePrice"`
	Tvl            *decimal.Decimal `gorm:"column:tvl;type:decimal(38,18);not null;comment:锁定总价值（Total Value Locked, USDC）" json:"tvl"`
	DepositorCount int              `gorm:"column:depositor_count;type:int(11);default:0;not null;comment:存款人数量" json:"depositorCount"`
	DeployedAmount *decimal.Decimal `gorm:"column:deployed_amount;type:decimal(38,18);default:0;not null;comment:已部署到链下策略的资金量（USDC）" json:"deployedAmount"`
	SnapshotAt     *time.Time       `gorm:"column:snapshot_at;type:datetime;not null;comment:快照时间戳" json:"snapshotAt"`
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
