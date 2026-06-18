package types

import (
	"time"

	"github.com/go-dev-frame/sponge/pkg/sgorm/query"
)

var _ time.Time

// Tip: suggested filling in the binding rules https://github.com/go-playground/validator in request struct fields tag.

// CreateVaultSnapshotsRequest request params
type CreateVaultSnapshotsRequest struct {
	TotalAssets    string     `json:"totalAssets" binding:""`    // 金库总资产（含链下策略债务, USDC）
	SharePrice     string     `json:"sharePrice" binding:""`     // 当前份额价格（USDC/份额）
	Tvl            string     `json:"tvl" binding:""`            // 锁定总价值（Total Value Locked, USDC）
	DepositorCount int        `json:"depositorCount" binding:""` // 存款人数量
	DeployedAmount string     `json:"deployedAmount" binding:""` // 已部署到链下策略的资金量（USDC）
	SnapshotAt     *time.Time `json:"snapshotAt" binding:""`     // 快照时间戳
}

// UpdateVaultSnapshotsByIDRequest request params
type UpdateVaultSnapshotsByIDRequest struct {
	ID uint64 `json:"id" binding:""` // uint64 id
	// 主键ID
	TotalAssets    string     `json:"totalAssets" binding:""`    // 金库总资产（含链下策略债务, USDC）
	SharePrice     string     `json:"sharePrice" binding:""`     // 当前份额价格（USDC/份额）
	Tvl            string     `json:"tvl" binding:""`            // 锁定总价值（Total Value Locked, USDC）
	DepositorCount int        `json:"depositorCount" binding:""` // 存款人数量
	DeployedAmount string     `json:"deployedAmount" binding:""` // 已部署到链下策略的资金量（USDC）
	SnapshotAt     *time.Time `json:"snapshotAt" binding:""`     // 快照时间戳
}

// VaultSnapshotsObjDetail detail
type VaultSnapshotsObjDetail struct {
	ID uint64 `json:"id"` // convert to uint64 id
	// 主键ID
	CreatedAt      *time.Time `json:"createdAt"`      // 创建时间
	UpdatedAt      *time.Time `json:"updatedAt"`      // 更新时间
	TotalAssets    string     `json:"totalAssets"`    // 金库总资产（含链下策略债务, USDC）
	SharePrice     string     `json:"sharePrice"`     // 当前份额价格（USDC/份额）
	Tvl            string     `json:"tvl"`            // 锁定总价值（Total Value Locked, USDC）
	DepositorCount int        `json:"depositorCount"` // 存款人数量
	DeployedAmount string     `json:"deployedAmount"` // 已部署到链下策略的资金量（USDC）
	SnapshotAt     *time.Time `json:"snapshotAt"`     // 快照时间戳
}

// CreateVaultSnapshotsReply only for api docs
type CreateVaultSnapshotsReply struct {
	Code int    `json:"code"` // return code
	Msg  string `json:"msg"`  // return information description
	Data struct {
		ID uint64 `json:"id"` // id
	} `json:"data"` // return data
}

// DeleteVaultSnapshotsByIDReply only for api docs
type DeleteVaultSnapshotsByIDReply struct {
	Code int      `json:"code"` // return code
	Msg  string   `json:"msg"`  // return information description
	Data struct{} `json:"data"` // return data
}

// UpdateVaultSnapshotsByIDReply only for api docs
type UpdateVaultSnapshotsByIDReply struct {
	Code int      `json:"code"` // return code
	Msg  string   `json:"msg"`  // return information description
	Data struct{} `json:"data"` // return data
}

// GetVaultSnapshotsByIDReply only for api docs
type GetVaultSnapshotsByIDReply struct {
	Code int    `json:"code"` // return code
	Msg  string `json:"msg"`  // return information description
	Data struct {
		VaultSnapshots VaultSnapshotsObjDetail `json:"vaultSnapshots"`
	} `json:"data"` // return data
}

// ListVaultSnapshotssRequest request params
type ListVaultSnapshotssRequest struct {
	query.Params
}

// ListVaultSnapshotssReply only for api docs
type ListVaultSnapshotssReply struct {
	Code int    `json:"code"` // return code
	Msg  string `json:"msg"`  // return information description
	Data struct {
		VaultSnapshotss []VaultSnapshotsObjDetail `json:"vaultSnapshotss"`
	} `json:"data"` // return data
}
