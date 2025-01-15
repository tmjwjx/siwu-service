package server

import (
	"fmt"
	"forum/pkg/globals"
)

// Run 启动路由
func Run() {

	// 启动处理函数
	SetupRouter()

	address := fmt.Sprintf("%s:%d", globals.AppConfig.App.Host, globals.AppConfig.App.Port)

	// 启动路由
	err := globals.Router.Run(address)
	if err != nil {
		globals.Log.Errorf("路由启动错误")
		return
	}
	// 运行结束时 缓存区的信息写入到文件中#
	defer globals.Log.Sync()

}
