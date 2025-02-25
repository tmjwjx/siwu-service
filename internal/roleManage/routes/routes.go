package routes

import (
	"forum/internal/roleManage/controllers"
	"forum/pkg/casbin"
	"forum/pkg/globals"
	"forum/pkg/token"
	"github.com/gin-gonic/gin"
)

// Role 角色分路由
func Role(e *gin.Engine) {
	// 分组
	r := e.Group("/role").Use(token.AuthMiddleware()).Use(casbin.CasbinAuth(globals.CasbinEnforcer))

	// 添加角色
	r.POST("/add_role", controllers.AddRole)
	// 删除角色
	r.DELETE("/delete_role", controllers.DeleteRole)
	// 检索角色
	r.POST("/search_role", controllers.SearchRole)
	// 更新角色
	r.POST("/update_role", controllers.UpdateRole)
	// 获取所有已启用的角色名称列表
	r.GET("/get_role_name", controllers.GetRoleName)
	// 获取当前角色详情
	r.GET("/get_detail", controllers.GetDetail)
	// 为用户分配角色
	r.POST("/dispatch_role", controllers.DispatchRole)
}
