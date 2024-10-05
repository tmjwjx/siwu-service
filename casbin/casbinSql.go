package casbin

import (
	"fmt"
	"forum/internal/internal_pkg/internal_utils"
	"forum/internal/models"
	"forum/pkg/globals"
	"strconv"
)

// GetApiPerm 获取当前角色的api权限
func (c *CasbinService) GetApiPerm(id string) ([]uint, error) {
	var resources []uint
	permissions, err := c.Enforcer.GetPermissionsForUser(id)
	if err != nil {
		return nil, fmt.Errorf("GetApiPerm -> 获取当前角色的api权限失败 -> %s", err)
	}
	for _, p := range permissions {
		obj := p[1] // p[1] 是资源字段
		apiId, err := internal_utils.ChangeType(obj)
		if err != nil {
			return nil, fmt.Errorf("GetApiPerm -> 获取当前角色的api权限失败 -> %s", err)
		}
		resources = append(resources, apiId)
	}
	return resources, nil
}

// ModifyRolePolicy 分配(修改)角色权限(策略)
func (c *CasbinService) ModifyRolePolicy(roleId uint, apiIds []uint) error {

	// 不直接操作数据库，利用enforcer简化操作
	err := c.Enforcer.LoadPolicy()
	if err != nil {
		return fmt.Errorf("ModifyRolePolicy -> 创建角色组权限， 已有的会忽略 -> %s", err)
	}

	_, err = c.Enforcer.DeletePermissionsForUser(strconv.FormatUint(uint64(roleId), 10))
	if err != nil {
		return fmt.Errorf("ModifyRolePolicy -> 删除角色组权限失败 -> %s", err)
	}

	for _, id := range apiIds {
		_, err = c.Enforcer.AddPolicy(roleId, id)
		if err != nil {
			return fmt.Errorf("ModifyRolePolicy -> 修改角色组权限成功 -> %s", err)
		}
	}

	return c.Enforcer.SavePolicy()

}

func SelApiId(requestUrl string) (uint, error) {
	var apiId uint
	// 根据请求中的路由接口,查询apiId
	err := globals.DB.Model(models.Api{}).Select("ID").Where("path = ?", requestUrl).Scan(&apiId).Error
	if err != nil {
		return 0, fmt.Errorf("SelApiId -> 根据请求路径，获取apiID失败 -> %s", err)
	}
	return apiId, nil
}

//// GetRoles 获取所有角色组
//func (c *CasbinService) GetRoles() ([]string, error) {
//	return c.Enforcer.GetAllRoles()
//}
//
//// GetRolePolicy 获取所有角色组权限
//func (c *CasbinService) GetRolePolicy() (roles []RolePolicy, err error) {
//	err = c.Adapter.GetDb().Model(&gormadapter.CasbinRule{}).Where("ptype = 'p'").Find(&roles).Error
//	if err != nil {
//		return nil, fmt.Errorf("CreateRolePolicy -> 获取所有角色组权限失败 -> %s", err)
//	}
//	return
//}

//// CreateRolePolicy 创建角色组权限， 已有的会忽略
//func (c *CasbinService) CreateRolePolicy(roleId, apiId uint) error {
//	// 不直接操作数据库，利用enforcer简化操作
//	err := c.Enforcer.LoadPolicy()
//	if err != nil {
//		return fmt.Errorf("CreateRolePolicy -> 创建角色组权限， 已有的会忽略 -> %s", err)
//	}
//	// 添加策略
//	_, err = c.Enforcer.AddPolicy(roleId, apiId)
//	if err != nil {
//		return fmt.Errorf("CreateRolePolicy -> 添加策略失败 -> %s", err)
//	}
//	// 保存策略
//	return c.Enforcer.SavePolicy()
//}
//
//// DeleteRolePolicy 删除角色组权限
//func (c *CasbinService) DeleteRolePolicy(r RolePolicy) error {
//	// 不直接操作数据库，利用enforcer简化操作
//	err := c.Enforcer.LoadPolicy()
//	if err != nil {
//		return fmt.Errorf("DeleteRolePolicy -> 创建角色组权限， 已有的会忽略 -> %s", err)
//	}
//	_, err = c.Enforcer.RemovePolicy(r.RoleId, r.SingleId, r.Kind)
//	if err != nil {
//		return fmt.Errorf("DeleteRolePolicy -> 删除角色组权限失败 -> %s", err)
//	}
//	return c.Enforcer.SavePolicy()
//}

//// CanAccess 验证用户权限
//func (c *CasbinService) CanAccess(username, menuId string) (ok bool, err error) {
//	return c.Enforcer.Enforce(username, menuId)
//}

//// SelMenId 根据请求路径，获取菜单ID
//func SelMenId(requestUrl, requestMethod string) (string, error) {
//	var menuId string
//	// 获取最后一个/后面的信息
//	re := regexp.MustCompile("/([^/]+)$")
//	matches := re.FindStringSubmatch(requestUrl)
//	// 判断是不是角色
//	err := SelMRole(matches[1])
//	if err != nil {
//		requestUrl = strings.Replace(requestUrl, matches[0], "", 1)
//	} else {
//		return "", err
//	}
//	// 根据请求中的路由接口,查询menuId
//	err = globals.DB.Model(models.Menu{}).Select("id").Where("request_url = ? and request_method", requestUrl, requestMethod).Scan(&menuId).Error
//	return menuId, fmt.Errorf("SelMenId -> 根据请求路径，获取菜单ID -> %s", err)
//}
//
//// SelMRole 模糊查询角色字段
//func SelMRole(role string) error {
//
//	var count int64
//
//	err := globals.DB.Table("casbin_rule").Where("v1 LIKE ?", role+"%").Count(&count).Error
//	if err != nil {
//		return fmt.Errorf("SelMRole -> 模糊查询角色字段失败 -> %s", err)
//	}
//
//	// 如果count>0,表示找到了记录
//	if count > 0 {
//		return nil
//	}
//
//	// 如果没有找到记录，返回false
//	return fmt.Errorf("SelMRole -> 该角色不存在")
//
//}