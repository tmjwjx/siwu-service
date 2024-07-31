package server

import (
	"forum/pkg/utils"
)

// Run 启动路由
func Run() {
	err := utils.Router.Run("0.0.0.0:8081")
	if err != nil {
		utils.Log.Errorf("路由启动错误")
		return
	}

}
