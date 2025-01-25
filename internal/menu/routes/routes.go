package routes

import (
	"forum/internal/menu/controllers"
	"github.com/gin-gonic/gin"
)

// 菜单管理 casbin鉴权

func Menu(e *gin.Engine) {
	/*casbinService, err := casbin_r_m_a.NewCasbinService(globals.DB)
	if err != nil {
		globals.Log.Errorf("casbin启动错误")
	}*/
	r := e.Group("/acl")

	{
		// 检索获取所有菜单列表
		r.POST("/menu/getlist", controllers.MenuSearchCtrl)

		// 获取所有菜单图标
		r.GET("/icon/list", controllers.GetMenuIconCtrl)

		// 新建菜单
		r.POST("/menu/add", controllers.CreateMenuCtrl)

		// 删除菜单
		r.DELETE("/menu/delete", controllers.DeleteMenuCtrl)

		// 修改菜单
		r.POST("/menu/update", controllers.UpdateMenuCtrl)

		// 获取当前菜单详情
		r.GET("/menu/detail", controllers.GetMenuDetailCtrl)

		// 获取所有type为1和2的菜单
		r.GET("/get_menu_type_list", controllers.GetSpecificMenuCtrl)
	}
}
