package routes

import (
	"forum/internal/administrator/controllers"
	"forum/pkg/token"
	"github.com/gin-gonic/gin"
)

// Administrator
// @Description: 管理员
// @param        e *gin.Engine
// @Author tianjiajie 2025-02-19 19:27:57
func Administrator(e *gin.Engine) {
	// 管理员
	administratorGroup := e.Group("/backstage_manager")

	// token 校验
	administratorGroup.Use(token.AuthMiddleware())
	{
		// 添加管理员
		administratorGroup.POST("/add", controllers.AddAdministratorCtrl)
		// 删除管理员
		administratorGroup.DELETE("/delete", controllers.DeleteAdministrator)
		//// 编辑管理员信息
		//administratorGroup.POST("/update", controllers.UpdateAdministrator)
		//// 查询管理员列表
		//administratorGroup.GET("/list", controllers.GetAdministratorList)
		//// 查询管理员详情
		//administratorGroup.GET("/info", controllers.GetAdministratorInfo)
		//// 重置管理员密码
		//administratorGroup.POST("/reset", controllers.ResetAdministratorPassword)
	}
}
