package model

import (
	"github.com/go-dev-frame/sponge/pkg/sgorm"
	"time"
)

// User 用户表（通过 MetaMask 钱包登录的用户）
type User struct {
	sgorm.Model `gorm:"embedded"` // embed id and time

	WalletAddress string     `gorm:"column:wallet_address;type:varchar(256);uniqueIndex;not null" json:"walletAddress"` // MetaMask 钱包地址
	Nickname      string     `gorm:"column:nickname;type:varchar(128)" json:"nickname"`                                 // 昵称
	Avatar        string     `gorm:"column:avatar;type:varchar(512)" json:"avatar"`                                     // 头像 URL
	LastLoginAt   *time.Time `gorm:"column:last_login_at;type:datetime" json:"lastLoginAt"`                             // 最后登录时间
}

// TableName 指定表名
func (u *User) TableName() string {
	return "users"
}

// UserColumnNames Whitelist for custom query fields to prevent sql injection attacks
var UserColumnNames = map[string]bool{
	"id":             true,
	"created_at":     true,
	"updated_at":     true,
	"deleted_at":     true,
	"wallet_address": true,
	"nickname":       true,
	"avatar":         true,
	"last_login_at":  true,
}
