package server

import (
	articleRouter "forum/internal/article/routes"
	"forum/internal/image/routes"
	tagRouter "forum/internal/tag/routes"
	userRouter "forum/internal/user/routes"
	"forum/pkg/globals"
)

// SetupRouter 启动处理函数
func SetupRouter() {

	// 用户分路由
	userRouter.User(globals.Router)

	// 搜索分路由
	articleRouter.Search()

	// 发布文章
	articleRouter.Publish()

	// 标签分路由
	tagRouter.Tag(globals.Router)

	// 评论分路由
	articleRouter.Comment(globals.Router)

	// 图片url分路由
	routes.ProduceImageUrl(globals.Router)
}
