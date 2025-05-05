package controllers

import (
	"fmt"
	"forum/internal/menu/logics"
	"forum/internal/menu/requests"
	"forum/pkg/globals"
	"forum/pkg/response"
	"github.com/gin-gonic/gin"
)

// MenuSearchCtrl 检索获取所有菜单列表
func MenuSearchCtrl(c *gin.Context) {

	// 获取参数
	var req requests.MenuSearchReq
	err := c.ShouldBindJSON(&req)
	if err != nil {
		e := response.NewAppErr(globals.StatusBadRequest, fmt.Errorf("MenuSearchCtrl -> 绑定请求结构失败"), nil)
		response.Failed(c, 400, e)
		return
	}

	// 逻辑处理
	menuSearchRes, err := logics.MenuSearchLogic(globals.DB, &req)

	// 返回响应
	if err != nil {
		e := response.NewAppErr(globals.StatusInternalServerError, err, nil)
		response.Failed(c, 500, e)
		return
	}
	d := response.NewAppData(globals.StatusOK, "检索获取所有菜单列表成功", menuSearchRes)
	response.Success(c, 200, d)

}

// GetMenuIconCtrl 获取所有菜单图标
func GetMenuIconCtrl(c *gin.Context) {

	// 逻辑处理
	res, err := logics.GetMenuIconLogic(globals.DB)

	// 返回响应
	if err != nil {
		e := response.NewAppErr(globals.StatusInternalServerError, err, nil)
		response.Failed(c, 500, e)
		return
	}
	d := response.NewAppData(globals.StatusOK, "获取所有菜单图标成功", res)
	response.Success(c, 200, d)

}

// CreateMenuCtrl 新建菜单
func CreateMenuCtrl(c *gin.Context) {

	// 获取参数
	var req requests.CreateMenuReq
	err := c.ShouldBindJSON(&req)
	if err != nil {
		e := response.NewAppErr(globals.StatusBadRequest, fmt.Errorf("CreateMenuCtrl -> 绑定请求结构失败"), nil)
		response.Failed(c, 400, e)
		return
	}

	// 逻辑处理
	err = logics.CreateMenuLogic(globals.DB, &req)

	// 返回响应
	if err != nil {
		e := response.NewAppErr(globals.StatusInternalServerError, err, nil)
		response.Failed(c, 500, e)
		return
	}
	d := response.NewAppData(globals.StatusOK, "新建菜单成功", nil)
	response.Success(c, 200, d)

}

// DeleteMenuCtrl 删除菜单
func DeleteMenuCtrl(c *gin.Context) {

	// 获取参数
	var req requests.DeleteMenuReq
	err := c.ShouldBindJSON(&req)
	if err != nil {
		e := response.NewAppErr(globals.StatusBadRequest, fmt.Errorf("DeleteMenuCtrl -> 绑定请求结构失败"), nil)
		response.Failed(c, 400, e)
		return
	}

	// 逻辑处理
	err = logics.DeleteMenuLogic(globals.DB, &req)

	// 返回响应
	if err != nil {
		e := response.NewAppErr(globals.StatusInternalServerError, err, nil)
		response.Failed(c, 500, e)
		return
	}
	d := response.NewAppData(globals.StatusOK, "删除菜单成功", nil)
	response.Success(c, 200, d)

}

// UpdateMenuCtrl 修改菜单
func UpdateMenuCtrl(c *gin.Context) {

	// 获取参数
	var req requests.UpdateMenuReq
	err := c.ShouldBindJSON(&req)
	if err != nil {
		e := response.NewAppErr(globals.StatusBadRequest, fmt.Errorf("MenuSearchCtrl -> 绑定请求结构失败"), nil)
		response.Failed(c, 400, e)
		return
	}

	// 逻辑处理
	err = logics.UpdateMenuLogic(globals.DB, &req)

	// 返回响应
	if err != nil {
		e := response.NewAppErr(globals.StatusInternalServerError, err, nil)
		response.Failed(c, 500, e)
		return
	}
	d := response.NewAppData(globals.StatusOK, "修改菜单成功", nil)
	response.Success(c, 200, d)

}

// GetMenuDetailCtrl 获取当前菜单详情
func GetMenuDetailCtrl(c *gin.Context) {

	// 获取参数
	id := c.Query("id")

	// 逻辑处理
	res, err := logics.GetMenuDetailLogic(globals.DB, id)

	// 返回响应
	if err != nil {
		e := response.NewAppErr(globals.StatusInternalServerError, err, nil)
		response.Failed(c, 500, e)
		return
	}
	d := response.NewAppData(globals.StatusOK, "获取当前菜单详情成功", res)
	response.Success(c, 200, d)

}

// GetSpecificMenuCtrl
// @Description: 获取所有type为1和2的菜单
// @Author wangyulong 2024-10-15 11:26:56
// @param        c *gin.Context
func GetSpecificMenuCtrl(c *gin.Context) {

	// 逻辑处理
	res, err := logics.GetSpecificMenuLogic(globals.DB)

	// 返回响应
	if err != nil {
		e := response.NewAppErr(globals.StatusInternalServerError, err, nil)
		response.Failed(c, 500, e)
		return
	}
	d := response.NewAppData(globals.StatusOK, "获取所有type为1和2的菜单成功", res)
	response.Success(c, 200, d)

}
