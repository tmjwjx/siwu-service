package inits

import (
	"forum/internal/internalPkg/internalUtils"
	"forum/pkg/globals"
	"github.com/spf13/viper"
	"log"
)

func StaticInit() {
	// 配置静态文件目录
	// 将文件系统中的目录映射到 URL 路径
	if err := viper.UnmarshalKey("static", &globals.SConfig); err != nil {
		log.Fatalf("Run -> 无法解码为结构: %s", err)
	}

	globals.Router.Static(globals.SConfig.Prefix, globals.SConfig.Path)

	// 创建存储静态文件的目录路径文件夹
	err := internalUtils.CreateFolder(globals.SConfig.Path)
	if err != nil {
		globals.Log.Errorf("创建存储静态文件的目录路径文件夹")
	}
}
