package server

import (
	"forum/internal/article/routes"
	userRouter "forum/internal/user/router"
	"forum/pkg/globals"
)

// SetupRouter 启动处理函数
func SetupRouter() {
	// 用户分路由
	userRouter.RouterInit(globals.Router)

	routes.Search()
}
