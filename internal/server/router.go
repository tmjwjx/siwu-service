package server

import (
	"forum/internal/article/routes"
	userControl "forum/internal/user/controllers"
	"forum/pkg/globals"
)

// SetupRouter 启动处理函数
func SetupRouter() {

	// 注册
	globals.Router.GET("/register", userControl.Register)

	routes.Search()
}
