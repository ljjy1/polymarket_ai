package routers

import (
	"github.com/gin-gonic/gin"

	"be/internal/handler"
)

func init() {
	// 外部数据接口（公开，无需 jwt 认证）
	apiV1PublicRouterFns = append(apiV1PublicRouterFns, func(group *gin.RouterGroup) {
		externalRouter(group, handler.NewExternalHandler())
	})
}

func externalRouter(group *gin.RouterGroup, h handler.ExternalHandler) {
	g := group.Group("/external")

	// 恐惧贪婪指数
	g.GET("/fear-greed-index", h.GetFearGreedIndex)           // GET /api/v1/external/fear-greed-index
	g.GET("/fear-greed-index/trend", h.GetFearGreedWithTrend) // GET /api/v1/external/fear-greed-index/trend?days=7

	// 新闻
	g.GET("/news", h.GetNews) // GET /api/v1/external/news?symbol=BTC&count=5

	// 链上数据
	g.GET("/on-chain", h.GetOnChainData) // GET /api/v1/external/on-chain

	// 技术指标
	g.POST("/technical-indicators", h.GetTechnicalIndicators)      // POST /api/v1/external/technical-indicators
	g.POST("/technical-indicators/1h", h.GetTechnicalIndicators1h) // POST /api/v1/external/technical-indicators/1h
}
