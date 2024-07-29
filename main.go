package main

import (
	"forum/inits"
	"forum/internal/server"
)

func main() {
	// 初始化
	inits.Init()

	// 启动处理函数
	server.SetupRouter()
	// 搜索
	server.Search()

	// 启动路由
	server.Run()

}
