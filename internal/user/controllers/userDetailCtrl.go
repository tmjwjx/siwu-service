package controllers

import (
	"fmt"
	"forum/internal/internalPkg/internalUtils"
	"forum/internal/user/logics"
	"forum/internal/user/requests"
	"forum/pkg/globals"
	"forum/pkg/response"
	"github.com/gin-gonic/gin"
)

// UserDataRequestCtrl 更新用户个人资料
func UserDataRequestCtrl(c *gin.Context) {

	// 获取参数
	var userDataReq requests.UserDataReq
	// 绑定 JSON 数据到结构体
	err := c.ShouldBind(&userDataReq)
	if err != nil {
		// 处理绑定错误
		e := response.NewAppErr(globals.StatusBadRequest, err, nil)
		response.Failed(c, 400, e)
		return
	}

	// 验证参数

	// 验证用户名是否合法
	res := internalUtils.IsValidNickname(userDataReq.Nickname)
	if !res {
		e := response.NewAppErr(globals.StatusBadRequest, fmt.Errorf("UserDataRequestCtrl -> 用户名格式不正确"), nil)
		response.Failed(c, 400, e)
		return
	}

	// 业务处理
	err, status := logics.PersonalDataLogic(&userDataReq, c, globals.DB)

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

// UserDataResponseCtrl 返回用户个人资料给前端
func UserDataResponseCtrl(c *gin.Context) {

	// 获取参数
	userID, exists := c.Get("id")
	if !exists {
		// 返回错误响应
		e := response.NewAppErr(globals.StatusInternalServerError, fmt.Errorf("UserDataResponseCtrl -> 从token中获取用户ID失败"), nil)
		response.Failed(c, 500, e)
		return
	}

	uintValue, err := internalUtils.ChangeAnyToUint(userID)
	if err != nil {
		// 返回错误响应
		e := response.NewAppErr(globals.StatusInternalServerError, err, nil)
		response.Failed(c, 500, e)
		return
	}

	// 业务处理
	userDataRes, err := logics.ResponsePersonDateLogic(uintValue, globals.DB)

	// 返回响应
	if err != nil {
		// 返回错误响应
		e := response.NewAppErr(globals.StatusInternalServerError, err, nil)
		response.Failed(c, 500, e)
	} else {
		// 返回用户数据
		d := response.NewAppData(globals.StatusOK, "用户个人资料响应成功", userDataRes)
		response.Success(c, 200, d)
	}
}

// UserAccountRequestCtrl 更新用户账号设置
func UserAccountRequestCtrl(c *gin.Context) {

	// 获取参数
	var userAccountReq requests.UserAccountReq

	// 绑定 JSON 数据到结构体
	err := c.ShouldBindJSON(&userAccountReq)
	if err != nil {
		// 处理绑定错误
		e := response.NewAppErr(globals.StatusBadRequest, err, nil)
		response.Failed(c, 400, e)
	}

	// 参数验证
	// 验证邮箱是否合法
	res := internalUtils.IsValidEmail(userAccountReq.Email)
	if !res {
		e := response.NewAppErr(globals.StatusBadRequest, fmt.Errorf("邮箱格式不正确"), nil)
		response.Failed(c, 400, e)
		return
	}
	// 验证密码是否合法
	res = internalUtils.IsValidPassword(userAccountReq.Password)
	if !res {
		e := response.NewAppErr(globals.StatusBadRequest, fmt.Errorf("密码格式不正确"), nil)
		response.Failed(c, 400, e)
		return
	}

	// 逻辑处理
	err, status := logics.UserAccountRequestLogic(&userAccountReq, globals.DB)

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
		// 返回成功响应
		d := response.NewAppData(globals.StatusOK, "用户账号设置更新成功", nil)
		response.Success(c, status, d)
	}
}

