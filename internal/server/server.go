package server

import (
	"fmt"
	"forum/pkg/globals"
	"github.com/spf13/viper"
	"log"
)

// Run 启动路由
func Run() {

	if err := viper.UnmarshalKey("app", &globals.AppConfig.App); err != nil {
		log.Fatalf("无法解码为结构: %s", err)
	}

	address := fmt.Sprintf("%s:%d", globals.AppConfig.App.Host, globals.AppConfig.App.Port)

	err := globals.Router.Run(address)
	if err != nil {
		globals.Log.Errorf("路由启动错误")
		return
	}

}
