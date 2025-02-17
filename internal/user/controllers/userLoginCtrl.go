package controllers

import (
	"fmt"
	"forum/internal/internalPkg/internalUtils"
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
	// 绑定数据
	var registerReq requests.RegisterReq
	if err := c.ShouldBind(&registerReq); err != nil {
		response.Failed(c, http.StatusBadRequest, response.NewAppErr(globals.StatusBadRequest, fmt.Errorf("Register() err: %v", err), nil))
		return
	}

	// 判断数据是否合法

	// 检验邮箱是否合法
	if !internalUtils.IsValidEmail(registerReq.Email) {
		response.Failed(c, http.StatusBadRequest, response.NewAppErr(globals.StatusBadRequest, fmt.Errorf("Register() : 邮箱不合法"), nil))
		return
	}

	// 核对两次输入的密码
	if registerReq.Password != registerReq.RePassword {
		response.Failed(c, http.StatusBadRequest, response.NewAppErr(globals.StatusBadRequest, fmt.Errorf("Register() : 两次输入的密码不相同"), nil))
		return
	}

	// 检验密码是否合法
	if !internalUtils.IsValidPassword(registerReq.Password) {
		response.Failed(c, http.StatusBadRequest, response.NewAppErr(globals.StatusBadRequest, fmt.Errorf("Register() : 密码必须要同时包含字母、数字、特殊字符，长度在8到20位之间"), nil))
		return
	}

	// req, ok := c.Get("req")
	// if !ok {
	// 	globals.Log.Errorf(response.ErrGetReqIsWrong)
	// 	response.Failed(c, http.StatusInternalServerError, response.NewAppErr(globals.StatusInternalServerError, fmt.Errorf(response.ErrGetReqIsWrong), nil))
	// 	return
	// }
	// // 类型断言
	// registerReq, ok := req.(requests.RegisterReq)
	// if !ok {
	// 	globals.Log.Errorf(response.ErrTypeAssertionFail)
	// 	response.Failed(c, http.StatusInternalServerError, response.NewAppErr(globals.StatusInternalServerError, fmt.Errorf(response.ErrTypeAssertionFail), nil))
	// 	return
	// }

	// 业务逻辑
	userReqContext := logics.NewUserReqContext(globals.DB, c, globals.SendEmailCfg)
	userInfo, err := userReqContext.Register(registerReq)
	if err != nil {
		// response.Failed(c, http.StatusInternalServerError, response.NewAppErr(globals.StatusInternalServerError, err, nil))
		globals.Log.Errorf(err.Error())
		response.Failed(c, http.StatusInternalServerError, response.NewAppErr(globals.StatusInternalServerError, err, nil))
		return
	}

	// // 通过email查询id
	// user := repositories.QueryUserByEmail(userReqContext.DB, registerReq.Email)
	// if user == nil {
	// 	globals.Log.Errorf(response.ErrEmailNotExist + ":" + registerReq.Email)
	// 	response.Failed(c, http.StatusInternalServerError, response.NewAppErr(globals.StatusInternalServerError, fmt.Errorf(response.ErrEmailNotExist+":"+registerReq.Email), nil))
	// 	return
	// }
	// // 生成token
	// tok, err := token.GenerateToken(user.ID)
	// if err != nil {
	// 	globals.Log.Errorf(err.Error())
	// 	response.Failed(c, http.StatusInternalServerError, response.NewAppErr(globals.StatusInternalServerError, err, nil))
	// 	return
	// }

	// 查询用户的信息

	// // 通过email查询id
	// user := repositories.QueryUserByEmail(userReqContext.DB, registerReq.Email)
	// if user == nil {
	// 	globals.Log.Errorf(response.ErrEmailNotExist + ":" + registerReq.Email)
	// 	response.Failed(c, http.StatusInternalServerError, response.NewAppErr(globals.StatusInternalServerError, fmt.Errorf(response.ErrEmailNotExist+":"+registerReq.Email), nil))
	// 	return
	// }
	// userImages, err := internalUtils.GetImages(u.DB, globals.UserHome, user.ID)
	// if err != nil {
	// 	globals.Log.Errorf(err.Error())
	// 	response.Failed(c, http.StatusInternalServerError, response.NewAppErr(globals.StatusInternalServerError, fmt.Errorf(response.ErrEmailNotExist+":"+registerReq.Email), nil))
	// 	return
	// 	// return nil, fmt.Errorf("UserReqContext.Login() %v", err)
	// }
	// // 没有图片
	// if userImages == nil {
	// 	// return nil, fmt.Errorf("UserReqContext.Login() err = 无法找到id为%d的用户头像图片", user.ID)
	// 	globals.Log.Errorf(response.ErrUnableFindUserAvatar + ":" + strconv.Itoa(int(user.ID)))
	//
	// 	return
	// }
	// avatarPath := (*userImages)[0]
	//
	// var userInfo = requests.LogicRes{
	// 	Id:         user.ID,
	// 	Nickname:   user.Nickname,
	// 	AvatarPath: avatarPath,
	// }

	// 生成token
	tok, err := token.GenerateToken(userInfo.Id)
	if err != nil {
		globals.Log.Errorf(err.Error())
		response.Failed(c, http.StatusInternalServerError, response.NewAppErr(globals.StatusInternalServerError, err, nil))
		return
	}
	response.Success(c, http.StatusOK, response.NewAppData(globals.StatusOK, response.DataSuccess, gin.H{"token": tok, "userinfo": userInfo}))

	// // 成功
	// response.Success(c, http.StatusOK, response.NewAppData(globals.StatusOK, response.DataSuccess, gin.H{"token": tok}))
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

	// req, ok := c.Get("req")
	// if !ok {
	// 	globals.Log.Errorf(response.ErrGetReqIsWrong)
	// 	response.Failed(c, http.StatusInternalServerError, response.NewAppErr(globals.StatusInternalServerError, fmt.Errorf(response.ErrGetReqIsWrong), nil))
	// 	return
	// }
	// // 类型断言
	// email, ok := req.(string)
	// if !ok {
	// 	globals.Log.Errorf(response.ErrTypeAssertionFail)
	// 	response.Failed(c, http.StatusInternalServerError, response.NewAppErr(globals.StatusInternalServerError, fmt.Errorf(response.ErrTypeAssertionFail), nil))
	// 	return
	// }

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
	// 绑定数据
	var loginReq requests.LogicReq
	if err := c.ShouldBind(&loginReq); err != nil {
		response.Failed(c, http.StatusBadRequest, response.NewAppErr(globals.StatusBadRequest, fmt.Errorf("Login() -> %v", err), nil))
		return
	}

	// 判断数据是否合法

	// 检验邮箱是否合法
	if !internalUtils.IsValidEmail(loginReq.Email) {
		response.Failed(c, http.StatusBadRequest, response.NewAppErr(globals.StatusBadRequest, fmt.Errorf("Login() : 邮箱不合法"), nil))
		return
	}

	// req, ok := c.Get("req")
	// if !ok {
	// 	globals.Log.Errorf(response.ErrGetReqIsWrong)
	// 	response.Failed(c, http.StatusInternalServerError, response.NewAppErr(globals.StatusInternalServerError, fmt.Errorf(response.ErrGetReqIsWrong), nil))
	// 	return
	// }
	// // 类型断言
	// loginReq, ok := req.(requests.LogicReq)
	// if !ok {
	// 	globals.Log.Errorf(response.ErrTypeAssertionFail)
	// 	response.Failed(c, http.StatusInternalServerError, response.NewAppErr(globals.StatusInternalServerError, fmt.Errorf(response.ErrTypeAssertionFail), nil))
	// 	return
	// }

	// 业务逻辑
	userReqContext := logics.NewUserReqContext(globals.DB, c, globals.SendEmailCfg)
	userInfo, err := userReqContext.Login(loginReq)
	if err != nil {
		globals.Log.Errorf(err.Error())
		response.Failed(c, http.StatusInternalServerError, response.NewAppErr(globals.StatusInternalServerError, err, nil))
		return
	}

	// 通过email查询id
	// user := repositories.QueryUserByEmail(userReqContext.DB, loginReq.Email)
	// if user == nil {
	// 	globals.Log.Errorf(response.ErrEmailNotExist + ":" + loginReq.Email)
	// 	response.Failed(c, http.StatusInternalServerError, response.NewAppErr(globals.StatusInternalServerError, fmt.Errorf(response.ErrEmailNotExist+":"+loginReq.Email), nil))
	// 	return
	// }
	// 生成token
	tok, err := token.GenerateToken(userInfo.Id)
	if err != nil {
		globals.Log.Errorf(err.Error())
		response.Failed(c, http.StatusInternalServerError, response.NewAppErr(globals.StatusInternalServerError, err, nil))
		return
	}
	response.Success(c, http.StatusOK, response.NewAppData(globals.StatusOK, response.DataSuccess, gin.H{"token": tok, "userinfo": userInfo}))
}

