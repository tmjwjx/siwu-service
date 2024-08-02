package server

import (
	articleControl "forum/internal/article/controllers"
	imageControl "forum/internal/image/controllers"
	userControl "forum/internal/user/controllers"
	"forum/pkg/globals"
)

// SetupRouter 启动处理函数
func SetupRouter() {

	// 注册
	globals.Router.GET("/register", userControl.Register)

	// 搜索
	userGroup := globals.Router.Group("/search")
	{
		userGroup.GET("/search_box", articleControl.Search)
	}

	// 上传图片
	globals.Router.POST("/upload", imageControl.UploadHandlerControllers)

	// 上传个人资料
	globals.Router.POST("/form_personal_data", userControl.PersonalDataHandler)
}
