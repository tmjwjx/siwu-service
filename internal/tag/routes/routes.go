package routes

import (
	"forum/internal/tag/controllers"
	"forum/pkg/casbin"
	"forum/pkg/globals"
	"github.com/gin-gonic/gin"
)

func Tag(e *gin.Engine) {

	// 前台分组
	r := e.Group("/tag")

	// 更新标签关注人数
	r.POST("/fan_count", controllers.UpdateTagUserCount)

	// 刷新前端标签页
	r.GET("/article_count", controllers.UpdateTag)

	// 存储新用户选择的标签
	r.POST("/random_tag", controllers.StorageTagCtrl)

	// 获取所有标签的id和name
	r.GET("/get_all_tags", controllers.GetAllTagCtrl)

	// 后台分组
	r2 := e.Group("/backstage_tag").Use(casbin.CasbinAuth(globals.CasbinEnforcer))

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
	r2.GET("/batch_query", controllers.BatchQueryTagCtrl)

}
