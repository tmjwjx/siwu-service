package controllers

import (
	"fmt"
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
	// // 绑定数据
	// var registerMsg requests.RegisterReq
	// if err := c.ShouldBind(&registerMsg); err != nil {
	// 	response.Failed(c, http.StatusBadRequest, response.NewAppErr(globals.StatusBadRequest, fmt.Errorf("Register() err: %v", err), nil))
	// 	return
	// }
	//
	// // 判断数据是否合法
	//
	// // 检验邮箱是否合法
	// if !internalUtils.IsValidEmail(registerMsg.Email) {
	// 	response.Failed(c, http.StatusBadRequest, response.NewAppErr(globals.StatusBadRequest, fmt.Errorf("Register() : 邮箱不合法"), nil))
	// 	return
	// }
	//
	// // 核对两次输入的密码
	// if registerMsg.Password != registerMsg.RePassword {
	// 	response.Failed(c, http.StatusBadRequest, response.NewAppErr(globals.StatusBadRequest, fmt.Errorf("Register() : 两次输入的密码不相同"), nil))
	// 	return
	// }
	//
	// // 检验密码是否合法
	// if !internalUtils.IsValidPassword(registerMsg.Password) {
	// 	response.Failed(c, http.StatusBadRequest, response.NewAppErr(globals.StatusBadRequest, fmt.Errorf("Register() : 密码必须要同时包含字母、数字、特殊字符，长度在8到20位之间"), nil))
	// 	return
	// }

	req, ok := c.Get("req")
	if !ok {
		globals.Log.Errorf(response.ErrGetReqIsWrong)
		response.Failed(c, http.StatusInternalServerError, response.NewAppErr(globals.StatusInternalServerError, fmt.Errorf(response.ErrGetReqIsWrong), nil))
		return
	}
	// 类型断言
	registerReq, ok := req.(requests.RegisterReq)
	if !ok {
		globals.Log.Errorf(response.ErrTypeAssertionFail)
		response.Failed(c, http.StatusInternalServerError, response.NewAppErr(globals.StatusInternalServerError, fmt.Errorf(response.ErrTypeAssertionFail), nil))
		return
	}

	// 业务逻辑
	userReqContext := logics.NewUserReqContext(globals.DB, c, globals.SendEmailCfg)
	if err := userReqContext.Register(registerReq); err != nil {
		// response.Failed(c, http.StatusInternalServerError, response.NewAppErr(globals.StatusInternalServerError, err, nil))
		globals.Log.Errorf(err.Error())
		response.Failed(c, http.StatusInternalServerError, response.NewAppErr(globals.StatusInternalServerError, err, nil))
		return
	}

	// 成功
	response.Success(c, http.StatusOK, response.NewAppData(globals.StatusOK, response.DataSuccess, nil))
}

// ReqVerifyCode 用户请求验证码
func ReqVerifyCode(c *gin.Context) {
	// // 绑定数据
	// email := c.Query("email")
	//
	// // 数据检验
	//
	// // 检验邮箱是否合法
	// if !internalUtils.IsValidEmail(email) {
	// 	response.Failed(c, http.StatusBadRequest, response.NewAppErr(globals.StatusBadRequest, fmt.Errorf("ReqVerifyCode() err: 邮箱不合法"), nil))
	// 	return
	// }

	req, ok := c.Get("req")
	if !ok {
		globals.Log.Errorf(response.ErrGetReqIsWrong)
		response.Failed(c, http.StatusInternalServerError, response.NewAppErr(globals.StatusInternalServerError, fmt.Errorf(response.ErrGetReqIsWrong), nil))
		return
	}
	// 类型断言
	email, ok := req.(string)
	if !ok {
		globals.Log.Errorf(response.ErrTypeAssertionFail)
		response.Failed(c, http.StatusInternalServerError, response.NewAppErr(globals.StatusInternalServerError, fmt.Errorf(response.ErrTypeAssertionFail), nil))
		return
	}

	// 业务逻辑
	userReqContext := logics.NewUserReqContext(globals.DB, c, globals.SendEmailCfg)
	if err := userReqContext.ReqVerifyCode(email); err != nil {
		// response.Failed(c, http.StatusInternalServerError, response.NewAppErr(globals.StatusInternalServerError, fmt.Errorf("ReqVerifyCode() -> %v", err), nil))
		globals.Log.Errorf(err.Error())
		response.Failed(c, http.StatusInternalServerError, response.NewAppErr(globals.StatusInternalServerError, err, nil))
		return
	}

	// 成功
	response.Success(c, http.StatusOK, response.NewAppData(globals.StatusOK, response.DataSuccess, nil))
}

