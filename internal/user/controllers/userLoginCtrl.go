package controllers

import (
	"fmt"
	"forum/internal/internal_pkg/internal_utils"
	"forum/internal/user/logics"
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
	userLogic := logics.NewUserLogic(globals.DB, c, globals.SendEmailCfg)
	// 绑定数据
	var registerMsg requests.RegisterMsg
	if err := c.ShouldBind(&registerMsg); err != nil {
		response.Failed(c, http.StatusBadRequest, response.NewAppErr(globals.StatusBadRequest, fmt.Errorf("Register() err: %v", err), nil))
		return
	}

	// 判断数据是否合法

	// 检验邮箱是否合法
	if !internal_utils.IsValidEmail(registerMsg.Email) {
		response.Failed(c, http.StatusBadRequest, response.NewAppErr(globals.StatusBadRequest, fmt.Errorf("Register() : 邮箱不合法"), nil))
		return
	}

	// 核对两次输入的密码
	if registerMsg.Password != registerMsg.RePassword {
		response.Failed(c, http.StatusBadRequest, response.NewAppErr(globals.StatusBadRequest, fmt.Errorf("Register() : 两次输入的密码不相同"), nil))
		return
	}

	// 检验密码是否合法
	if !internal_utils.IsValidPassword(registerMsg.Password) {
		response.Failed(c, http.StatusBadRequest, response.NewAppErr(globals.StatusBadRequest, fmt.Errorf("Register() : 密码过于简单"), nil))
		return
	}

	// 业务逻辑
	err := userLogic.Register(registerMsg)
	if err != nil {
		response.Failed(c, http.StatusInternalServerError, response.NewAppErr(globals.StatusInternalServerError, fmt.Errorf("Register() -> %v", err), nil))
		return
	}

	// 成功
	response.Success(c, http.StatusOK, response.NewAppData(globals.StatusOK, "成功", nil))
}

// ReqVerifyCode 用户请求验证码
func ReqVerifyCode(c *gin.Context) {
	userLogic := logics.NewUserLogic(globals.DB, c, globals.SendEmailCfg)
	// 绑定数据
	var verifyCodeMsg requests.VerifyCodeMsg
	if err := c.ShouldBind(&verifyCodeMsg); err != nil {
		response.Failed(c, http.StatusBadRequest, response.NewAppErr(globals.StatusBadRequest, fmt.Errorf("ReqVerifyCode() -> %v", err), nil))
		return
	}

	// 检验邮箱是否合法
	if !internal_utils.IsValidEmail(verifyCodeMsg.Email) {
		response.Failed(c, http.StatusBadRequest, response.NewAppErr(globals.StatusBadRequest, fmt.Errorf("ReqVerifyCode() err: 邮箱不合法"), nil))
		return
	}

	// 业务逻辑
	if err := userLogic.ReqVerifyCode(verifyCodeMsg); err != nil {
		response.Failed(c, http.StatusInternalServerError, response.NewAppErr(globals.StatusInternalServerError, fmt.Errorf("ReqVerifyCode() -> %v", err), nil))
		return
	}

	// 成功
	response.Success(c, http.StatusOK, response.NewAppData(globals.StatusOK, "成功", nil))
}

// Login 登录
func Login(c *gin.Context) {
	userLogic := logics.NewUserLogic(globals.DB, c, globals.SendEmailCfg)
	// 绑定数据
	var logicMsg requests.LogicMsg
	if err := c.ShouldBind(&logicMsg); err != nil {
		response.Failed(c, http.StatusBadRequest, response.NewAppErr(globals.StatusBadRequest, fmt.Errorf("Login() -> %v", err), nil))
		return
	}

	email := logicMsg.Email

	// 判断数据是否合法

	// 检验邮箱是否合法
	if !internal_utils.IsValidEmail(logicMsg.Email) {
		response.Failed(c, http.StatusBadRequest, response.NewAppErr(globals.StatusBadRequest, fmt.Errorf("Login() : 邮箱不合法"), nil))
		return
	}
	// 检验密码是否合法
	if !internal_utils.IsValidPassword(logicMsg.Password) {
		response.Failed(c, http.StatusBadRequest, response.NewAppErr(globals.StatusBadRequest, fmt.Errorf("Login() : 密码不合法"), nil))
		return
	}

	// 业务逻辑
	if err := userLogic.Login(logicMsg); err != nil {
		response.Failed(c, http.StatusInternalServerError, response.NewAppErr(globals.StatusInternalServerError, fmt.Errorf("Login() -> %v", err), nil))
		return
	}

	// 成功
	// response.Success(c, http.StatusOK, response.NewAppData(globals.StatusOK, "成功", nil))

	// 生成token
	tok, err := token.GenerateToken(email)
	fmt.Println("生成的token为：", tok)
	if err != nil {
		response.Failed(c, http.StatusInternalServerError, response.NewAppErr(globals.StatusInternalServerError, fmt.Errorf("Login() -> %v", err), nil))
		return
	}
	response.Success(c, http.StatusOK, response.NewAppData(globals.StatusOK, "成功", tok))
}

// Follow 关注
func Follow(c *gin.Context) {
	userLogic := logics.NewUserLogic(globals.DB, c, globals.SendEmailCfg)
	// 绑定数据
	var followMsg requests.FollowMsg
	if err := c.ShouldBind(&followMsg); err != nil {
		response.Failed(c, http.StatusBadRequest, response.NewAppErr(globals.StatusBadRequest, fmt.Errorf("Follow() -> %v", err), nil))
		return
	}

	// 简单检验数据
	if followMsg.FollowerId == followMsg.FollowedId {
		response.Failed(c, http.StatusBadRequest, response.NewAppErr(globals.StatusBadRequest, fmt.Errorf("Follow() : id%d不能关注%d", followMsg.FollowerId, followMsg.FollowedId), nil))
		return
	}

	// 业务逻辑
	if err := userLogic.Follow(followMsg); err != nil {
		response.Failed(c, http.StatusInternalServerError, response.NewAppErr(globals.StatusInternalServerError, fmt.Errorf("Follow() -> %v", err), nil))
		return
	}

	// 成功
	response.Success(c, http.StatusOK, response.NewAppData(globals.StatusOK, "成功", nil))
}
