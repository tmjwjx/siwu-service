package server

import (
	articleControl "forum/internal/article/controllers"
	userControl "forum/internal/user/controllers"
	"forum/pkg/globals"
)

// SetupRouter 启动处理函数
func SetupRouter() {

	// 注册
	globals.Router.GET("/register", userControl.Register)

	// 搜索
	userGroup := globals.Router.Group("/search")
	{
		userGroup.GET("/search_box", articleControl.Search)
	}
}
