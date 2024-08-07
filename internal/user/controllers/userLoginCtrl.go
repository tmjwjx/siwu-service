package controllers

import (
	"fmt"
	"forum/internal/internal_pkg/internal_utils"
	"forum/internal/models"
	"forum/internal/user/logics"
	"forum/internal/user/repositories"
	"forum/internal/user/requests"
	"forum/pkg/globals"
	"forum/pkg/response"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

// 注册、登陆、验证码

// Register 用户注册
func Register(c *gin.Context) {
	// 绑定数据
	var registerMsg requests.RegisterMsg
	if err := c.ShouldBind(&registerMsg); err != nil {
		response.Failed(c, http.StatusBadRequest, response.NewAppErr(globals.StatusBadRequest, fmt.Errorf("Register() err: %s", err.Error()), nil))
		return
	}

	// 判断数据是否合法

	// 检验邮箱是否合法
	if !internal_utils.IsValidEmail(registerMsg.Email) {
		response.Failed(c, http.StatusBadRequest, response.NewAppErr(globals.StatusBadRequest, fmt.Errorf("Register() : 邮箱不合法"), nil))
		return
	}

	// 检验密码是否合法
	if !internal_utils.IsValidPassword(registerMsg.Password) {
		response.Failed(c, http.StatusBadRequest, response.NewAppErr(globals.StatusBadRequest, fmt.Errorf("Register() : 密码过于简单"), nil))
		return
	}

	// 核对两次输入的密码
	if registerMsg.Password != registerMsg.RePassword {
		response.Failed(c, http.StatusBadRequest, response.NewAppErr(globals.StatusBadRequest, fmt.Errorf("Register() : 两次输入的密码不相同"), nil))
		return
	}

	// 验证码核对（忽略大小写）
	result, err := repositories.Query(globals.DB, models.UserVerifyCode{}, map[string]interface{}{"email": registerMsg.Email})
	if err != nil {
		response.Failed(c, http.StatusBadRequest, response.NewAppErr(globals.StatusBadRequest, fmt.Errorf("Register() -> %s"), err.Error()))
	}
	// 类型断言，将 interface{} 转换为具体的切片类型
	userVerifyCode, ok := result.(*[]models.UserVerifyCode)
	if !ok {
		response.Failed(c, http.StatusInternalServerError, response.NewAppErr(globals.StatusInternalServerError, fmt.Errorf("Register() : 类型断言错误"), nil))
		return
	}
	verifyCode := (*userVerifyCode)[0].VerifyCode

	fmt.Printf("查询的验证码为 %s，传入的验证码为 %s\n", verifyCode, registerMsg.VerifyCode)
	if !strings.EqualFold(verifyCode, registerMsg.VerifyCode) {
		response.Failed(c, http.StatusBadRequest, response.NewAppErr(globals.StatusBadRequest, fmt.Errorf("Register() : 验证码错误"), nil))
		return
	}

	// 业务逻辑
	userLogic := logics.NewUserLogic(globals.DB, c, globals.VerifyCode)
	err = userLogic.Register(&registerMsg)
	if err != nil {
		response.Failed(c, http.StatusInternalServerError, response.NewAppErr(globals.StatusInternalServerError, fmt.Errorf("Register() -> %s", err.Error()), nil))
		return
	}

	// 成功
	response.Success(c, http.StatusOK, response.NewAppData(globals.StatusOK, "成功", nil))
}

// ReqVerifyCode 用户请求验证码
func ReqVerifyCode(c *gin.Context) {
	userLogic := logics.NewUserLogic(globals.DB, c, globals.VerifyCode)

	// 绑定数据
	var reqVerifyCode requests.ReqVerifyCode
	if err := c.ShouldBind(&reqVerifyCode); err != nil {
		response.Failed(c, http.StatusBadRequest, response.NewAppErr(globals.StatusBadRequest, fmt.Errorf("ReqVerifyCode() -> %s", err.Error()), nil))
		return
	}

	// 检验邮箱是否合法
	if !internal_utils.IsValidEmail(reqVerifyCode.Email) {
		response.Failed(c, http.StatusBadRequest, response.NewAppErr(globals.StatusBadRequest, fmt.Errorf("Register() err: 邮箱不合法"), nil))
		return
	}

	// 业务逻辑
	if err := userLogic.ReqVerifyCode(&reqVerifyCode); err != nil {
		response.Failed(c, http.StatusInternalServerError, response.NewAppErr(globals.StatusInternalServerError, fmt.Errorf("ReqVerifyCode() -> %s", err.Error()), nil))
		return
	}

	// 成功
	response.Success(c, http.StatusOK, response.NewAppData(globals.StatusOK, "成功", nil))
}
