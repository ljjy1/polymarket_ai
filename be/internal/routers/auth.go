package routers

import (
	"github.com/gin-gonic/gin"

	"be/internal/handler"
)

func init() {
	apiV1PublicRouterFns = append(apiV1PublicRouterFns, func(group *gin.RouterGroup) {
		authRouter(group, handler.NewAuthHandler())
	})
}

func authRouter(group *gin.RouterGroup, h handler.AuthHandler) {
	g := group.Group("/auth")

	g.POST("/nonce", h.Nonce) // [post] /api/v1/auth/nonce
	g.POST("/login", h.Login) // [post] /api/v1/auth/login
}
