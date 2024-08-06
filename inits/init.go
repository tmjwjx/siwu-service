package inits

import (
	"forum/pkg/globals"
	"github.com/gin-gonic/gin"
)

func init() {
	// 初始化环境
	EnvInit()

	// 优先初始化配置文件（给mysql，redis赋上配置信息）
	ConfigInit()

	// 初始化 mysql
	DBInit()

	// 初始化 redis
	// RedisInit()

	// 初始化表
	TableInit()

	// 初始化日志文件
	LogInit("logs", "forum")

	// 初始化路由 Router
	globals.Router = gin.Default()

	// 某一个控制器报错，不影响整体
	// 日志自动记录
	globals.Router.Use(gin.Recovery())
}
