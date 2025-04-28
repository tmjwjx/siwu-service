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
		// 销毁虚拟机毁掉
		vmGroup.POST("/destroy_callback", controllers.DestroyVMCallback)
		// 虚拟机通知信息
		vmGroup.POST("/notice_callback", controllers.NoticeVMCallback)
	}

	// token 校验
	vmGroup.Use(token.AuthMiddleware())
	{
		vmGroup.POST("/create", controllers.CreateVM)

		// 用户手动删除虚拟机
		vmGroup.POST("/destroy", controllers.DestroyVM)
	}

}
