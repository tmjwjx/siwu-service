package inits

import (
	"forum/pkg/globals"
	"github.com/spf13/viper"
	"log"
)

// ServerInit
// @Description: 服务初始化
// @Author tianjiajie 2024-10-11 09:22:49
func ServerInit() {

	// viper 提取配置文件中的app 即提取ip和端口
	if err := viper.UnmarshalKey("app", &globals.AppConfig.App); err != nil {
		log.Fatalf("无法解码为结构: %s", err)
	}

}