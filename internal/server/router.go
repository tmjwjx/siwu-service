package server

import (
	articleControl "forum/internal/article/controllers"
	userControl "forum/internal/user/controllers"
)

// SetupRouter 启动处理函数
func SetupRouter() {

	// 注册
	Router.GET("/register", userControl.Register)

	// 搜索
	userGroup := Router.Group("/search")
	{
		userGroup.GET("/search_box", articleControl.Search)
	}

}
