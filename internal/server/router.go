package server

import (
	articleRouter "forum/internal/article/routes"
	userRouter "forum/internal/user/routes"
	"forum/pkg/globals"
)

// SetupRouter 启动处理函数
func SetupRouter() {
	
	// 用户分路由
	userRouter.RouterInit(globals.Router)
	
	//搜索
	articleRouter.Search()
	
}