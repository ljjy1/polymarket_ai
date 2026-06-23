package model

import (
	"github.com/go-dev-frame/sponge/pkg/sgorm"
	"gorm.io/datatypes"
)

// SystemLogs 系统日志表（同时用作暂停恢复标志位存储）
type SystemLogs struct {
	sgorm.Model `gorm:"embedded"` // embed id and time

	Level   string          `gorm:"column:level;type:varchar(16);not null;comment:日志级别: INFO, WARNING, ERROR, DEBUG" json:"level"`
	Source  string          `gorm:"column:source;type:varchar(64);not null;comment:日志来源（模块名, 如 scanner, predictor, executor）" json:"source"`
	Message string          `gorm:"column:message;type:text;not null;comment:日志消息内容" json:"message"`
	Context *datatypes.JSON `gorm:"column:context;type:json;not null;comment:日志上下文信息（JSON格式, 包含额外结构化数据）" json:"context"`
	TraceID string          `gorm:"column:trace_id;type:varchar(64);comment:链路追踪ID（用于关联同一请求链中的多条日志）" json:"traceID"`
}

// SystemLogsColumnNames Whitelist for custom query fields to prevent sql injection attacks
var SystemLogsColumnNames = map[string]bool{
	"id":         true,
	"created_at": true,
	"updated_at": true,
	"deleted_at": true,
	"level":      true,
	"source":     true,
	"message":    true,
	"context":    true,
	"trace_id":   true,
}
