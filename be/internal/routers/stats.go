package routers

import (
	"github.com/gin-gonic/gin"

	"be/internal/handler"
)

func init() {
	apiV1RouterFns = append(apiV1RouterFns, func(group *gin.RouterGroup) {
		statsRouter(group, handler.NewStatsHandler())
	})
}

func statsRouter(group *gin.RouterGroup, h handler.StatsHandler) {
	g := group.Group("/stats")
	g.GET("/overview", h.Overview) // GET /api/v1/stats/overview
	g.GET("/daily", h.Daily)       // GET /api/v1/stats/daily
}
