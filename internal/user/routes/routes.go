package routes

import (
	"forum/internal/user/controllers"
	"github.com/gin-gonic/gin"
)

// User 用户路由
func User(e *gin.Engine) {
	// 注册
	e.POST("/user/register", controllers.Register)
	// 忘记密码
	e.POST("/user/forgot_password", controllers.Register)
	// 用户请求验证码
	e.GET("/user/req_verify_code", controllers.ReqVerifyCode)
	// 登录
	e.POST("/user/login", controllers.Login)

	// 分组
	r := e.Group("/user")
	// token 校验
	// r.Use(token.AuthMiddleware())

	// 关注
	r.POST("/follow", controllers.Follow)
	// 用户排行
	r.GET("/rank", controllers.UserRank)

	// 上传用户个人资料
	r.POST("/form_personal_data", controllers.UserDataRequestCtrl)

	// 前端获取用户个人资料
	r.GET("/form_personal_data/:id", controllers.UserDataResponseCtrl)

	// 上传用户账号设置
	r.POST("/account_settings", controllers.UserAccountRequestCtrl)

	// 前端获取用户账号设置数据
	r.GET("/account_settings/:id", controllers.UserAccountResponseCtrl)

	// 上传用户私信设置
	r.POST("/private_settings", controllers.UserPrivateSetRequestCtrl)

	// 前端获取用户私信设置数据
	r.GET("/private_settings/:id", controllers.UserPrivateSetResponseCtrl)

	// 铃铛消息
	// e.GET("/event", controllers.MessagePushCtrl)

	// 用户管理
	// 重置用户密码
	r.POST("/reset", controllers.Reset)
	// 添加用户
	r.POST("/add", controllers.Add)
	// 删除用户
	r.DELETE("/delete", controllers.Delete)
	// 编辑用户
	r.POST("/edit", controllers.Edit)
	// 获取所有用户列表
	r.POST("/list", controllers.List)
	// 导入用户表
	r.POST("/import", controllers.Import)
	// 导出用户表
	r.GET("/export", controllers.Export)
	// 下载导入用户模版excel
	r.GET("/import/template", controllers.ImportTemplate)
	// 获取当前用户基本信息
	r.GET("/getInfo", controllers.GetInfo)
	// 上传用户头像
	r.POST("/upload/headshot", controllers.UploadHeadshot)

}
