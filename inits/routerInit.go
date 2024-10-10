package inits

import (
	"forum/pkg/globals"
	"github.com/gin-gonic/gin"
)

// RouterInit
// @Description: 初始化路由
// @Author tianjiajie 2024-10-05 16:00:37
func RouterInit() {
	// 初始化路由 Router
	globals.Router = gin.Default()

	// 某一个控制器报错，不影响整体
	globals.Router.Use(gin.Recovery())
}