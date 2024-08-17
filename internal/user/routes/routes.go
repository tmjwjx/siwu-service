package routes

import (
	"forum/internal/user/controllers"
	"github.com/gin-gonic/gin"
)

// User 用户路由
func User(e *gin.Engine) {
	// 分组
	r := e.Group("/user")
	// 注册
	r.POST("/register", controllers.Register)
	// 用户请求验证码
	r.GET("req_verify_code", controllers.ReqVerifyCode)

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
}