// UserAccountResponseCtrl 返回用户账号设置数据给前端
func UserAccountResponseCtrl(c *gin.Context) {

	// 获取参数
	userID, exists := c.Get("id")
	if !exists {
		// 返回错误响应
		e := response.NewAppErr(globals.StatusInternalServerError, fmt.Errorf("UserAccountResponseCtrl -> 从token中获取用户ID失败"), nil)
		response.Failed(c, 500, e)
		return
	}

	uintValue, err := internalUtils.ChangeAnyToUint(userID)
	if err != nil {
		// 返回错误响应
		e := response.NewAppErr(globals.StatusInternalServerError, err, nil)
		response.Failed(c, 500, e)
		return
	}

	// 逻辑处理
	userAccountRes, err := logics.UserAccountResponseLogic(uintValue, globals.DB)

	// 返回响应
	if err != nil {
		// 返回错误响应
		e := response.NewAppErr(globals.StatusInternalServerError, err, nil)
		response.Failed(c, 500, e)
	} else {
		// 返回用户数据
		d := response.NewAppData(globals.StatusOK, "用户账号设置数据响应成功", userAccountRes)
		response.Success(c, 200, d)
	}
}

// UserPrivateSetRequestCtrl 更新用户私信设置
func UserPrivateSetRequestCtrl(c *gin.Context) {

	// 获取参数
	var userPrivateSetReq *requests.UserPrivateSettingsReq

	// 绑定 JSON 数据到结构体
	err := c.ShouldBindJSON(&userPrivateSetReq)
	if err != nil {
		// 处理绑定错误
		e := response.NewAppErr(globals.StatusBadRequest, err, nil)
		response.Failed(c, 400, e)
		return
	}

	// 获取参数
	userID, exists := c.Get("id")
	if !exists {
		// 返回错误响应
		e := response.NewAppErr(globals.StatusInternalServerError, fmt.Errorf("UserPrivateSetRequestCtrl -> 从token中获取用户ID失败"), nil)
		response.Failed(c, 500, e)
		return
	}

	uintValue, err := internalUtils.ChangeAnyToUint(userID)
	if err != nil {
		// 返回错误响应
		e := response.NewAppErr(globals.StatusInternalServerError, err, nil)
		response.Failed(c, 500, e)
	}

	err = logics.UserPrivateSetRequestLogic(uintValue, userPrivateSetReq, globals.DB)
	if err != nil {
		// 返回错误响应
		e := response.NewAppErr(globals.StatusInternalServerError, err, nil)
		response.Failed(c, 500, e)
	} else {
		// 返回成功响应
		d := response.NewAppData(globals.StatusOK, "用户私信设置更新成功", nil)
		response.Success(c, 200, d)
	}
}

// UserPrivateSetResponseCtrl 返回用户私信设置数据给前端
func UserPrivateSetResponseCtrl(c *gin.Context) {

	// 获取参数
	userID, exists := c.Get("id")
	if !exists {
		// 返回错误响应
		e := response.NewAppErr(globals.StatusInternalServerError, fmt.Errorf("UserPrivateSetResponseCtrl -> 从token中获取用户ID失败"), nil)
		response.Failed(c, 500, e)
		return
	}

	uintValue, err := internalUtils.ChangeAnyToUint(userID)
	if err != nil {
		// 返回错误响应
		e := response.NewAppErr(globals.StatusInternalServerError, err, nil)
		response.Failed(c, 500, e)
		return
	}

	// 逻辑处理
	userPrivateSetRes, err := logics.UserPrivateSetResponseLogic(uintValue, globals.DB)

	// 返回响应
	if err != nil {
		// 返回错误响应
		e := response.NewAppErr(globals.StatusInternalServerError, err, nil)
		response.Failed(c, 500, e)
	} else {
		// 返回用户数据
		d := response.NewAppData(globals.StatusOK, "用户私信设置数据响应成功", userPrivateSetRes)
		response.Success(c, 200, d)
	}
}
