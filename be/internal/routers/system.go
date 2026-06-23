package routers

import (
	"github.com/gin-gonic/gin"

	"be/internal/handler"
)

func init() {
	// 系统路由无需 jwt 认证，注册为 public router
	apiV1PublicRouterFns = append(apiV1PublicRouterFns, func(group *gin.RouterGroup) {
		systemRouter(group, handler.NewSystemHandler())
	})
}

func systemRouter(group *gin.RouterGroup, h handler.SystemHandler) {
	g := group.Group("/system")
	g.GET("/status", h.GetStatus) // GET /api/v1/system/status
	g.POST("/pause", h.Pause)     // POST /api/v1/system/pause
	g.POST("/resume", h.Resume)   // POST /api/v1/system/resume
}
