package routes

import (
	"forum/internal/virtualMachine/controllers"
	"forum/pkg/token"
	"github.com/gin-gonic/gin"
)

// VirtualMachine
// @Description: 虚拟机相关
// @param        e *gin.Engine
// @Author tianjiajie 2025-03-15 16:13:45
func VirtualMachine(e *gin.Engine) {

	vmGroup := e.Group("/vm")
	// 不需要 token 校验
	{
		// pve回调api

		// 创建虚拟机回调
		vmGroup.POST("/create_callback", controllers.CreateVMCallback)
	}
	// token 校验
	vmGroup.Use(token.AuthMiddleware())
	{
		vmGroup.POST("/create", controllers.CreateVM)
	}

}
