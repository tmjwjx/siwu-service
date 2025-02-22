package controllers

import (
	"errors"
	"forum/internal/administrator/logics"
	"forum/internal/administrator/requests"
	"forum/pkg/globals"
	"forum/pkg/response"
	"github.com/gin-gonic/gin"
	"net/http"
)

/*
	// Ctrl模板
	// 初始化需要的变量
	db := globals.DB
	var req *requests.XXXReq

	// 绑定查询参数到 req 变量，如果绑定失败，返回错误信息
	if err := c.ShouldBind(&req); err != nil {
		// 日志记录错误信息
		globals.Log.Errorf("绑定req失败 err = %s", err)
		// 返回错误响应
		data := response.NewAppErr(globals.StatusBadRequest, err, nil)
		response.Failed(c, http.StatusBadRequest, data)
		return // 结束函数执行
	}

	// 进入业务层
	articleList, err := logics.XXXLogic(db, req)
	if err != nil {
		globals.Log.Errorf("获取数据失败 err = %s", err)
		data := response.NewAppErr(globals.StatusInternalServerError, err, nil)
		response.Failed(c, http.StatusInternalServerError, data)
		return
	}

	// 返回响应
	data := response.NewAppData(globals.StatusOK, "成功", articleList)
	response.Success(c, http.StatusOK, data)
*/

// UpdateAdministratorCtrl
// @Description: 编辑管理员信息
// @param        c *gin.Context
// @Author tianjiajie 2025-02-21 17:42:20
func UpdateAdministratorCtrl(c *gin.Context) {
	// 初始化需要的变量
	db := globals.DB
	req := requests.UpdateAdministratorReq{}
	id, ok := c.Get("id")
	globals.Log.Info(id)
	if !ok {
		globals.Log.Errorf("身份验证失败")
		data := response.NewAppErr(globals.StatusBadRequest, errors.New("获取id失败"), nil)
		response.Failed(c, http.StatusBadRequest, data)
		return
	}

	// 绑定查询参数到 req 变量
	if err := c.ShouldBind(&req); err != nil {
		// 日志记录错误信息
		globals.Log.Errorf("绑定req失败 err = %s", err)
		// 返回错误响应
		data := response.NewAppErr(globals.StatusBadRequest, err, nil)
		response.Failed(c, http.StatusBadRequest, data)
		return // 结束函数执行
	}

	// 进入业务层
	err := logics.UpdateAdministratorLogic(db, req, id)
	if err != nil {
		globals.Log.Errorf("获取管理员信息失败 err = %s", err)
		data := response.NewAppErr(globals.StatusInternalServerError, err, nil)
		response.Failed(c, http.StatusInternalServerError, data)
		return
	}

	// 返回响应
	data := response.NewAppData(globals.StatusOK, "成功", nil)
	response.Success(c, http.StatusOK, data)

}

// ResetAdministratorPasswordCtrl
// @Description: 重置管理员密码
// @param        c *gin.Context
// @Author tianjiajie 2025-02-21 21:27:59
func ResetAdministratorPasswordCtrl(c *gin.Context) {
	// 初始化需要的变量
	db := globals.DB
	req := requests.AdministratorReq{}

	// 绑定查询参数到 req 变量
	if err := c.ShouldBind(&req); err != nil {
		// 日志记录错误信息
		globals.Log.Errorf("绑定req失败 err = %s", err)
		// 返回错误响应
		data := response.NewAppErr(globals.StatusBadRequest, err, nil)
		response.Failed(c, http.StatusBadRequest, data)
		return // 结束函数执行
	}

	// 进入业务层
	err := logics.ResetAdministratorPasswordLogic(db, req)
	if err != nil {
		globals.Log.Errorf("获取管理员信息失败 err = %s", err)
		data := response.NewAppErr(globals.StatusInternalServerError, err, nil)
		response.Failed(c, http.StatusInternalServerError, data)
		return
	}

	// 返回响应
	data := response.NewAppData(globals.StatusOK, "成功", nil)
	response.Success(c, http.StatusOK, data)

}

// GetAdministratorInfoCtrl
// @Description: 查询管理员详情
// @param        c *gin.Context
// @Author tianjiajie 2025-02-21 21:19:57
func GetAdministratorInfoCtrl(c *gin.Context) {
	// 初始化需要的变量
	db := globals.DB

	// 绑定查询参数到 req 变量
	id := c.Query("id")

	// 进入业务层
	administratorInfo, err := logics.GetAdministratorInfoLogic(db, id)
	if err != nil {
		globals.Log.Errorf("获取管理员信息失败 err = %s", err)
		data := response.NewAppErr(globals.StatusInternalServerError, err, nil)
		response.Failed(c, http.StatusInternalServerError, data)
		return
	}

	// 返回响应
	data := response.NewAppData(globals.StatusOK, "成功", administratorInfo)
	response.Success(c, http.StatusOK, data)
}

