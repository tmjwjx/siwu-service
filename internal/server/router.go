package server

import (
	"forum/internal/article/routes"
	tagRoutes "forum/internal/tag/routes"
	userRouter "forum/internal/user/routes"
	"forum/pkg/globals"
)

// SetupRouter 启动处理函数
func SetupRouter() {
	// 用户分路由
	userRouter.RouterInit(globals.Router)

	routes.Search()

	// 标签路由
	tagRoutes.RouterInit(globals.Router)

}
