package main

import (
	_ "forum/inits" // 空导入，初始化
	"forum/internal/server"
)

func main() {
	
	//启动路由
	server.Run()
	
}