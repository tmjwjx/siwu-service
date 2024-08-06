package main

import (
	_ "forum/inits" // 空导入，初始化
	"forum/internal/server"
	"forum/pkg/globals"
)

func main() {
	
	globals.Log.Info("main运行")
	
	//启动路由
	server.Run()
	
}