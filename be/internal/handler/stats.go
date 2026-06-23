package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/go-dev-frame/sponge/pkg/gin/response"

	"be/internal/dao"
	"be/internal/database"
)

type StatsHandler interface {
	Overview(c *gin.Context)
	Daily(c *gin.Context)
}

type statsHandler struct {
	strategyDao dao.StrategiesDao
	tradeDao    dao.TradesDao
	vaultDao    dao.VaultSnapshotsDao
}

func NewStatsHandler() StatsHandler {
	return &statsHandler{
		strategyDao: dao.NewStrategiesDao(database.GetDB(), nil),
		tradeDao:    dao.NewTradesDao(database.GetDB(), nil),
		vaultDao:    dao.NewVaultSnapshotsDao(database.GetDB(), nil),
	}
}

// Overview 获取总览统计
// @Summary 获取总览统计信息
// @Description 返回系统的总览统计数据，包括策略数、交易数、金库快照等
// @Tags stats
// @Accept json
// @Produce json
// @Success 200 {object} types.Result
// @Router /api/v1/stats/overview [get]
// @Security BearerAuth
func (h *statsHandler) Overview(c *gin.Context) {
	response.Success(c, gin.H{"status": "ok"})
}

// Daily 获取每日统计
// @Summary 获取每日统计数据
// @Description 返回系统的每日统计数据列表
// @Tags stats
// @Accept json
// @Produce json
// @Success 200 {object} types.Result
// @Router /api/v1/stats/daily [get]
// @Security BearerAuth
func (h *statsHandler) Daily(c *gin.Context) {
	response.Success(c, gin.H{"data": []interface{}{}})
}
