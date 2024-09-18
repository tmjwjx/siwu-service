package controllers

import (
	"fmt"
	"forum/internal/internal_pkg/internal_utils"
	"forum/internal/user/logics"
	"forum/internal/user/requests"
	"forum/pkg/globals"
	"forum/pkg/response"
	"github.com/gin-gonic/gin"
	"net/http"
)

// Reset 重置用户密码
func Reset(c *gin.Context) {
	// 绑定数据
	var reseatReq requests.ReseatReq
	err := c.ShouldBind(&reseatReq)
	if err != nil {
		response.Failed(c, http.StatusBadRequest, response.NewAppErr(globals.StatusBadRequest, fmt.Errorf("Reset() err: %v", err), nil))
		return
	}

	// 业务逻辑
	userReqContext := logics.NewUserReqContext(globals.DB, c, globals.SendEmailCfg)
	if err = userReqContext.Reset(reseatReq); err != nil {
		response.Failed(c, http.StatusInternalServerError, response.NewAppErr(globals.StatusInternalServerError, fmt.Errorf("Reset() -> %v", err), nil))
		return
	}
	response.Success(c, http.StatusOK, response.NewAppData(globals.StatusOK, "成功", nil))
}

// Add 添加用户
func Add(c *gin.Context) {
	// 绑定数据
	var addReq requests.AddAndEditReq
	err := c.ShouldBind(&addReq)
	if err != nil {
		response.Failed(c, http.StatusBadRequest, response.NewAppErr(globals.StatusBadRequest, fmt.Errorf("Add() err: %v", err), nil))
		return
	}

	// 数据检验
	// 检验邮箱是否合法
	if !internal_utils.IsValidEmail(addReq.Email) {
		response.Failed(c, http.StatusBadRequest, response.NewAppErr(globals.StatusBadRequest, fmt.Errorf("Add() err: 邮箱不合法"), nil))
		return
	}

	// 业务逻辑
	userReqContext := logics.NewUserReqContext(globals.DB, c, globals.SendEmailCfg)
	if err = userReqContext.Add(addReq); err != nil {
		response.Failed(c, http.StatusInternalServerError, response.NewAppErr(globals.StatusInternalServerError, fmt.Errorf("Add() -> %v", err), nil))
		return
	}
	response.Success(c, http.StatusOK, response.NewAppData(globals.StatusOK, "成功", nil))
}

// Delete 删除用户
func Delete(c *gin.Context) {
	// 绑定数据
	var deleteReq requests.DeleteReq
	err := c.ShouldBind(&deleteReq)
	if err != nil {
		response.Failed(c, http.StatusBadRequest, response.NewAppErr(globals.StatusBadRequest, fmt.Errorf("Delete() err: %v", err), nil))
		return
	}

	// 业务逻辑
	userReqContext := logics.NewUserReqContext(globals.DB, c, globals.SendEmailCfg)
	if err = userReqContext.Delete(deleteReq); err != nil {
		response.Failed(c, http.StatusInternalServerError, response.NewAppErr(globals.StatusInternalServerError, fmt.Errorf("Delete() -> %v", err), nil))
		return
	}
	response.Success(c, http.StatusOK, response.NewAppData(globals.StatusOK, "成功", nil))
}

// Edit 编辑用户
func Edit(c *gin.Context) {
	// 绑定数据
	var editReq requests.AddAndEditReq
	err := c.ShouldBind(&editReq)
	if err != nil {
		response.Failed(c, http.StatusBadRequest, response.NewAppErr(globals.StatusBadRequest, fmt.Errorf("Edit() err: %v", err), nil))
		return
	}

	// 业务逻辑
	userReqContext := logics.NewUserReqContext(globals.DB, c, globals.SendEmailCfg)
	if err := userReqContext.Edit(editReq); err != nil {
		response.Failed(c, http.StatusInternalServerError, response.NewAppErr(globals.StatusInternalServerError, fmt.Errorf("Edit() -> %v", err), nil))
		return
	}

	response.Success(c, http.StatusOK, response.NewAppData(globals.StatusOK, "成功", nil))
}

// List 获取所有用户列表
func List(c *gin.Context) {
	// 绑定数据
	var listReq requests.ListReq
	err := c.ShouldBind(&listReq)
	if err != nil {
		response.Failed(c, http.StatusBadRequest, response.NewAppErr(globals.StatusBadRequest, fmt.Errorf("List() err: %v", err), nil))
		return
	}

	// 检验数据
	if listReq.Page <= 0 {
		response.Failed(c, http.StatusBadRequest, response.NewAppErr(globals.StatusBadRequest, fmt.Errorf("List() err: Page参数必须为正数"), nil))
		return
	}
	if listReq.Limit <= 0 {
		response.Failed(c, http.StatusBadRequest, response.NewAppErr(globals.StatusBadRequest, fmt.Errorf("List() err: limit参数必须为正数"), nil))
		return
	}

	// 业务逻辑
	userReqContext := logics.NewUserReqContext(globals.DB, c, globals.SendEmailCfg)
	listRes, err := userReqContext.List(listReq)
	if err != nil {
		response.Failed(c, http.StatusInternalServerError, response.NewAppErr(globals.StatusInternalServerError, fmt.Errorf("List() -> %v", err), nil))
		return
	}

	response.Success(c, http.StatusOK, response.NewAppData(globals.StatusOK, "成功", gin.H{"user_list": listRes}))
}

// Import 导入用户表
func Import(c *gin.Context) {

}

// Export 导出用户表
func Export(c *gin.Context) {

}

// ImportTemplate 下载导入用户模版excel
func ImportTemplate(c *gin.Context) {

}

// GetInfo 获取当前用户基本信息
func GetInfo(c *gin.Context) {

}
