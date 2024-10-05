package server

import (
	"fmt"
	"forum/internal/internal_pkg/internal_utils"
	"forum/pkg/globals"
	"github.com/spf13/viper"
	"log"
)

// Run 启动路由
func Run() {

	//globals.Router.Static("/images", "./static/images")

	// 启动处理函数
	SetupRouter()

	// 配置静态文件目录
	// 将文件系统中的目录映射到 URL 路径
	if err := viper.UnmarshalKey("static", &globals.SConfig); err != nil {
		log.Fatalf("Run -> 无法解码为结构: %s", err)
	}

	globals.Router.Static(globals.SConfig.Prefix, globals.SConfig.Path)

	// 创建存储静态文件的目录路径文件夹
	err := internal_utils.CreateFolder(globals.SConfig.Path)
	if err != nil {
		globals.Log.Errorf("创建存储静态文件的目录路径文件夹")
	}

	// viper 提取配置文件中的app 即提取ip和端口
	if err := viper.UnmarshalKey("app", &globals.AppConfig.App); err != nil {
		log.Fatalf("无法解码为结构: %s", err)
	}
	address := fmt.Sprintf("%s:%d", globals.AppConfig.App.Host, globals.AppConfig.App.Port)

	// 启动路由
	err = globals.Router.Run(address)
	if err != nil {
		globals.Log.Errorf("路由启动错误")
		return
	}
	// 运行结束时 缓存区的信息写入到文件中#
	defer globals.Log.Sync()

}