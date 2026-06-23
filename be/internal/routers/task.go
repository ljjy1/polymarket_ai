package routers

import (
	"github.com/gin-gonic/gin"

	"be/internal/handler"
)

func init() {
	// 任务手动触发接口（需要 jwt 认证）
	apiV1RouterFns = append(apiV1RouterFns, func(group *gin.RouterGroup) {
		taskRouter(group, handler.NewTaskHandler())
	})
}

func taskRouter(group *gin.RouterGroup, h handler.TaskHandler) {
	g := group.Group("/tasks")

	// 日常任务
	g.POST("/scan", h.TriggerScan)         // POST /api/v1/tasks/scan
	g.POST("/predict", h.TriggerPredict)   // POST /api/v1/tasks/predict
	g.POST("/strategy", h.TriggerStrategy) // POST /api/v1/tasks/strategy
	g.POST("/execute", h.TriggerExecute)   // POST /api/v1/tasks/execute

	// 监控任务
	g.POST("/positions", h.TriggerMonitorPositions)   // POST /api/v1/tasks/positions
	g.POST("/vault-snapshot", h.TriggerVaultSnapshot) // POST /api/v1/tasks/vault-snapshot
	g.POST("/health-check", h.TriggerHealthCheck)     // POST /api/v1/tasks/health-check
	g.POST("/settlement", h.TriggerSettlement)        // POST /api/v1/tasks/settlement
}
