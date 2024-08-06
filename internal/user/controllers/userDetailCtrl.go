package controllers

import (
	"forum/internal/pkg/utils"
	"forum/internal/user/logics"
	"forum/internal/user/requests"
	"forum/pkg/globals"
	"forum/pkg/response"
	"github.com/gin-gonic/gin"
)

// PersonalDataHandler 个人信息的更新处理
func PersonalDataHandler(c *gin.Context) {

	// 获取参数
	var user requests.User
	err := c.ShouldBind(&user)
	if err != nil {
		// 处理绑定错误
		e := response.NewAppErr(globals.StatusBadRequest, err, nil)
		response.Failed(c, 400, e)
	}

	// 验证参数
	err = utils.UserDateVerify(&user)
	if err != nil {
		e := response.NewAppErr(globals.StatusBadRequest, err, nil)
		response.Failed(c, 400, e)
	}

	// 业务处理
	err, status := logics.PersonalDataLogic(&user, c)

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

// ResponsePersonDate 返回用户数据前端
func ResponsePersonDate(c *gin.Context) {

	// 获取参数
	userID := c.Param("id")

	// 业务处理
	user, err := logics.ResponsePersonDateLogic(userID, c)

	// 返回响应
	if err != nil {
		// 返回错误响应
		e := response.NewAppErr(globals.StatusInternalServerError, err, nil)
		response.Failed(c, 500, e)
	} else {
		// 返回用户数据
		d := response.NewAppData(globals.StatusOK, "用户数据响应成功", user)
		response.Success(c, 200, d)
	}

}
