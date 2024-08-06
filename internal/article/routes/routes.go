package routes

import (
	articleControl "forum/internal/article/controllers"
	"forum/pkg/globals"
)

func Search() {
	// 搜索
	userGroup := globals.Router.Group("/search")
	{
		userGroup.GET("/query", articleControl.SearchHandler)
		userGroup.GET("/test", articleControl.Test1)
	}
	
}