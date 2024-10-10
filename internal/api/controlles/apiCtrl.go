package controlles

import (
	"fmt"
	"forum/internal/api/logics"
	"forum/internal/api/requests"
	"forum/internal/internalPkg/internalUtils"
	"forum/pkg/globals"
	"forum/pkg/response"
	"github.com/gin-gonic/gin"
)

// ApiInitCtrl Api列表初始化
func ApiInitCtrl(c *gin.Context) {

	// 逻辑处理
	res, err := logics.ApiInitLogic(globals.DB)

	// 返回响应
	if err != nil {
		e := response.NewAppErr(globals.StatusInternalServerError, err, nil)
		response.Failed(c, 500, e)
		return
	}
	d := response.NewAppData(globals.StatusOK, "获取所有api列表成功", res)
	response.Success(c, 200, d)

}

// GetApiDetailsCtrl 获取当前api详情
func GetApiDetailsCtrl(c *gin.Context) {
	// 获取参数
	idStr := c.Param("id")
	id, err := internalUtils.ChangeType(idStr)
	if err != nil {
		e := response.NewAppErr(globals.StatusInternalServerError, err, nil)
		response.Failed(c, 500, e)
		return
	}

	// 逻辑处理
	apiDetailsRes, err := logics.GetApiDetailsLogic(globals.DB, id)

	// 返回响应
	if err != nil {
		e := response.NewAppErr(globals.StatusInternalServerError, err, nil)
		response.Failed(c, 500, e)
		return
	}
	d := response.NewAppData(globals.StatusOK, "获取所有api列表成功", apiDetailsRes)
	response.Success(c, 200, d)

}

// GetGroupListCtrl 获取所有api分组列表
func GetGroupListCtrl(c *gin.Context) {

	// 逻辑处理
	apiGroupRes, err := logics.GetGroupListLogic(globals.DB)

	// 返回响应
	if err != nil {
		e := response.NewAppErr(globals.StatusInternalServerError, err, nil)
		response.Failed(c, 500, e)
		return
	}
	d := response.NewAppData(globals.StatusOK, "获取所有api分组列表成功", apiGroupRes)
	response.Success(c, 200, d)

}

// GetRequestMethodCtrl 获取所有请求方法
func GetRequestMethodCtrl(c *gin.Context) {

	// 逻辑处理
	apiReqMethodRes, err := logics.GetRequestMethodLogic(globals.DB)

	// 返回响应
	if err != nil {
		e := response.NewAppErr(globals.StatusInternalServerError, err, nil)
		response.Failed(c, 500, e)
		return
	}
	d := response.NewAppData(globals.StatusOK, "获取所有请求方法成功", apiReqMethodRes)
	response.Success(c, 200, d)

}

// DeleteApiCtrl 删除api
func DeleteApiCtrl(c *gin.Context) {

	// 获取参数
	var req requests.DeleteApiReq
	err := c.ShouldBindJSON(&req)
	if err != nil {
		e := response.NewAppErr(globals.StatusBadRequest, fmt.Errorf("DeleteApiCtrl -> 删除api失败 -> %s", err), nil)
		response.Failed(c, 400, e)
		return
	}

	// 逻辑处理
	err = logics.DeleteApiLogic(&req, globals.DB)
	// 返回响应
	if err != nil {
		e := response.NewAppErr(globals.StatusInternalServerError, err, nil)
		response.Failed(c, 500, e)
		return
	}
	d := response.NewAppData(globals.StatusOK, "删除api成功", nil)
	response.Success(c, 200, d)

}

// UpdateApiCtrl 编辑api
func UpdateApiCtrl(c *gin.Context) {

	// 获取参数
	var req requests.UpdateApiReq
	err := c.ShouldBindJSON(&req)
	if err != nil {
		e := response.NewAppErr(globals.StatusBadRequest, fmt.Errorf("UpdateApiCtrl -> 编辑api失败 -> %s", err), nil)
		response.Failed(c, 400, e)
		return
	}

	// 逻辑处理
	err = logics.UpdateApiLogic(globals.DB, &req)

	// 返回响应
	if err != nil {
		e := response.NewAppErr(globals.StatusInternalServerError, err, nil)
		response.Failed(c, 500, e)
		return
	}
	d := response.NewAppData(globals.StatusOK, "编辑api成功", nil)
	response.Success(c, 200, d)

}

// CreateApiCtrl 添加api
func CreateApiCtrl(c *gin.Context) {

	// 获取参数
	var req requests.CreateApiReq
	err := c.ShouldBindJSON(&req)
	if err != nil {
		e := response.NewAppErr(globals.StatusBadRequest, fmt.Errorf("CreateApiCtrl -> 添加api失败 -> %s", err), nil)
		response.Failed(c, 400, e)
		return
	}

	// 逻辑处理
	err = logics.CreateApiLogic(globals.DB, &req)

	// 返回响应
	if err != nil {
		e := response.NewAppErr(globals.StatusInternalServerError, err, nil)
		response.Failed(c, 500, e)
		return
	}
	d := response.NewAppData(globals.StatusOK, "添加api成功", nil)
	response.Success(c, 200, d)

}

// SearchApiListCtrl 检索api列表
func SearchApiListCtrl(c *gin.Context) {

	// 获取参数
	var req requests.SearchApiListReq
	err := c.ShouldBindJSON(&req)
	if err != nil {
		e := response.NewAppErr(globals.StatusBadRequest, fmt.Errorf("SearchApiListCtrl -> 检索api列表失败 -> %s", err), nil)
		response.Failed(c, 400, e)
		return
	}

	// 逻辑处理
	res, err := logics.SearchApiListLogic(globals.DB, &req)

	// 返回响应
	if err != nil {
		e := response.NewAppErr(globals.StatusInternalServerError, err, nil)
		response.Failed(c, 500, e)
		return
	}
	d := response.NewAppData(globals.StatusOK, "添加api成功", res)
	response.Success(c, 200, d)

}