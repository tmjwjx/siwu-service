package controllers

import (
	"fmt"
	"forum/internal/user/logics"
	"forum/internal/user/requests"
	"forum/pkg/globals"
	"forum/pkg/response"
	"github.com/gin-gonic/gin"
	"net/http"
)

// 注册、登陆、验证码

// Register 用户注册
func Register(c *gin.Context) {
	userLogic := logics.NewUserLogic(globals.DB, c)

	// 绑定数据
	var registerMsg *requests.RegisterMsg
	if err := c.ShouldBind(registerMsg); err != nil {
		response.Failed(c, http.StatusBadRequest, response.NewAppErr(globals.StatusBadRequest, fmt.Errorf("Register() -> %s", err.Error()), nil))
		return
	}

	// 业务逻辑
	err := userLogic.Register(registerMsg)
	if err != nil {
		response.Failed(c, http.StatusInternalServerError, response.NewAppErr(globals.StatusInternalServerError, fmt.Errorf("Register() -> %s", err.Error()), nil))
		return
	}

	// 成功
	response.Success(c, http.StatusOK, response.NewAppData(globals.StatusOK, "成功", nil))
}

// ReqVerifyCode 用户请求验证码
func ReqVerifyCode(c *gin.Context) {
	userLogic := logics.NewUserLogic(globals.DB, c)

	// 绑定数据
	var reqVerifyCode *requests.ReqVerifyCode
	if err := c.ShouldBind(reqVerifyCode); err != nil {
		response.Failed(c, http.StatusBadRequest, response.NewAppErr(globals.StatusBadRequest, fmt.Errorf("ReqVerifyCode() -> %s", err.Error()), nil))
		return
	}

	// 业务逻辑
	err := userLogic.ReqVerifyCode(reqVerifyCode)
	if err != nil {
		response.Failed(c, http.StatusInternalServerError, response.NewAppErr(globals.StatusInternalServerError, fmt.Errorf("ReqVerifyCode() -> %s", err.Error()), nil))
		return
	}

	// 成功
	response.Success(c, http.StatusOK, response.NewAppData(globals.StatusOK, "成功", nil))
}