// ForgotPassword 忘记密码
func ForgotPassword(c *gin.Context) {
	// 绑定数据
	var forgotPasswordReq requests.ForgotPasswordReq
	if err := c.ShouldBind(&forgotPasswordReq); err != nil {
		response.Failed(c, http.StatusBadRequest, response.NewAppErr(globals.StatusBadRequest, fmt.Errorf("Register() err: %v", err), nil))
		return
	}

	// 判断数据是否合法

	// 检验邮箱是否合法
	if !internalUtils.IsValidEmail(forgotPasswordReq.Email) {
		response.Failed(c, http.StatusBadRequest, response.NewAppErr(globals.StatusBadRequest, fmt.Errorf("Register() : 邮箱不合法"), nil))
		return
	}

	// 核对两次输入的密码
	if forgotPasswordReq.Password != forgotPasswordReq.RePassword {
		response.Failed(c, http.StatusBadRequest, response.NewAppErr(globals.StatusBadRequest, fmt.Errorf("Register() : 两次输入的密码不相同"), nil))
		return
	}

	// 检验密码是否合法
	if !internalUtils.IsValidPassword(forgotPasswordReq.Password) {
		response.Failed(c, http.StatusBadRequest, response.NewAppErr(globals.StatusBadRequest, fmt.Errorf("Register() : 密码必须要同时包含字母、数字、特殊字符，长度在8到20位之间"), nil))
		return
	}

	// req, ok := c.Get("req")
	// if !ok {
	// 	globals.Log.Errorf(response.ErrGetReqIsWrong)
	// 	response.Failed(c, http.StatusInternalServerError, response.NewAppErr(globals.StatusInternalServerError, fmt.Errorf(response.ErrGetReqIsWrong), nil))
	// 	return
	// }
	// // 类型断言
	// forgotPasswordReq, ok := req.(requests.ForgotPasswordReq)
	// if !ok {
	// 	globals.Log.Errorf(response.ErrTypeAssertionFail)
	// 	response.Failed(c, http.StatusInternalServerError, response.NewAppErr(globals.StatusInternalServerError, fmt.Errorf(response.ErrTypeAssertionFail), nil))
	// 	return
	// }

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
