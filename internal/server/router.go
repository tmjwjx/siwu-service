package server

import (
	articleRouter "forum/internal/article/routes"
	"forum/internal/image/routes"
	dictRouter "forum/internal/dictManage/routes"
	roleRouter "forum/internal/roleManage/routes"
	tagRouter "forum/internal/tag/routes"
	userRouter "forum/internal/user/routes"
	"forum/pkg/globals"
)

// SetupRouter 启动处理函数
func SetupRouter() {

	// 用户分路由
	userRouter.User(globals.Router)

	// 文章分路由
	articleRouter.Article(globals.Router)

	// 标签分路由
	tagRouter.Tag(globals.Router)

	// 评论分路由
	articleRouter.Comment(globals.Router)

	// 图片url分路由
	routes.ProduceImageUrl(globals.Router)
	// 角色管理分路由
	roleRouter.Role(globals.Router)

	// 字典管理分路由
	dictRouter.Dict(globals.Router)
}
