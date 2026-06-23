package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/go-dev-frame/sponge/pkg/errcode"
	"github.com/go-dev-frame/sponge/pkg/gin/response"

	"be/internal/tasks"
)

// ---------- 请求体 ----------

type marketIDReq struct {
	MarketID int `json:"marketId" binding:"required"`
}

type strategyReq struct {
	PredictionID int `json:"predictionId" binding:"required"`
	MarketID     int `json:"marketId" binding:"required"`
}

type executeReq struct {
	MarketID   int `json:"marketId" binding:"required"`
	StrategyID int `json:"strategyId" binding:"required"`
}

// ---------- Handler ----------

type TaskHandler interface {
	// 日常任务
	TriggerScan(c *gin.Context)
	TriggerPredict(c *gin.Context)
	TriggerStrategy(c *gin.Context)
	TriggerExecute(c *gin.Context)
	// 监控任务
	TriggerMonitorPositions(c *gin.Context)
	TriggerVaultSnapshot(c *gin.Context)
	TriggerHealthCheck(c *gin.Context)
	TriggerSettlement(c *gin.Context)
}

type taskHandler struct{}

func NewTaskHandler() TaskHandler {
	return &taskHandler{}
}

// ---------- 日常任务 ----------

// TriggerScan  POST /api/v1/tasks/scan
// @Summary 手动触发市场扫描
// @Description 手动触发市场扫描任务，从 Polymarket 拉取最新市场数据
// @Tags tasks
// @Accept json
// @Produce json
// @Success 200 {object} types.Result
// @Router /api/v1/tasks/scan [post]
// @Security BearerAuth
func (h *taskHandler) TriggerScan(c *gin.Context) {
	if err := tasks.TriggerScan(c.Request.Context()); err != nil {
		response.Error(c, errcode.InternalServerError)
		return
	}
	response.Success(c, gin.H{"message": "市场扫描任务执行完成"})
}

// TriggerPredict  POST /api/v1/tasks/predict  需要 marketId
// @Summary 手动触发 AI 预测
// @Description 手动触发 AI 预测任务，对指定市场进行预测分析
// @Tags tasks
// @Accept json
// @Produce json
// @Param data body marketIDReq true "市场ID"
// @Success 200 {object} types.Result
// @Router /api/v1/tasks/predict [post]
// @Security BearerAuth
func (h *taskHandler) TriggerPredict(c *gin.Context) {
	var req marketIDReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, errcode.InvalidParams)
		return
	}
	if err := tasks.TriggerPredict(c.Request.Context(), req.MarketID); err != nil {
		response.Error(c, errcode.InternalServerError)
		return
	}
	response.Success(c, gin.H{"message": "AI 预测任务执行完成"})
}

// TriggerStrategy  POST /api/v1/tasks/strategy  需要 predictionId, marketId
// @Summary 手动触发策略生成
// @Description 手动触发策略生成任务，基于预测结果生成交易策略
// @Tags tasks
// @Accept json
// @Produce json
// @Param data body strategyReq true "预测ID和市场ID"
// @Success 200 {object} types.Result
// @Router /api/v1/tasks/strategy [post]
// @Security BearerAuth
func (h *taskHandler) TriggerStrategy(c *gin.Context) {
	var req strategyReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, errcode.InvalidParams)
		return
	}
	if err := tasks.TriggerStrategy(c.Request.Context(), req.PredictionID, req.MarketID); err != nil {
		response.Error(c, errcode.InternalServerError)
		return
	}
	response.Success(c, gin.H{"message": "策略生成任务执行完成"})
}

// TriggerExecute  POST /api/v1/tasks/execute  需要 marketId, strategyId
// @Summary 手动触发交易执行
// @Description 手动触发交易执行任务，根据策略执行实际交易
// @Tags tasks
// @Accept json
// @Produce json
// @Param data body executeReq true "市场ID和策略ID"
// @Success 200 {object} types.Result
// @Router /api/v1/tasks/execute [post]
// @Security BearerAuth
func (h *taskHandler) TriggerExecute(c *gin.Context) {
	var req executeReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, errcode.InvalidParams)
		return
	}
	if err := tasks.TriggerExecute(c.Request.Context(), req.MarketID, req.StrategyID); err != nil {
		response.Error(c, errcode.InternalServerError)
		return
	}
	response.Success(c, gin.H{"message": "交易执行任务执行完成"})
}

// ---------- 监控任务 ----------

// TriggerMonitorPositions  POST /api/v1/tasks/positions
// @Summary 手动触发持仓监控
// @Description 手动触发持仓监控任务，检查当前持仓状态
// @Tags tasks
// @Accept json
// @Produce json
// @Success 200 {object} types.Result
// @Router /api/v1/tasks/positions [post]
// @Security BearerAuth
func (h *taskHandler) TriggerMonitorPositions(c *gin.Context) {
	if err := tasks.TriggerMonitorPositions(c.Request.Context()); err != nil {
		response.Error(c, errcode.InternalServerError)
		return
	}
	response.Success(c, gin.H{"message": "持仓监控任务执行完成"})
}

// TriggerVaultSnapshot  POST /api/v1/tasks/vault-snapshot
// @Summary 手动触发金库快照
// @Description 手动触发金库快照任务，记录当前金库状态
// @Tags tasks
// @Accept json
// @Produce json
// @Success 200 {object} types.Result
// @Router /api/v1/tasks/vault-snapshot [post]
// @Security BearerAuth
func (h *taskHandler) TriggerVaultSnapshot(c *gin.Context) {
	if err := tasks.TriggerVaultSnapshot(c.Request.Context()); err != nil {
		response.Error(c, errcode.InternalServerError)
		return
	}
	response.Success(c, gin.H{"message": "金库快照任务执行完成"})
}

// TriggerHealthCheck  POST /api/v1/tasks/health-check
// @Summary 手动触发健康检查
// @Description 手动触发健康检查任务，检查系统各组件运行状态
// @Tags tasks
// @Accept json
// @Produce json
// @Success 200 {object} types.Result
// @Router /api/v1/tasks/health-check [post]
// @Security BearerAuth
func (h *taskHandler) TriggerHealthCheck(c *gin.Context) {
	if err := tasks.TriggerHealthCheck(c.Request.Context()); err != nil {
		response.Error(c, errcode.InternalServerError)
		return
	}
	response.Success(c, gin.H{"message": "健康检查任务执行完成"})
}

// TriggerSettlement  POST /api/v1/tasks/settlement
// @Summary 手动触发结算检查
// @Description 手动触发结算检查任务，检查已过期市场的结算状态
// @Tags tasks
// @Accept json
// @Produce json
// @Success 200 {object} types.Result
// @Router /api/v1/tasks/settlement [post]
// @Security BearerAuth
func (h *taskHandler) TriggerSettlement(c *gin.Context) {
	if err := tasks.TriggerSettlement(c.Request.Context()); err != nil {
		response.Error(c, errcode.InternalServerError)
		return
	}
	response.Success(c, gin.H{"message": "结算检查任务执行完成"})
}
