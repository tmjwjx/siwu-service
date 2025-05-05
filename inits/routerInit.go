package inits

import (
	"forum/pkg/globals"
	"forum/pkg/recovery"
	"github.com/gin-gonic/gin"
)

// RouterInit
// @Description: 初始化路由
// @Author tianjiajie 2024-10-05 16:00:37
func RouterInit() {
	// 初始化路由 Router
	globals.Router = gin.Default()

	// 某一个控制器报错，不影响整体（在发生 panic 时，将错误打印到日志中，并返回 500。）
	globals.Router.Use(recovery.CustomRecovery())

}