// GetAdministratorListCtrl
// @Description: 查询管理员列表
// @param        c *gin.Context
// @Author tianjiajie 2025-02-21 20:30:46
func GetAdministratorListCtrl(c *gin.Context) {
	// 初始化需要的变量
	db := globals.DB
	req := requests.GetAdministratorListReq{}

	// 绑定查询参数到 req 变量
	if err := c.ShouldBind(&req); err != nil {
		// 日志记录错误信息
		globals.Log.Errorf("绑定req失败 err = %s", err)
		// 返回错误响应
		data := response.NewAppErr(globals.StatusBadRequest, err, nil)
		response.Failed(c, http.StatusBadRequest, data)
		return // 结束函数执行
	}

	// 进入业务层
	administratorList, err := logics.GetAdministratorListLogic(db, req)
	if err != nil {
		globals.Log.Errorf("获取管理员列表失败 err = %s", err)
		data := response.NewAppErr(globals.StatusInternalServerError, err, nil)
		response.Failed(c, http.StatusInternalServerError, data)
		return
	}

	// 返回响应
	data := response.NewAppData(globals.StatusOK, "成功", administratorList)
	response.Success(c, http.StatusOK, data)
}

// DeleteAdministratorCtrl
// @Description: 删除管理员
// @param        c *gin.Context
// @Author tianjiajie 2025-02-21 14:59:46
func DeleteAdministratorCtrl(c *gin.Context) {
	// 初始化需要的变量
	db := globals.DB
	req := requests.AdministratorReq{}
	id, ok := c.Get("id")
	globals.Log.Info(id)
	if !ok {
		globals.Log.Errorf("身份验证失败")
		data := response.NewAppErr(globals.StatusBadRequest, errors.New("获取id失败"), nil)
		response.Failed(c, http.StatusBadRequest, data)
		return
	}
	//else if id != 1 { // 是否写死 待定
	//	globals.Log.Errorf("没有权限")
	//	data := response.NewAppErr(globals.StatusBadRequest, errors.New("没有权限"), nil)
	//	response.Failed(c, http.StatusBadRequest, data)
	//	return
	//}

	// 绑定查询参数到 req 变量
	if err := c.ShouldBind(&req); err != nil {
		// 日志记录错误信息
		globals.Log.Errorf("绑定req失败 err = %s", err)
		// 返回错误响应
		data := response.NewAppErr(globals.StatusBadRequest, err, nil)
		response.Failed(c, http.StatusBadRequest, data)
		return // 结束函数执行
	}

	// 进入业务层
	err := logics.DeleteAdministratorLogic(db, req)
	if err != nil {
		globals.Log.Errorf("删除管理员失败 err = %s", err)
		data := response.NewAppErr(globals.StatusInternalServerError, err, nil)
		response.Failed(c, http.StatusInternalServerError, data)
		return
	}

	// 返回响应
	data := response.NewAppData(globals.StatusOK, "成功", nil)
	response.Success(c, http.StatusOK, data)
}

// AddAdministratorCtrl
// @Description: 添加管理员
// @param        c *gin.Context
// @Author tianjiajie 2025-02-20 20:11:54
func AddAdministratorCtrl(c *gin.Context) {
	// 初始化需要的变量
	db := globals.DB
	req := requests.AddAdministratorReq{}
	id, ok := c.Get("id")
	globals.Log.Info(id)
	if !ok {
		globals.Log.Errorf("身份验证失败")
		data := response.NewAppErr(globals.StatusBadRequest, errors.New("获取id失败"), nil)
		response.Failed(c, http.StatusBadRequest, data)
		return
	}
	//else if id != 1 { // 是否写死 待定
	//	globals.Log.Errorf("没有权限")
	//	data := response.NewAppErr(globals.StatusBadRequest, errors.New("没有权限"), nil)
	//	response.Failed(c, http.StatusBadRequest, data)
	//	return
	//}

	// 绑定查询参数到 req 变量
	if err := c.ShouldBind(&req); err != nil {
		// 日志记录错误信息
		globals.Log.Errorf("绑定req失败 err = %s", err)
		// 返回错误响应
		data := response.NewAppErr(globals.StatusBadRequest, err, nil)
		response.Failed(c, http.StatusBadRequest, data)
		return // 结束函数执行
	}

	// 进入业务层
	administrator, err := logics.AddAdministratorLogic(db, req)
	if err != nil {
		globals.Log.Errorf("创建管理员失败 err = %s", err)
		data := response.NewAppErr(globals.StatusInternalServerError, err, nil)
		response.Failed(c, http.StatusInternalServerError, data)
		return
	}

	// 返回响应
	data := response.NewAppData(globals.StatusOK, "成功", administrator)
	response.Success(c, http.StatusOK, data)

}
