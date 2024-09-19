package logics

import (
	"fmt"
	"forum/internal/models"
	"forum/internal/roleManage/repositories"
	"forum/internal/roleManage/requests"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// RoleReqContext 用于在处理请求时传递数据库连接和请求上下文信息
type RoleReqContext struct {
	DB  *gorm.DB
	Ctx *gin.Context
}

func NewRoleReqContext(db *gorm.DB, c *gin.Context) *RoleReqContext {
	return &RoleReqContext{
		DB:  db,
		Ctx: c,
	}
}

// AddRole 添加角色
func (r *RoleReqContext) AddRole(role requests.RoleReq) error {
	err := repositories.InsertObject(r.DB, &models.Role{Name: role.Name, Code: role.Code, Status: role.Status, Sort: role.Sort})
	if err != nil {
		return fmt.Errorf("RoleReqContext.AddRole() -> %v", err)
	}
	return nil
}

// DeleteRole 删除角色
func (r *RoleReqContext) DeleteRole(ids requests.RoleIdsReq) error {
	var num int64
	for _, v := range ids.Ids {
		n, err := repositories.DeleteObjectsByModel(r.DB, &models.Role{}, map[string]interface{}{"id": v})
		if err != nil {
			return fmt.Errorf("RoleReqContext.DeleteRole() -> %v", err)
		}
		num += n
	}
	// fmt.Printf("RoleReqContext.DeleteRole() 应该删除%v个角色，实际删除%v个角色\n", len(ids.Ids), num)

	return nil
}

// SearchRole 检索数据
func (r *RoleReqContext) SearchRole(searchRole requests.SearchRoleReq) ([]*models.Role, error) {
	conditions := map[string]interface{}{}

	// 判断是否该添加某些查询条件（如果某些条件为空，那么就不查询这个条件）
	if searchRole.Name != "" {
		conditions["name"] = searchRole.Name
	}
	if searchRole.Code != "" {
		conditions["code"] = searchRole.Code
	}
	// status == 0 代表着全部
	if searchRole.Status != 0 {
		conditions["status"] = searchRole.Status
	}

	// 查询
	roleSli, err := repositories.QueryRolesByPage(r.DB, conditions, searchRole.Page, searchRole.Limit)
	if err != nil {
		return nil, fmt.Errorf("RoleReqContext.SearchRoleReq() -> %v", err)
	}

	return roleSli, nil
}

// UpdateRole 更新角色
func (r *RoleReqContext) UpdateRole(role requests.RoleReq) error {
	// 判断该角色id是否存在
	role2 := repositories.QueryRoleById(r.DB, role.Id)
	if role2 == nil {
		return fmt.Errorf("RoleReqContext.UpdateRole err = 没有查询到id为 %v 的角色", role.Id)
	}

	m := map[string]interface{}{
		"name":   role.Name,
		"code":   role.Code,
		"status": role.Status,
		"sort":   role.Sort,
	}
	err := repositories.UpdateObjects(r.DB, &models.Role{Model: gorm.Model{ID: role.Id}}, m)
	if err != nil {
		return fmt.Errorf("RoleReqContext.UpdateRole() -> %v", err)
	}
	return nil
}

// GetRoleName 获取所有已启用的角色名称列表
func (r *RoleReqContext) GetRoleName(status int) ([]*models.Role, error) {
	// 查询
	roleSli, err := repositories.QueryRoles(r.DB, map[string]interface{}{"status": status})
	if err != nil {
		return nil, fmt.Errorf("RoleReqContext.GetRoleName() -> %v", err)
	}

	return roleSli, nil
}

// GetDetail 获取当前角色详情
func (r *RoleReqContext) GetDetail(id uint) (*models.Role, error) {
	// 查询
	role := repositories.QueryRoleById(r.DB, id)
	if role == nil {
		return nil, fmt.Errorf("RoleReqContext.GetDetail err = 没有查询到id为 %v 的角色", id)
	}
	return role, nil
}

// DispatchRole 给用户分配角色
func (r *RoleReqContext) DispatchRole(dispatchRole requests.DispatchRoleReq) error {
	// 判断数据是否正确
	// 判断userId是否存在
	user := repositories.QueryUserById(r.DB, dispatchRole.UserId)
	if user == nil {
		return fmt.Errorf("RoleReqContext.DispatchRoleReq() err = 没有查询到id为 %v 的用户", dispatchRole.UserId)
	}

	adminRoleSli := make([]*models.AdminRole, 0)
	userId := dispatchRole.UserId
	for _, id := range dispatchRole.Ids {
		// 判断角色Id是否存在
		role := repositories.QueryRoleById(r.DB, id)
		if role == nil {
			return fmt.Errorf("RoleReqContext.DispatchRoleReq err = 没有查询到id为 %v 的角色", id)
		}

		// 添加的切片里
		adminRoleSli = append(adminRoleSli, &models.AdminRole{
			AdminId: userId,
			RoleId:  id,
		})
	}

	// 插入数据
	err := repositories.InsertObjects(r.DB, adminRoleSli)
	if err != nil {
		return fmt.Errorf("RoleReqContext.DispatchRoleReq() -> %v", err)
	}
	return nil
}
