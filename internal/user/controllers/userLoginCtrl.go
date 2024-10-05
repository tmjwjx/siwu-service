package controllers

import (
	"fmt"
	"forum/internal/internalPkg/internalUtils"
	"forum/internal/user/logics"
	"forum/internal/user/repositories"
	"forum/internal/user/requests"
	"forum/pkg/globals"
	"forum/pkg/response"
	"forum/pkg/token"
	"github.com/gin-gonic/gin"
	"net/http"
)

// 注册、登陆、验证码

// Register 用户注册
func Register(c *gin.Context) {
	// 绑定数据
	var registerMsg requests.RegisterReq
	if err := c.ShouldBind(&registerMsg); err != nil {
		response.Failed(c, http.StatusBadRequest, response.NewAppErr(globals.StatusBadRequest, fmt.Errorf("Register() err: %v", err), nil))
		return
	}

	// 判断数据是否合法

	// 检验邮箱是否合法
	if !internalUtils.IsValidEmail(registerMsg.Email) {
		response.Failed(c, http.StatusBadRequest, response.NewAppErr(globals.StatusBadRequest, fmt.Errorf("Register() : 邮箱不合法"), nil))
		return
	}

	// 核对两次输入的密码
	if registerMsg.Password != registerMsg.RePassword {
		response.Failed(c, http.StatusBadRequest, response.NewAppErr(globals.StatusBadRequest, fmt.Errorf("Register() : 两次输入的密码不相同"), nil))
		return
	}

	// 检验密码是否合法
	if !internalUtils.IsValidPassword(registerMsg.Password) {
		response.Failed(c, http.StatusBadRequest, response.NewAppErr(globals.StatusBadRequest, fmt.Errorf("Register() : 密码过于简单"), nil))
		return
	}

	// 业务逻辑
	userReqContext := logics.NewUserReqContext(globals.DB, c, globals.SendEmailCfg)
	err := userReqContext.Register(registerMsg)
	if err != nil {
		response.Failed(c, http.StatusInternalServerError, response.NewAppErr(globals.StatusInternalServerError, fmt.Errorf("Register() -> %v", err), nil))
		return
	}

	// 成功
	response.Success(c, http.StatusOK, response.NewAppData(globals.StatusOK, "成功", nil))
}

// ReqVerifyCode 用户请求验证码
func ReqVerifyCode(c *gin.Context) {
	// 绑定数据
	email := c.Query("email")

	// 数据检验

	// 检验邮箱是否合法
	if !internalUtils.IsValidEmail(email) {
		response.Failed(c, http.StatusBadRequest, response.NewAppErr(globals.StatusBadRequest, fmt.Errorf("ReqVerifyCode() err: 邮箱不合法"), nil))
		return
	}

	// 业务逻辑
	userReqContext := logics.NewUserReqContext(globals.DB, c, globals.SendEmailCfg)
	if err := userReqContext.ReqVerifyCode(email); err != nil {
		response.Failed(c, http.StatusInternalServerError, response.NewAppErr(globals.StatusInternalServerError, fmt.Errorf("ReqVerifyCode() -> %v", err), nil))
		return
	}

	// 成功
	response.Success(c, http.StatusOK, response.NewAppData(globals.StatusOK, "成功", nil))
}

// Login 登录
func Login(c *gin.Context) {
	// 绑定数据
	var logicMsg requests.LogicReq
	if err := c.ShouldBind(&logicMsg); err != nil {
		response.Failed(c, http.StatusBadRequest, response.NewAppErr(globals.StatusBadRequest, fmt.Errorf("Login() -> %v", err), nil))
		return
	}

	// 判断数据是否合法

	// 检验邮箱是否合法
	if !internalUtils.IsValidEmail(logicMsg.Email) {
		response.Failed(c, http.StatusBadRequest, response.NewAppErr(globals.StatusBadRequest, fmt.Errorf("Login() : 邮箱不合法"), nil))
		return
	}
	// 检验密码是否合法
	if !internalUtils.IsValidPassword(logicMsg.Password) {
		response.Failed(c, http.StatusBadRequest, response.NewAppErr(globals.StatusBadRequest, fmt.Errorf("Login() : 密码不合法"), nil))
		return
	}

	// 业务逻辑
	userReqContext := logics.NewUserReqContext(globals.DB, c, globals.SendEmailCfg)
	if err := userReqContext.Login(logicMsg); err != nil {
		response.Failed(c, http.StatusInternalServerError, response.NewAppErr(globals.StatusInternalServerError, fmt.Errorf("Login() -> %v", err), nil))
		return
	}

	// 通过email查询id
	user := repositories.QueryUserByEmail(userReqContext.DB, logicMsg.Email)
	if user == nil {
		response.Failed(c, http.StatusInternalServerError, response.NewAppErr(globals.StatusInternalServerError, fmt.Errorf("Login() err: 不存在email为 %v 的用户", logicMsg.Email), nil))
		return
	}
	// 生成token
	tok, err := token.GenerateToken(user.ID)
	if err != nil {
		response.Failed(c, http.StatusInternalServerError, response.NewAppErr(globals.StatusInternalServerError, fmt.Errorf("Login() -> %v", err), nil))
		return
	}
	response.Success(c, http.StatusOK, response.NewAppData(globals.StatusOK, "成功", gin.H{"token": tok}))
}