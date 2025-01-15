package casbin

import (
	"fmt"
	"forum/internal/internalPkg/internalUtils"
	"forum/internal/models"
	"forum/pkg/globals"
	"github.com/casbin/casbin/v2"
	"github.com/casbin/casbin/v2/model"
	gormadapter "github.com/casbin/gorm-adapter/v3"
	"gorm.io/gorm"
)

// casbin_r_m_a 结构体
type CasbinService struct {
	Enforcer *casbin.Enforcer
	Adapter  *gormadapter.Adapter
}

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

// GetApiPerm 获取当前角色的api权限
func (c *CasbinService) GetApiPerm(id string) ([]uint, error) {
	var resources []uint
	permissions, err := c.Enforcer.GetPermissionsForUser(id)
	if err != nil {
		return nil, fmt.Errorf("GetApiPerm -> 获取当前角色的api权限失败 -> %s", err)
	}
	for _, p := range permissions {
		obj := p[1] // p[1] 是资源字段
		apiId, err := internalUtils.ChangeStringToUint(obj)
		if err != nil {
			return nil, fmt.Errorf("GetApiPerm -> 获取当前角色的api权限失败 -> %s", err)
		}
		resources = append(resources, apiId)
	}
	return resources, nil
}

// ModifyRolePolicy 重置角色权限
func (c *CasbinService) ModifyRolePolicy(roleId string, apiIds []string) error {

	// 不直接操作数据库，利用enforcer简化操作
	err := c.Enforcer.LoadPolicy()
	if err != nil {
		return fmt.Errorf("ModifyRolePolicy -> 创建角色组权限， 已有的会忽略 -> %s", err)
	}

	// 将角色的所有api权限全部删除
	//_, err = c.Enforcer.DeletePermissionsForUser(strconv.FormatUint(uint64(roleId), 10))
	_, err = c.Enforcer.DeletePermissionsForUser(roleId)
	if err != nil {
		return fmt.Errorf("ModifyRolePolicy -> 删除角色组权限失败 -> %s", err)
	}

	// 为角色分配新的api权限
	for _, id := range apiIds {
		_, err = c.Enforcer.AddPolicy(roleId, id)
		if err != nil {
			return fmt.Errorf("ModifyRolePolicy -> 修改角色组权限成功 -> %s", err)
		}
	}

	// 保存策略(保存到存储)
	return c.Enforcer.SavePolicy()

}

// SelApiId 根据请求路径查找api的id
func SelApiId(requestUrl string) (uint, error) {
	var apiId uint
	// 根据请求中的路由接口,查询apiId
	err := globals.DB.Model(models.Api{}).Select("ID").Where("path = ?", requestUrl).Scan(&apiId).Error
	if err != nil {
		return 0, fmt.Errorf("SelApiId -> 根据请求路径，获取apiID失败，该用户没有该权限 -> %s", err)
	}
	if err == nil && apiId == 0 {
		return 0, fmt.Errorf("SelApiId -> 该api不存在，该用户没有该权限")
	}
	return apiId, nil
}

// AssignRolesForUser
// @Description: 为用户分配角色
// @Author wangyulong 2025-01-15 11:32:04
// @receiver     c
// @param        userId uint 用户ID
// @param        ids []uint 角色IDs
// @return       error
func (c *CasbinService) AssignRolesForUser(userId uint, ids []uint) error {
	// 确保最新的策略数据
	err := c.Enforcer.LoadPolicy()
	if err != nil {
		return fmt.Errorf("ModifyRolePolicy -> 创建角色组权限， 已有的会忽略 -> %s", err)
	}

	// 遍历要分配的角色
	for _, id := range ids {
		role := fmt.Sprintf("%v", id)

		// 为用户添加单个角色
		ok, err := c.Enforcer.AddRoleForUser(fmt.Sprintf("%v", userId), role)
		if err != nil {
			return fmt.Errorf("(c *CasbinService) AssignRolesForUser -> 为用户分配角色失败: %s", err)
		}
		if !ok {

		}
	}

	// 如果需要持久化到数据库
	if err := c.Enforcer.SavePolicy(); err != nil {
		return fmt.Errorf("(c *CasbinService) AssignRolesForUser -> 保存策略失败: %s", err)
	}

	return nil

	//var tempId string
	//var tempIDS []string
	//for _, id := range ids {
	//	tempId = fmt.Sprintf("%v", id)
	//	tempIDS = append(tempIDS, tempId)
	//}
	//
	//ok, err := c.Enforcer.AddRolesForUser(fmt.Sprintf("%v", userId), tempIDS)
	//if err != nil {
	//	return fmt.Errorf("(c *CasbinService) AssignRolesForUser -> 为用户分配角色失败 %s", err)
	//}
	//if !ok {
	//	return fmt.Errorf("(c *CasbinService) AssignRolesForUser -> 该用户已经拥有该角色")
	//}
	//// 保存策略(保存到存储)
	//return c.Enforcer.SavePolicy()

}