// Login 登录
func Login(c *gin.Context) {
	// // 绑定数据
	// var logicMsg requests.LogicReq
	// if err := c.ShouldBind(&logicMsg); err != nil {
	// 	response.Failed(c, http.StatusBadRequest, response.NewAppErr(globals.StatusBadRequest, fmt.Errorf("Login() -> %v", err), nil))
	// 	return
	// }
	//
	// // 判断数据是否合法
	//
	// // 检验邮箱是否合法
	// if !internalUtils.IsValidEmail(logicMsg.Email) {
	// 	response.Failed(c, http.StatusBadRequest, response.NewAppErr(globals.StatusBadRequest, fmt.Errorf("Login() : 邮箱不合法"), nil))
	// 	return
	// }

	req, ok := c.Get("req")
	if !ok {
		globals.Log.Errorf(response.ErrGetReqIsWrong)
		response.Failed(c, http.StatusInternalServerError, response.NewAppErr(globals.StatusInternalServerError, fmt.Errorf(response.ErrGetReqIsWrong), nil))
		return
	}
	// 类型断言
	loginReq, ok := req.(requests.LogicReq)
	if !ok {
		globals.Log.Errorf(response.ErrTypeAssertionFail)
		response.Failed(c, http.StatusInternalServerError, response.NewAppErr(globals.StatusInternalServerError, fmt.Errorf(response.ErrTypeAssertionFail), nil))
		return
	}

	// 业务逻辑
	userReqContext := logics.NewUserReqContext(globals.DB, c, globals.SendEmailCfg)
	userInfo, err := userReqContext.Login(loginReq)
	if err != nil {
		globals.Log.Errorf(err.Error())
		response.Failed(c, http.StatusInternalServerError, response.NewAppErr(globals.StatusInternalServerError, err, nil))
		return
	}

	// 通过email查询id
	user := repositories.QueryUserByEmail(userReqContext.DB, loginReq.Email)
	if user == nil {
		globals.Log.Errorf(response.ErrEmailNotExist + ":" + loginReq.Email)
		response.Failed(c, http.StatusInternalServerError, response.NewAppErr(globals.StatusInternalServerError, fmt.Errorf(response.ErrEmailNotExist+":"+loginReq.Email), nil))
		return
	}
	// 生成token
	tok, err := token.GenerateToken(user.ID)
	if err != nil {
		globals.Log.Errorf(err.Error())
		response.Failed(c, http.StatusInternalServerError, response.NewAppErr(globals.StatusInternalServerError, err, nil))
		return
	}
	response.Success(c, http.StatusOK, response.NewAppData(globals.StatusOK, response.DataSuccess, gin.H{"token": tok, "userinfo": userInfo}))
}

// ForgotPassword 忘记密码
func ForgotPassword(c *gin.Context) {
	// // 绑定数据
	// var forgotPasswordMsg requests.ForgotPasswordReq
	// if err := c.ShouldBind(&forgotPasswordMsg); err != nil {
	// 	response.Failed(c, http.StatusBadRequest, response.NewAppErr(globals.StatusBadRequest, fmt.Errorf("Register() err: %v", err), nil))
	// 	return
	// }
	//
	// // 判断数据是否合法
	//
	// // 检验邮箱是否合法
	// if !internalUtils.IsValidEmail(forgotPasswordMsg.Email) {
	// 	response.Failed(c, http.StatusBadRequest, response.NewAppErr(globals.StatusBadRequest, fmt.Errorf("Register() : 邮箱不合法"), nil))
	// 	return
	// }
	//
	// // 核对两次输入的密码
	// if forgotPasswordMsg.Password != forgotPasswordMsg.RePassword {
	// 	response.Failed(c, http.StatusBadRequest, response.NewAppErr(globals.StatusBadRequest, fmt.Errorf("Register() : 两次输入的密码不相同"), nil))
	// 	return
	// }
	//
	// // 检验密码是否合法
	// if !internalUtils.IsValidPassword(forgotPasswordMsg.Password) {
	// 	response.Failed(c, http.StatusBadRequest, response.NewAppErr(globals.StatusBadRequest, fmt.Errorf("Register() : 密码必须要同时包含字母、数字、特殊字符，长度在8到20位之间"), nil))
	// 	return
	// }

	req, ok := c.Get("req")
	if !ok {
		globals.Log.Errorf(response.ErrGetReqIsWrong)
		response.Failed(c, http.StatusInternalServerError, response.NewAppErr(globals.StatusInternalServerError, fmt.Errorf(response.ErrGetReqIsWrong), nil))
		return
	}
	// 类型断言
	forgotPasswordReq, ok := req.(requests.ForgotPasswordReq)
	if !ok {
		globals.Log.Errorf(response.ErrTypeAssertionFail)
		response.Failed(c, http.StatusInternalServerError, response.NewAppErr(globals.StatusInternalServerError, fmt.Errorf(response.ErrTypeAssertionFail), nil))
		return
	}

	// 业务逻辑
	userReqContext := logics.NewUserReqContext(globals.DB, c, globals.SendEmailCfg)
	if err := userReqContext.ForgotPassword(forgotPasswordReq); err != nil {
		// response.Failed(c, http.StatusInternalServerError, response.NewAppErr(globals.StatusInternalServerError, fmt.Errorf("Register() -> %v", err), nil))
		globals.Log.Errorf(err.Error())
		response.Failed(c, http.StatusInternalServerError, response.NewAppErr(globals.StatusInternalServerError, err, nil))
		return
	}

	// 成功
	response.Success(c, http.StatusOK, response.NewAppData(globals.StatusOK, response.DataSuccess, nil))
}

// Logout 登出
func Logout(c *gin.Context) {
	tokenString := c.GetHeader("Authorization")
	if tokenString == "" {
		globals.Log.Errorf(response.ErrMissAuthorizationHeader)
		response.Failed(c, http.StatusUnauthorized, response.NewAppErr(globals.StatusUnauthorized, fmt.Errorf(response.ErrMissAuthorizationHeader), nil))
		return
	}
	tokenString = tokenString[len("Bearer "):]

	// 业务逻辑
	bsManageContext := logics.NewUserReqContext(globals.DB, c, globals.SendEmailCfg)
	if err := bsManageContext.Logout(tokenString); err != nil {
		globals.Log.Errorf(err.Error())
		response.Failed(c, http.StatusInternalServerError, response.NewAppErr(globals.StatusInternalServerError, err, nil))
		return
	}

	response.Success(c, http.StatusOK, response.NewAppData(globals.StatusOK, response.DataSuccess, nil))
}
