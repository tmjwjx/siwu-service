package server

import (
	articleRouter "forum/internal/article/routes"
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

	// 角色分路由
	roleRouter.Role(globals.Router)
}
