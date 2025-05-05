package controllers

import (
	"fmt"
	"forum/internal/casbin_r_m_a/logics"
	"forum/internal/casbin_r_m_a/requests"
	"forum/pkg/globals"
	"forum/pkg/response"
	"github.com/gin-gonic/gin"
)

// AssignMenuPermCtrl 为角色分配菜单权限
func AssignMenuPermCtrl(c *gin.Context) {

	// 获取参数
	var req requests.AssignMenuPermReq
	if err := c.ShouldBind(&req); err != nil {
		e := response.NewAppErr(globals.StatusBadRequest, err, nil)
		response.Failed(c, 400, e)
		return
	}

	// 业务处理
	err := logics.AssignMenuPermLogic(globals.DB, &req)

	// 返回响应
	if err != nil {
		e := response.NewAppErr(globals.StatusInternalServerError, err, nil)
		response.Failed(c, 500, e)
		return
	}

	d := response.NewAppData(globals.StatusOK, "为角色分配菜单权限成功", nil)
	response.Success(c, 200, d)

}

// GetMenuPermCtrl 获取当前角色的菜单权限(用于渲染侧边栏，只要type1和2)
func GetMenuPermCtrl(c *gin.Context) {

	// 获取参数
	id := c.Query("id") // 角色ID

	// 业务处理
	res, err := logics.GetMenuPermLogic(globals.DB, id)

	// 返回响应
	if err != nil {
		e := response.NewAppErr(globals.StatusInternalServerError, err, nil)
		response.Failed(c, 500, e)
		return
	}

	d := response.NewAppData(globals.StatusOK, "获取当前角色的菜单权限成功", res)
	response.Success(c, 200, d)

}

// GetApiPermCtrl 获取当前角色的api权限
func GetApiPermCtrl(c *gin.Context) {

	// 获取参数
	id := c.Query("id")

	// 业务处理
	res, err := logics.GetApiPermLogic(globals.CasbinEnforcer, id)

	// 返回响应
	if err != nil {
		e := response.NewAppErr(globals.StatusInternalServerError, err, nil)
		response.Failed(c, 500, e)
		return
	}

	d := response.NewAppData(globals.StatusOK, "获取当前角色的api权限成功", res)
	response.Success(c, 200, d)

}

// AssignApiPermCtrl 为角色分配api权限
func AssignApiPermCtrl(c *gin.Context) {

	// 获取参数
	var req requests.AssignApiPermReq
	err := c.ShouldBindJSON(&req)
	if err != nil {
		e := response.NewAppErr(globals.StatusBadRequest, fmt.Errorf("MenuSearchCtrl -> 绑定请求结构失败"), nil)
		response.Failed(c, 400, e)
		return
	}

	// 逻辑处理
	err = logics.AssignApiPermLogic(globals.DB, &req, globals.CasbinEnforcer)

	// 返回响应
	if err != nil {
		e := response.NewAppErr(globals.StatusInternalServerError, err, nil)
		response.Failed(c, 500, e)
		return
	}
	d := response.NewAppData(globals.StatusOK, "为角色分配api权限成功", nil)
	response.Success(c, 200, d)

}

// GetPermCodeCtrl 获取当前角色的所有权限标识
func GetPermCodeCtrl(c *gin.Context) {

	// 获取参数
	id := c.Query("id")

	// 业务处理
	res, err := logics.GetPermCodeLogic(globals.CasbinEnforcer, globals.DB, id)

	// 返回响应
	if err != nil {
		e := response.NewAppErr(globals.StatusInternalServerError, err, nil)
		response.Failed(c, 500, e)
		return
	}

	d := response.NewAppData(globals.StatusOK, "获取当前角色的菜单权限成功", res)
	response.Success(c, 200, d)

}
