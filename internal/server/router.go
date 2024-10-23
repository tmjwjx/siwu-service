package server

import (
	ApiRouter "forum/internal/api/routes"
	articleRouter "forum/internal/article/routes"
	bsLoginRouter "forum/internal/backstage/routes"
	casbinRouter "forum/internal/casbin_r_m_a/routes"
	dictRouter "forum/internal/dictManage/routes"
	imageRouter "forum/internal/image/routes"
	MenuRouter "forum/internal/menu/routes"
	messageRouter "forum/internal/message/routes"
	roleRouter "forum/internal/roleManage/routes"
	tagRouter "forum/internal/tag/routes"
	userRouter "forum/internal/user/routes"
	"forum/pkg/globals"
)

// SetupRouter 启动处理函数
func SetupRouter() {

	// 用户分路由
	userRouter.User(globals.Router)

	// 消息分路由
	messageRouter.Message(globals.Router)

	// 文章分路由
	articleRouter.Article(globals.Router)

	// 标签分路由
	tagRouter.Tag(globals.Router)

	// 评论分路由
	articleRouter.Comment(globals.Router)

	// 图片url分路由
	imageRouter.ProduceImageUrl(globals.Router)

	// 角色管理分路由
	roleRouter.Role(globals.Router)

	// 字典管理分路由
	dictRouter.Dict(globals.Router)

	// api管理分录由
	ApiRouter.Api(globals.Router)

	// 菜单管理分录由
	MenuRouter.Menu(globals.Router)

	// casbin权限管理分录由
	casbinRouter.CasbinRMA(globals.Router)

	// 后台登陆分路由
	bsLoginRouter.Backstage(globals.Router)
}
