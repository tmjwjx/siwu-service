package routes

import (
	"forum/internal/article/controllers"
	"forum/pkg/globals"
)

func Search() {
	// 搜索
	userGroup := globals.Router.Group("/search")
	{
		userGroup.GET("/query", controllers.ArticleSearchCtrl)
	}

}

func Publish() {

	// 搜索
	articleGroup := globals.Router.Group("/article")
	{
		// 编辑界面
		articleGroup.GET("/edit", controllers.ArticleEditCtrl)
		// 发布文章
		articleGroup.POST("/publish", controllers.ArticlePublishCtrl)
	}

}