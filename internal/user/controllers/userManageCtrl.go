package controllers

import (
	"fmt"
	"forum/internal/internalPkg/internalUtils"
	"forum/internal/user/logics"
	"forum/internal/user/requests"
	"forum/pkg/globals"
	"forum/pkg/response"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
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
	var addReq requests.AddReq
	err := c.ShouldBind(&addReq)
	if err != nil {
		response.Failed(c, http.StatusBadRequest, response.NewAppErr(globals.StatusBadRequest, fmt.Errorf("Add() err: %v", err), nil))
		return
	}

	// 数据检验
	// 检验邮箱是否合法
	if !internalUtils.IsValidEmail(addReq.Email) {
		response.Failed(c, http.StatusBadRequest, response.NewAppErr(globals.StatusBadRequest, fmt.Errorf("Add() err: 邮箱不合法"), nil))
		return
	}

	// 业务逻辑
	userReqContext := logics.NewUserReqContext(globals.DB, c, globals.SendEmailCfg)
	id, err := userReqContext.Add(addReq)
	if err != nil {
		response.Failed(c, http.StatusInternalServerError, response.NewAppErr(globals.StatusInternalServerError, fmt.Errorf("Add() -> %v", err), nil))
		return
	}
	response.Success(c, http.StatusOK, response.NewAppData(globals.StatusOK, "成功", gin.H{"id": id}))
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
	var editReq requests.EditReq
	err := c.ShouldBind(&editReq)
	if err != nil {
		response.Failed(c, http.StatusBadRequest, response.NewAppErr(globals.StatusBadRequest, fmt.Errorf("Edit() err: %v", err), nil))
		return
	}

	// 数据校验
	// 判断邮箱是否合法
	if !internalUtils.IsValidEmail(editReq.Email) {
		response.Failed(c, http.StatusBadRequest, response.NewAppErr(globals.StatusBadRequest, fmt.Errorf("Edit() err: 邮箱不合法"), nil))
		return
	}

	// 业务逻辑
	userReqContext := logics.NewUserReqContext(globals.DB, c, globals.SendEmailCfg)
	if err = userReqContext.Edit(editReq); err != nil {
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
	listRes, total, err := userReqContext.List(listReq)
	if err != nil {
		response.Failed(c, http.StatusInternalServerError, response.NewAppErr(globals.StatusInternalServerError, fmt.Errorf("List() -> %v", err), nil))
		return
	}

	response.Success(c, http.StatusOK, response.NewAppData(globals.StatusOK, "成功", gin.H{"user_list": listRes, "total": total}))
}

// Import 导入用户表
func Import(c *gin.Context) {
	// 获取上传的文件
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to get file"})
		return
	}

	// 业务逻辑
	userReqContext := logics.NewUserReqContext(globals.DB, c, globals.SendEmailCfg)
	if err = userReqContext.Import(file); err != nil {
		response.Failed(c, http.StatusInternalServerError, response.NewAppErr(globals.StatusInternalServerError, fmt.Errorf("Export() err: %v", err), nil))
		return
	}
	response.Success(c, http.StatusOK, response.NewAppData(globals.StatusOK, "成功", nil))
}

// Export 导出用户表
func Export(c *gin.Context) {

	// 设置响应头，返回 Excel 文件
	fileName := "users.xls"
	c.Header("Content-Disposition", "attachment; filename="+fileName)
	c.Header("Content-Type", "application/vnd.ms-excel")
	c.Header("Content-Transfer-Encoding", "binary")

	// 业务逻辑
	userReqContext := logics.NewUserReqContext(globals.DB, c, globals.SendEmailCfg)
	if err := userReqContext.Export(); err != nil {
		response.Failed(c, http.StatusInternalServerError, response.NewAppErr(globals.StatusInternalServerError, fmt.Errorf("Export() err: %v", err), nil))
		return
	}
	// response.Success(c, http.StatusOK, response.NewAppData(globals.StatusOK, "成功", nil))
}

// DownloadTemplate 下载导入用户模版excel
func DownloadTemplate(c *gin.Context) {
	// 设置响应头，返回 Excel 文件
	c.Header("Content-Disposition", "attachment; filename=user_template.xlsx")
	c.Header("Content-Type", "application/octet-stream")
	c.Header("Content-Transfer-Encoding", "binary")

	// 业务逻辑
	userReqContext := logics.NewUserReqContext(globals.DB, c, globals.SendEmailCfg)
	if err := userReqContext.DownloadTemplate(); err != nil {
		response.Failed(c, http.StatusInternalServerError, response.NewAppErr(globals.StatusInternalServerError, fmt.Errorf("DownloadTemplate() err: %v", err), nil))
		return
	}
	response.Success(c, http.StatusOK, response.NewAppData(globals.StatusOK, "成功", nil))
}

// GetInfo 获取当前用户基本信息
func GetInfo(c *gin.Context) {
	// 获取参数
	id, err := strconv.Atoi(c.Query("id"))
	if err != nil {
		response.Failed(c, http.StatusBadRequest, response.NewAppErr(globals.StatusBadRequest, fmt.Errorf("GetDetail() err: 数据错误"), nil))
		return
	}

	// 业务逻辑
	userReqContext := logics.NewUserReqContext(globals.DB, c, globals.SendEmailCfg)
	info, err := userReqContext.GetInfo(uint(id))
	if err != nil {
		response.Failed(c, http.StatusInternalServerError, response.NewAppErr(globals.StatusInternalServerError, fmt.Errorf("GetInfo() -> %v", err), nil))
		return
	}
	response.Success(c, http.StatusOK, response.NewAppData(globals.StatusOK, "成功", info))
}
