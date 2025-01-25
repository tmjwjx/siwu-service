package routes

import (
	"forum/internal/backstage/controllers"
	"forum/pkg/token"
	"github.com/gin-gonic/gin"
)

// Backstage
// @Description: 后台路由
// @Author lizhuang 2024-10-21 21:21:39
// @param        e *gin.Engine
func Backstage(e *gin.Engine) {
	// 后台登陆
	e.POST("/backstage/login", controllers.BsLogin)

	// 分组
	r := e.Group("/backstage")
	// token 校验
	r.Use(token.AuthMiddleware())

	// 后台登出
	r.POST("/logout", controllers.BsLogout)
}
