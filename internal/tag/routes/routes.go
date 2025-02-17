package routes

import (
	"forum/internal/tag/controllers"
	"forum/pkg/casbin"
	"forum/pkg/globals"
	"forum/pkg/token"
	"github.com/gin-gonic/gin"
)

func Tag(e *gin.Engine) {

	// 刷新前端标签页
	e.GET("/article_count", controllers.UpdateTag)

	// 获取所有标签的id和name
	e.GET("/get_all_tags", controllers.GetAllTagCtrl)

	// 前台分组
	r := e.Group("/tag").Use(token.AuthMiddleware())
	{
		// 更新标签关注人数
		r.POST("/fan_count", controllers.UpdateTagUserCount)

		// 存储新用户选择的标签
		r.POST("/random_tag", controllers.StorageTagCtrl)
	}

	// 后台分组
	r2 := e.Group("/backstage_tag").Use(token.AuthMiddleware()).Use(casbin.CasbinAuth(globals.CasbinEnforcer))
	{
		// 新增标签
		r2.POST("/add", controllers.AddTagCtrl)

		// 删除标签
		r2.POST("/delete", controllers.DeleteTagCtrl)

		// 批量删除标签
		r2.POST("/batch_delete", controllers.BatchDelTagCtrl)

		// 更新标签
		r2.POST("/update", controllers.UpdateTagCtrl)

		// 查询标签
		r2.GET("/query", controllers.QueryTagCtrl)

		// 批量查询标签
		r2.POST("/batch_query", controllers.BatchQueryTagCtrl)
	}

}
