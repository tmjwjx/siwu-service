package server

import (
	articleContro "forum/internal/article/controllers"
	userContro "forum/internal/user/controllers"
)

// SetupRouter 启动处理函数
func SetupRouter() {
	// 注册
	Router.GET("/register", userContro.Register)

}

// Search 搜索
func Search() {

	userGroup := Router.Group("/search")
	{
		userGroup.GET("/search_box", articleContro.Search)
	}
	//r.GET("/search_box", controllers.Search)

}
