package controllers

import (
	"fmt"
	"forum/internal/roleManage/logics"
	"forum/internal/roleManage/requests"
	"forum/pkg/globals"
	"forum/pkg/response"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

// AddRole 添加角色
func AddRole(c *gin.Context) {
	// 绑定数据
	var role requests.RoleReq
	if err := c.ShouldBind(&role); err != nil {
		response.Failed(c, http.StatusBadRequest, response.NewAppErr(globals.StatusBadRequest, fmt.Errorf("AddRole() err: 绑定数据错误"), nil))
		return
	}

	// 业务逻辑
	roleReqContext := logics.NewRoleReqContext(globals.DB, c)
	if err := roleReqContext.AddRole(role); err != nil {
		response.Failed(c, http.StatusInternalServerError, response.NewAppErr(globals.StatusInternalServerError, fmt.Errorf("AddRole() -> %v", err), nil))
		return
	}

	response.Success(c, http.StatusOK, response.NewAppData(globals.StatusOK, "成功", nil))
}

// DeleteRole 删除角色
func DeleteRole(c *gin.Context) {
	// 绑定数据
	var ids requests.RoleIdsReq
	if err := c.ShouldBind(&ids); err != nil {
		response.Failed(c, http.StatusBadRequest, response.NewAppErr(globals.StatusBadRequest, fmt.Errorf("DeleteRole() err: 绑定数据错误"), nil))
		return
	}

	// 业务逻辑
	roleReqContext := logics.NewRoleReqContext(globals.DB, c)
	if err := roleReqContext.DeleteRole(ids); err != nil {
		response.Failed(c, http.StatusInternalServerError, response.NewAppErr(globals.StatusInternalServerError, fmt.Errorf("DeleteRole() -> %v", err), nil))
		return
	}

	response.Success(c, http.StatusOK, response.NewAppData(globals.StatusOK, "成功", nil))
}

// SearchRole 检索角色
func SearchRole(c *gin.Context) {
	// 绑定数据
	var searchRole requests.SearchRoleReq
	if err := c.ShouldBind(&searchRole); err != nil {
		response.Failed(c, http.StatusBadRequest, response.NewAppErr(globals.StatusBadRequest, fmt.Errorf("SearchRoleReq() err: 绑定数据错误"), nil))
		return
	}

	// 检验数据
	if searchRole.Page <= 0 {
		response.Failed(c, http.StatusBadRequest, response.NewAppErr(globals.StatusBadRequest, fmt.Errorf("SearchRoleReq() err: Page参数必须为正数"), nil))
		return
	}
	if searchRole.Limit <= 0 {
		response.Failed(c, http.StatusBadRequest, response.NewAppErr(globals.StatusBadRequest, fmt.Errorf("SearchRoleReq() err: limit参数必须为正数"), nil))
		return
	}

	// 业务逻辑
	roleReqContext := logics.NewRoleReqContext(globals.DB, c)
	roles, err := roleReqContext.SearchRole(searchRole)
	if err != nil {
		response.Failed(c, http.StatusInternalServerError, response.NewAppErr(globals.StatusInternalServerError, fmt.Errorf("SearchRoleReq() -> %v", err), nil))
		return
	}

	response.Success(c, http.StatusOK, response.NewAppData(globals.StatusOK, "成功", gin.H{"role_list": roles}))
}

// UpdateRole 更新角色
func UpdateRole(c *gin.Context) {
	// 绑定数据
	var role requests.RoleReq
	if err := c.ShouldBind(&role); err != nil {
		response.Failed(c, http.StatusBadRequest, response.NewAppErr(globals.StatusBadRequest, fmt.Errorf("UpdateRole() err: 绑定数据错误"), nil))
		return
	}

	// 业务逻辑
	roleReqContext := logics.NewRoleReqContext(globals.DB, c)
	if err := roleReqContext.UpdateRole(role); err != nil {
		response.Failed(c, http.StatusInternalServerError, response.NewAppErr(globals.StatusInternalServerError, fmt.Errorf("UpdateRole() -> %v", err), nil))
		return
	}

	response.Success(c, http.StatusOK, response.NewAppData(globals.StatusOK, "成功", nil))
}

// GetRoleName 获取所有已启用的角色名称列表
func GetRoleName(c *gin.Context) {
	// status: 0全部 1正常 2封禁
	status := 1

	// 业务逻辑
	roleReqContext := logics.NewRoleReqContext(globals.DB, c)
	roles, err := roleReqContext.GetRoleName(status)
	if err != nil {
		response.Failed(c, http.StatusInternalServerError, response.NewAppErr(globals.StatusInternalServerError, fmt.Errorf("SearchRoleReq() -> %v", err), nil))
		return
	}

	response.Success(c, http.StatusOK, response.NewAppData(globals.StatusOK, "成功", gin.H{"role_list": roles}))
}

// GetDetail 获取当前角色详情
func GetDetail(c *gin.Context) {
	// 查询参数 Query
	id, err := strconv.Atoi(c.Query("id"))
	if err != nil {
		response.Failed(c, http.StatusBadRequest, response.NewAppErr(globals.StatusBadRequest, fmt.Errorf("GetDetail() err: 数据错误"), nil))
		return
	}

	// 业务逻辑
	roleReqContext := logics.NewRoleReqContext(globals.DB, c)
	roles, err := roleReqContext.GetDetail(uint(id))
	if err != nil {
		response.Failed(c, http.StatusInternalServerError, response.NewAppErr(globals.StatusInternalServerError, fmt.Errorf("GetDetail() -> %v", err), nil))
		return
	}

	response.Success(c, http.StatusOK, response.NewAppData(globals.StatusOK, "成功", roles))
}

// DispatchRole 给用户分配角色
func DispatchRole(c *gin.Context) {
	// 绑定数据
	var dispatchRole requests.DispatchRoleReq
	if err := c.ShouldBind(&dispatchRole); err != nil {
		response.Failed(c, http.StatusBadRequest, response.NewAppErr(globals.StatusBadRequest, fmt.Errorf("DispatchRoleReq() err: 绑定数据错误"), nil))
		return
	}

	// 业务逻辑
	roleReqContext := logics.NewRoleReqContext(globals.DB, c)
	if err := roleReqContext.DispatchRole(dispatchRole); err != nil {
		response.Failed(c, http.StatusInternalServerError, response.NewAppErr(globals.StatusInternalServerError, fmt.Errorf("DispatchRoleReq() -> %v", err), nil))
		return
	}

	response.Success(c, http.StatusOK, response.NewAppData(globals.StatusOK, "成功", nil))
}
