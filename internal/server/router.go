package server

import (
	tagRouter "forum/internal/tag/routes"
	userRouter "forum/internal/user/routes"
	"forum/pkg/globals"
)

// SetupRouter 启动处理函数
func SetupRouter() {

	// 用户分路由
	userRouter.User(globals.Router)

	// 搜索分路由
	//articleRouter.Search()

	// 标签分路由
	tagRouter.Tag(globals.Router)

}
