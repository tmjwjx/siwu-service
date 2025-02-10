package routes

import (
	"forum/internal/casbin_r_m_a/controllers"
	"forum/pkg/casbin"
	"forum/pkg/globals"
	"forum/pkg/token"
	"github.com/gin-gonic/gin"
)

// CasbinRMA
// @Description: casbin权限路由管理
// @Author wangyulong 2024-10-15 11:38:59
// @param        e *gin.Engine
func CasbinRMA(e *gin.Engine) {

	r := e.Group("/acl").Use(token.AuthMiddleware()).Use(casbin.CasbinAuth(globals.CasbinEnforcer))

	// 为角色分配菜单权限
	r.POST("/dispatch/role_menu", controllers.AssignMenuPermCtrl)

	// 获取当前角色的菜单权限(用于渲染侧边栏，只要type1和2)
	r.GET("/get_role_menu", controllers.GetMenuPermCtrl)

	// 获取当前角色的api权限
	r.GET("/get_role_api", controllers.GetApiPermCtrl)

	// 为角色分配api权限
	r.POST("/dispatch/role_api", controllers.AssignApiPermCtrl)

	// 获取当前角色的所有权限标识
	r.GET("/get_role_code", controllers.GetPermCodeCtrl)
}
