package types

import (
	"time"

	"github.com/go-dev-frame/sponge/pkg/sgorm/query"
)

var _ time.Time

// Tip: suggested filling in the binding rules https://github.com/go-playground/validator in request struct fields tag.

// CreateSystemLogsRequest request params
type CreateSystemLogsRequest struct {
	Level   string `json:"level" binding:""`   // 日志级别: INFO, WARNING, ERROR, DEBUG
	Source  string `json:"source" binding:""`  // 日志来源（模块名, 如 scanner, predictor, executor）
	Message string `json:"message" binding:""` // 日志消息内容
	Context string `json:"context" binding:""` // 日志上下文信息（JSON格式, 包含额外结构化数据）
	TraceID string `json:"traceID" binding:""` // 链路追踪ID（用于关联同一请求链中的多条日志）
}

// UpdateSystemLogsByIDRequest request params
type UpdateSystemLogsByIDRequest struct {
	ID uint64 `json:"id" binding:""` // uint64 id
	// 主键ID
	Level   string `json:"level" binding:""`   // 日志级别: INFO, WARNING, ERROR, DEBUG
	Source  string `json:"source" binding:""`  // 日志来源（模块名, 如 scanner, predictor, executor）
	Message string `json:"message" binding:""` // 日志消息内容
	Context string `json:"context" binding:""` // 日志上下文信息（JSON格式, 包含额外结构化数据）
	TraceID string `json:"traceID" binding:""` // 链路追踪ID（用于关联同一请求链中的多条日志）
}

// SystemLogsObjDetail detail
type SystemLogsObjDetail struct {
	ID uint64 `json:"id"` // convert to uint64 id
	// 主键ID
	CreatedAt *time.Time `json:"createdAt"` // 创建时间
	UpdatedAt *time.Time `json:"updatedAt"` // 更新时间
	Level     string     `json:"level"`     // 日志级别: INFO, WARNING, ERROR, DEBUG
	Source    string     `json:"source"`    // 日志来源（模块名, 如 scanner, predictor, executor）
	Message   string     `json:"message"`   // 日志消息内容
	Context   string     `json:"context"`   // 日志上下文信息（JSON格式, 包含额外结构化数据）
	TraceID   string     `json:"traceID"`   // 链路追踪ID（用于关联同一请求链中的多条日志）
}

// CreateSystemLogsReply only for api docs
type CreateSystemLogsReply struct {
	Code int    `json:"code"` // return code
	Msg  string `json:"msg"`  // return information description
	Data struct {
		ID uint64 `json:"id"` // id
	} `json:"data"` // return data
}

// DeleteSystemLogsByIDReply only for api docs
type DeleteSystemLogsByIDReply struct {
	Code int      `json:"code"` // return code
	Msg  string   `json:"msg"`  // return information description
	Data struct{} `json:"data"` // return data
}

// UpdateSystemLogsByIDReply only for api docs
type UpdateSystemLogsByIDReply struct {
	Code int      `json:"code"` // return code
	Msg  string   `json:"msg"`  // return information description
	Data struct{} `json:"data"` // return data
}

// GetSystemLogsByIDReply only for api docs
type GetSystemLogsByIDReply struct {
	Code int    `json:"code"` // return code
	Msg  string `json:"msg"`  // return information description
	Data struct {
		SystemLogs SystemLogsObjDetail `json:"systemLogs"`
	} `json:"data"` // return data
}

// ListSystemLogssRequest request params
type ListSystemLogssRequest struct {
	query.Params
}

// ListSystemLogssReply only for api docs
type ListSystemLogssReply struct {
	Code int    `json:"code"` // return code
	Msg  string `json:"msg"`  // return information description
	Data struct {
		SystemLogss []SystemLogsObjDetail `json:"systemLogss"`
	} `json:"data"` // return data
}
