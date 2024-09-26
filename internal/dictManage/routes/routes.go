package routes

import (
	"forum/internal/dictManage/controllers"
	"github.com/gin-gonic/gin"
)

// Dict 字典管理分路由
func Dict(e *gin.Engine) {
	// 分组
	r := e.Group("/dict")
	// token 校验
	// r.Use(token.AuthMiddleware())

	// 新增字典类型
	r.POST("/add_type", controllers.AddType)
	// 批量删除字典类型
	r.DELETE("/delete_type", controllers.DeleteType)
	// 修改字典类型
	r.POST("/update_type", controllers.UpdateType)
	// 获取字典类型
	r.GET("/get_type", controllers.GetType)

	// 新增字典项
	r.POST("/add_item", controllers.AddItem)
	// 批量删除字典项
	r.DELETE("/delete_item", controllers.DeleteItem)
	// 修改字典项
	r.POST("/update_item", controllers.UpdateItem)
	// 获取字典项
	r.GET("/get_item", controllers.GetItem)
}
