package routes

import (
	"forum/internal/tag/controllers"
	"github.com/gin-gonic/gin"
)

func Tag(e *gin.Engine) {
	// 前台分组
	r := e.Group("/tag")

	// 更新标签关注人数
	r.POST("/fan_count", controllers.UpdateTagUserCount)

	// 刷新前端标签页
	r.GET("/article_count", controllers.UpdateTag)

	// 后台分组
	r2 := e.Group("/backstage_tag")

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
