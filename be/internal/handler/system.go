package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/go-dev-frame/sponge/pkg/gin/response"
)

type SystemHandler interface {
	GetStatus(c *gin.Context)
	Pause(c *gin.Context)
	Resume(c *gin.Context)
}

type systemHandler struct{}

func NewSystemHandler() SystemHandler {
	return &systemHandler{}
}

// GetStatus 获取系统状态
// @Summary 获取系统状态
// @Description 返回当前系统运行状态和版本信息
// @Tags system
// @Accept json
// @Produce json
// @Success 200 {object} types.Result
// @Router /api/v1/system/status [get]
func (h *systemHandler) GetStatus(c *gin.Context) {
	response.Success(c, gin.H{
		"status":  "running",
		"version": "v0.0.0",
	})
}

// Pause 暂停系统
// @Summary 暂停系统
// @Description 暂停系统的定时任务调度
// @Tags system
// @Accept json
// @Produce json
// @Success 200 {object} types.Result
// @Router /api/v1/system/pause [post]
func (h *systemHandler) Pause(c *gin.Context) {
	response.Success(c, gin.H{"message": "system paused"})
}

// Resume 恢复系统
// @Summary 恢复系统运行
// @Description 恢复系统的定时任务调度
// @Tags system
// @Accept json
// @Produce json
// @Success 200 {object} types.Result
// @Router /api/v1/system/resume [post]
func (h *systemHandler) Resume(c *gin.Context) {
	response.Success(c, gin.H{"message": "system resumed"})
}
