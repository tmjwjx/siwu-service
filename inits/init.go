package inits

import (
	"forum/pkg/globals"
	"github.com/gin-gonic/gin"
)

func init() {
	// 优先初始化配置文件（给mysql，redis赋上配置信息）
	viperInit()

	// 初始化 mysql
	DBInit()

	// 初始化 redis
	RedisInit()

	// 初始化表
	TableInit()

	// 初始化日志文件
	logInit("logs", "forum")

	// 初始化路由 Router
	globals.Router = gin.Default()
}
