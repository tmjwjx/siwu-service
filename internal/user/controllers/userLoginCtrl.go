package controllers

import (
	"encoding/json"
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

	// 获取参数，检验参数

	// 获取数据包
	data, err := c.GetRawData()
	// 请求语法错误或无效参数
	if err != nil {
		response.Failed(c, http.StatusBadRequest, response.NewAppErr(globals.StatusBadRequest, err, nil))
		return
	}

	// fmt.Println("Register 接收到的消息为:")
	// fmt.Println(string(data))

	// 将数据包反序列化
	var registerMsg requests.RegisterMsg
	err = json.Unmarshal(data, &registerMsg)
	if err != nil {
		response.Failed(c, http.StatusInternalServerError, response.NewAppErr(globals.StatusInternalServerError, err, nil))
		return
	}

	// 具体的业务逻辑
	appErr := userLogic.Register(registerMsg)
	if appErr != nil {
		response.Failed(c, http.StatusInternalServerError, response.NewAppErr(globals.StatusInternalServerError, err, nil))
		return
	}

	// 成功
	response.Success(c, http.StatusOK, response.NewAppData(globals.StatusOK, "", nil))
}

// ReqVerifyCode 用户请求验证码
func ReqVerifyCode(c *gin.Context) {
	// userLogic := logic.NewUserLogic(globals.DB, c)
	//
	// // 具体的业务逻辑
	// verifycode, err := userLogic.ReqVerifyCode()
	// if err != nil {
	// 	response.Failed(c, http.StatusInternalServerError, response.NewAppErr(globals.StatusInternalServerError, err, nil))
	// 	return
	// }
	//
	// // 成功
	// response.Success(c, http.StatusOK, response.NewAppData(globals.StatusOK, ""))
}
