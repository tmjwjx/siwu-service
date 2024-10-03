package casbin

import (
	"fmt"
	"github.com/casbin/casbin/v2"
	"github.com/casbin/casbin/v2/model"
	gormadapter "github.com/casbin/gorm-adapter/v3"
	"gorm.io/gorm"
)

func NewCasbinService(db *gorm.DB) (*CasbinService, error) {
	// 创建适配器
	a, err := gormadapter.NewAdapterByDB(db)
	if err != nil {
		return nil, fmt.Errorf("NewCasbinService -> 创建适配器失败 -> %s", err)
	}
	// 创建模型
	m, err := model.NewModelFromString(`
[request_definition]
r = sub, obj


[policy_definition]
p = sub, obj

[role_definition]
g = _,_

[policy_effect]
e = some(where (p.eft == allow))

[matchers]
m = r.sub == p.sub && r.obj == p.obj`)
	if err != nil {
		return nil, fmt.Errorf("NewCasbinService -> 创建模型失败 -> %s", err)
	}
	// 创建执行器
	e, err := casbin.NewEnforcer(m, a)
	if err != nil {
		return nil, fmt.Errorf("NewCasbinService -> 创建执行器失败 -> %s", err)
	}

	return &CasbinService{
		Enforcer: e,
		Adapter:  a,
	}, nil
}

//// GetRoles 获取所有角色组
//func (c *models.CasbinService) GetRoles() ([]string, error) {
//	return c.Enforcer.GetAllRoles()
//}
//
//// GetRolePolicy 获取所有角色组权限
//func (c *models.CasbinService) GetRolePolicy() (roles []models.RolePolicy, err error) {
//	err = c.Adapter.GetDb().Model(&gormadapter.CasbinRule{}).Where("ptype = 'p'").Find(&roles).Error
//	if err != nil {
//		return nil, fmt.Errorf("CreateRolePolicy -> 获取所有角色组权限失败 -> %s", err)
//	}
//	return
//}
//
//// CreateRolePolicy 创建角色组权限， 已有的会忽略
//func (c *models.CasbinService) CreateRolePolicy(r models.RolePolicy) error {
//	// 不直接操作数据库，利用enforcer简化操作
//	err := c.Enforcer.LoadPolicy()
//	if err != nil {
//		return fmt.Errorf("CreateRolePolicy -> 创建角色组权限， 已有的会忽略 -> %s", err)
//	}
//	// 添加策略
//	_, err = c.Enforcer.AddPolicy(r.RoleId, r.SingleId, r.Kind)
//	if err != nil {
//		return fmt.Errorf("CreateRolePolicy -> 添加策略失败 -> %s", err)
//	}
//	// 保存策略
//	return c.Enforcer.SavePolicy()
//}
//
//// DeleteRolePolicy 删除角色组权限
//func (c *models.CasbinService) DeleteRolePolicy(r models.RolePolicy) error {
//	// 不直接操作数据库，利用enforcer简化操作
//	err := c.Enforcer.LoadPolicy()
//	if err != nil {
//		return fmt.Errorf("DeleteRolePolicy -> 创建角色组权限， 已有的会忽略 -> %s", err)
//	}
//	_, err = c.Enforcer.RemovePolicy(r.RoleName, r.MenuId)
//	if err != nil {
//		return fmt.Errorf("DeleteRolePolicy -> 删除角色组权限失败 -> %s", err)
//	}
//	return c.Enforcer.SavePolicy()
//}
//
//// CanAccess 验证用户权限
//func (c *models.CasbinService) CanAccess(username, menuId string) (ok bool, err error) {
//	return c.Enforcer.Enforce(username, menuId)
//}
//
//// ModifyRolePolicy 修改角色策略
//func (c *models.CasbinService) ModifyRolePolicy(RoleName string, OldMenuId string, NewMenuId string) error {
//
//	// 不直接操作数据库，利用enforcer简化操作
//	err := c.Enforcer.LoadPolicy()
//	if err != nil {
//		return fmt.Errorf("ModifyRolePolicy -> 创建角色组权限， 已有的会忽略 -> %s", err)
//	}
//
//	_, err = c.Enforcer.RemovePolicy(RoleName, OldMenuId)
//	if err != nil {
//		return fmt.Errorf("ModifyRolePolicy -> 删除角色组权限失败 -> %s", err)
//	}
//
//	_, err = c.Enforcer.AddPolicy(RoleName, NewMenuId)
//	if err != nil {
//		return fmt.Errorf("ModifyRolePolicy -> 修改角色组权限成功 -> %s", err)
//	}
//
//	return c.Enforcer.SavePolicy()
//
//}
