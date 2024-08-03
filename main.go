package main

import (
	"fmt"
	_ "forum/inits" // 空导入，初始化
	"forum/internal/server"
)

func init() {
	fmt.Println("hello")
}
func main() {
	//启动路由
	server.Run()
}
