package server

import (
	"forum/internal/user/controllers"
	"github.com/gin-gonic/gin"
)

var Router *gin.Engine

// SetupRouter 启动处理函数
func SetupRouter() {
	// 注册
	Router.GET("/register", controllers.Register)

}
