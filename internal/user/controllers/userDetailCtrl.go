package controllers

import (
	"forum/internal/internal_pkg/internal_utils"
	"forum/internal/user/logics"
	"forum/internal/user/requests"
	"forum/pkg/globals"
	"forum/pkg/response"
	"github.com/gin-gonic/gin"
)

// PersonalDataHandler 个人信息的更新处理
func PersonalDataHandler(c *gin.Context) {

	// 获取参数
	var userRequest requests.UserRequest
	// 绑定 JSON 数据到结构体
	err := c.ShouldBindJSON(&userRequest)
	if err != nil {
		// 处理绑定错误
		e := response.NewAppErr(globals.StatusBadRequest, err, nil)
		response.Failed(c, 400, e)
	}

	// 验证参数

	// 验证用户名是否合法
	res := internal_utils.IsValidNickname(userRequest.User.Nickname)
	if !res {
		e := response.NewAppErr(globals.StatusBadRequest, err, nil)
		response.Failed(c, 400, e)
	}
	// 验证邮箱是否合法
	res = internal_utils.IsValidEmail(userRequest.User.Email)
	if !res {
		e := response.NewAppErr(globals.StatusBadRequest, err, nil)
		response.Failed(c, 400, e)
	}
	// 验证密码是否合法
	res = internal_utils.IsValidPassword(userRequest.User.Password)
	if !res {
		e := response.NewAppErr(globals.StatusBadRequest, err, nil)
		response.Failed(c, 400, e)
	}

	// 业务处理
	err, status := logics.PersonalDataLogic(&userRequest, c)

	// 返回响应
	if err != nil {
		var state globals.AppCode
		if status == 400 {
			state = globals.StatusBadRequest
		} else if status == 500 {
			state = globals.StatusInternalServerError
		}
		// 返回错误响应
		e := response.NewAppErr(state, err, nil)
		response.Failed(c, status, e)
	} else {
		d := response.NewAppData(globals.StatusOK, "用户信息更新成功", nil)
		response.Success(c, 200, d)
	}

}

// ResponsePersonDate 返回用户数据给前端
func ResponsePersonDate(c *gin.Context) {

	// 获取参数
	userID := c.Param("id")

	// 业务处理
	userResponse, err := logics.ResponsePersonDateLogic(userID)

	// 返回响应
	if err != nil {
		// 返回错误响应
		e := response.NewAppErr(globals.StatusInternalServerError, err, nil)
		response.Failed(c, 500, e)
	} else {
		// 返回用户数据
		d := response.NewAppData(globals.StatusOK, "用户数据响应成功", userResponse)
		response.Success(c, 200, d)
	}

}
