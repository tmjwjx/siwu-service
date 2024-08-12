package routes

import (
	"forum/internal/user/controllers"
	"github.com/gin-gonic/gin"
)

// User 用户路由
func User(e *gin.Engine) {

	// 铃铛消息
	//e.GET("/event", controllers.MessagePushCtrl)

	// 分组
	r := e.Group("/user")
	// 注册
	r.POST("/register", controllers.Register)
	// 用户请求验证码
	r.GET("req_verify_code", controllers.ReqVerifyCode)
	// 登录
	r.POST("/login", controllers.Login)
	// 关注
	r.POST("/follow", controllers.Follow)

	// 上传个人资料
	e.POST("/form_personal_data", controllers.PersonalDataHandler)

	// 前端获取用户资料
	e.GET("/form_personal_data/:id", controllers.ResponsePersonDate)
}
