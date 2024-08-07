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

	globals.Router.GET("/publish", controllers.ArticlePublishCtrl)

}