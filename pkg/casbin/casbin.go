package casbin

import (
	"fmt"
	"forum/internal/models"
	"forum/pkg/globals"
	"forum/pkg/utils"
	"github.com/casbin/casbin/v2"
	"github.com/casbin/casbin/v2/model"
	gormadapter "github.com/casbin/gorm-adapter/v3"
	"gorm.io/gorm"
)

// casbin_r_m_a 结构体
type CasbinService struct {
	Enforcer *casbin.Enforcer
	//Adapter  *gormadapter.Adapter
}

// NewCasbinObject
// @Description: 生成casbinService结构体
// @Author wangyulong 2025-02-25 20:23:25
// @return       CasbinService
func NewCasbinObject() CasbinService {
	return CasbinService{
		Enforcer: globals.CasbinEnforcer,
	}
}

// NewCasbinService
// @Description: 创建 casbin_r_m_a 结构体
// @Author wangyulong 2025-01-17 09:31:09
// @param        db *gorm.DB
// @return       *CasbinService
// @return       error
func NewCasbinService(db *gorm.DB) (*CasbinService, error) {

	// 确保数据库连接有效
	if db == nil {
		return nil, fmt.Errorf("NewCasbinService -> 数据库连接无效")
	}

	// 创建适配器
	a, err := gormadapter.NewAdapterByDB(db)
	if err != nil {
		return nil, fmt.Errorf("NewCasbinService -> 创建适配器失败 -> %v", err)
	}

	// 检查适配器是否成功创建
	if a == nil {
		return nil, fmt.Errorf("NewCasbinService -> 创建适配器失败，适配器为空")
	}

	// 创建模型
	modelString := `
    [request_definition]
    r = sub, obj

    [policy_definition]
    p = sub, obj

    [role_definition]
    g = _,_

    [policy_effect]
    e = some(where (p.eft == allow))

    [matchers]
    m = g(r.sub, p.sub) && r.obj == p.obj
    `
	m, err1 := model.NewModelFromString(modelString)
	if err1 != nil {
		return nil, fmt.Errorf("NewCasbinService -> 创建模型失败 -> %s", err1)
	}

	// 确保模型被正确加载
	if m == nil {
		return nil, fmt.Errorf("NewCasbinService -> 创建模型失败，模型为空")
	}

	// 创建执行器
	e, err2 := casbin.NewEnforcer(m, a)
	if err2 != nil {
		return nil, fmt.Errorf("NewCasbinService -> 创建执行器失败 -> %s", err2)
	}

	return &CasbinService{
		Enforcer: e,
	}, nil
}

// GetApiPermForRole 获取当前角色的api权限
func (c *CasbinService) GetApiPermForRole(email string) ([]uint, error) {

	// 确保最新的策略数据
	err := c.Enforcer.LoadPolicy()
	if err != nil {
		return nil, fmt.Errorf("(c *CasbinService) GetApiPermForRole -> 策略加载失败， 已有的会忽略 -> %s", err)
	}

	var resources []uint
	permissions, err := c.Enforcer.GetPermissionsForUser(email)
	if err != nil {
		return nil, fmt.Errorf("GetApiPermForRole -> 获取当前角色的api权限失败 -> %s", err)
	}
	for _, p := range permissions {
		obj := p[1] // p[1] 是资源字段
		apiId, err := utils.ChangeStringToUint(obj)
		if err != nil {
			return nil, fmt.Errorf("GetApiPermForRole -> 获取当前角色的api权限失败 -> %s", err)
		}
		resources = append(resources, apiId)
	}
	return resources, nil
}

// ModifyRolePolicy 重置角色权限
func (c *CasbinService) ModifyRolePolicy(email string, apiIds []string) error {

	// 不直接操作数据库，利用enforcer简化操作
	err := c.Enforcer.LoadPolicy()
	if err != nil {
		return fmt.Errorf("(c *CasbinService) ModifyRolePolicy -> 策略加载失败 -> %s", err)
	}

	// 将角色的所有api权限全部删除
	//_, err = c.Enforcer.DeletePermissionsForUser(strconv.FormatUint(uint64(roleId), 10))
	_, err = c.Enforcer.DeletePermissionsForUser(email)
	if err != nil {
		return fmt.Errorf("(c *CasbinService) ModifyRolePolicy -> 删除角色组权限失败 -> %s", err)
	}

	// 为角色分配新的api权限
	for _, id := range apiIds {
		_, err = c.Enforcer.AddPolicy(email, id)
		if err != nil {
			return fmt.Errorf("(c *CasbinService) ModifyRolePolicy -> 修改角色组权限成功 -> %s", err)
		}
	}

	// 保存策略(保存到存储)
	return c.Enforcer.SavePolicy()

}

/*// SelApiId 根据请求路径查找api的id
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
}*/

// AssignRolesForAdminOrUser
// @Description: 为管理员或用户分配角色
// @Author wangyulong 2025-01-15 11:32:04
// @receiver     c
// @param        userId uint 用户ID
// @param        ids []uint 角色IDs
// @return       error
func (c *CasbinService) AssignRolesForAdminOrUser(email string, ids []uint) error {
	// 确保最新的策略数据
	err := c.Enforcer.LoadPolicy()
	if err != nil {
		return fmt.Errorf("(c *CasbinService) ModifyRolePolicy -> 策略加载失败， 已有的会忽略 -> %s", err)
	}

	// 遍历要分配的角色
	for _, id := range ids {
		role := fmt.Sprintf("%v", id)

		// 为用户添加单个角色
		ok, err := c.Enforcer.AddRoleForUser(email, role)
		if err != nil {
			return fmt.Errorf("(c *CasbinService) AssignRolesForAdminOrUser -> 为用户分配角色失败: %s", err)
		}
		if !ok {

		}
	}

	// 如果需要持久化到数据库
	if err := c.Enforcer.SavePolicy(); err != nil {
		return fmt.Errorf("(c *CasbinService) AssignRolesForAdminOrUser -> 保存策略失败: %s", err)
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
	//	return fmt.Errorf("(c *CasbinService) AssignRolesForAdminOrUser -> 为用户分配角色失败 %s", err)
	//}
	//if !ok {
	//	return fmt.Errorf("(c *CasbinService) AssignRolesForAdminOrUser -> 该用户已经拥有该角色")
	//}
	//// 保存策略(保存到存储)
	//return c.Enforcer.SavePolicy()

}

// GetRolesForAdminOrUser
// @Description: 获取管理员或用户拥有的角色ID
// @Author wangyulong 2025-01-16 16:14:23
// @receiver     c
// @param        email uint
// @return       []uint
// @return       error
func (c *CasbinService) GetRolesForAdminOrUser(email string) ([]uint, error) {

	/*// 用于存储要返回的用户拥有的角色id
	var ids []uint*/

	// 确保最新的策略数据
	err := c.Enforcer.LoadPolicy()
	if err != nil {
		return nil, fmt.Errorf("(c *CasbinService) GetRolesForAdminOrUser -> 策略加载失败， 已有的会忽略 -> %s", err)
	}
	// 获取用户拥有的角色
	roles, err := c.Enforcer.GetRolesForUser(email)
	if err != nil {
		return nil, fmt.Errorf("c *CasbinService) GetRolesForAdminOrUser -> 获取用户拥有的角色异常 -> %s", err)
	}

	// 根据角色Name查询角色ID
	roleIds, err := c.GetIdForRole(roles)
	if err != nil {
		return nil, fmt.Errorf("(c *CasbinService) GetRolesForAdminOrUser -> 策略加载失败， 已有的会忽略 -> %s", err)
	}

	/*var id uint
	for _, roleId := range roles {
		id, err = utils.ChangeStringToUint(roleId)
		if err != nil {
			return nil, fmt.Errorf("c *CasbinService) GetRolesForAdminOrUser -> %s", err)
		}
		ids = append(ids, id)
	}*/

	return roleIds, nil
}

// DeleteRoleForAdminOrUser
// @Description: 删除管理员或用户对应的角色ids
// @Author wangyulong 2025-01-17 09:30:30
// @receiver     c
// @param        userId uint
// @param        roleIds []uint
// @return       error
func (c *CasbinService) DeleteRoleForAdminOrUser(email string, roleIds []uint) error {

	// 根据角色ID查询其对应的Name
	roles, err := c.GetRoleForId(roleIds)
	if err != nil {
		return fmt.Errorf("(c *CasbinService) DeleteRoleForAdminOrUser -> 策略加载失败， 已有的会忽略 -> %s", err)
	}

	// 确保最新的策略数据
	err = c.Enforcer.LoadPolicy()
	if err != nil {
		return fmt.Errorf("(c *CasbinService) DeleteRoleForAdminOrUser -> 策略加载失败， 已有的会忽略 -> %s", err)
	}

	for _, role := range roles {
		ok, err := c.Enforcer.DeleteRoleForUser(email, role)
		if err != nil {
			return fmt.Errorf("c *CasbinService) DeleteRoleForAdminOrUser -> 删除用户对应的角色异常 -> %s", err)
		}
		if !ok {
			return fmt.Errorf("c *CasbinService) DeleteRoleForAdminOrUser -> 删除用户对应的角色失败 -> %s", err)
		}
	}

	// 如果需要持久化到数据库
	if err := c.Enforcer.SavePolicy(); err != nil {
		return fmt.Errorf("(c *CasbinService) DeleteRoleForAdminOrUser -> 保存策略失败: %s", err)
	}

	return nil
}

// UpdateRoleForAdminOrUser
// @Description: 更改管理员或用户的角色id
// @Author wangyulong 2025-01-17 09:29:37
// @receiver     c
// @param        email uint
// @param        roleIds []uint
// @return       error
func (c *CasbinService) UpdateRoleForAdminOrUser(email string, roleIds []uint) error {

	// 根据角色ID查询其对应的Name
	roles, err := c.GetRoleForId(roleIds)
	if err != nil {
		return fmt.Errorf("(c *CasbinService) DeleteRoleForAdminOrUser -> 策略加载失败， 已有的会忽略 -> %s", err)
	}

	// 确保最新的策略数据
	err = c.Enforcer.LoadPolicy()
	if err != nil {
		return fmt.Errorf("(c *CasbinService) UpdateRoleForAdminOrUser -> 策略加载失败， 已有的会忽略 -> %s", err)
	}

	// 删除用户所有对应的角色
	_, err = c.Enforcer.DeleteRolesForUser(email)
	if err != nil {
		return fmt.Errorf("c *CasbinService) UpdateRoleForAdminOrUser -> 删除用户所有对应的角色异常 -> %s", err)
	}

	// 遍历要分配的角色
	for _, role := range roles {

		// 为用户添加单个角色
		ok, err := c.Enforcer.AddRoleForUser(email, role)
		if err != nil {
			return fmt.Errorf("(c *CasbinService) UpdateRoleForAdminOrUser -> 为用户分配角色失败: %s", err)
		}
		if !ok {
			// 不用做任何处理
		}
	}

	// 如果需要持久化到数据库
	if err := c.Enforcer.SavePolicy(); err != nil {
		return fmt.Errorf("(c *CasbinService) UpdateRoleForAdminOrUser -> 保存策略失败: %s", err)
	}

	return nil
}

// VerifySuperAdministrator
// @Description: 验证管理员或用户是否是超级管理员
// @Author wangyulong 2025-01-24 18:28:02
// @receiver     c
// @param        email string
// @param        role string
// @return       bool
// @return       error
func (c *CasbinService) VerifySuperAdministrator(email string, role string) (bool, error) {
	// 确保最新的策略数据
	err := c.Enforcer.LoadPolicy()
	if err != nil {
		return false, fmt.Errorf("(c *CasbinService) VerifySuperAdministrator -> 策略加载失败， 已有的会忽略 -> %s", err)
	}

	// 验证用户是否是超级管理员
	ok, err := c.Enforcer.HasRoleForUser(email, role)
	if err != nil {
		return false, fmt.Errorf("(c *CasbinService) VerifySuperAdministrator -> 验证用户是否是超级管理员异常 -> %s", err)
	}
	if ok {
		return true, nil
	} else {
		return false, nil
	}
}

// DeletePermForAdminOrUser
// @Description: 删除角色拥有的权限
// @Author wangyulong 2025-01-24 20:26:47
// @receiver     c
// @param        role string
// @return       error
func (c *CasbinService) DeletePermForAdminOrUser(role string) error {
	// 删除角色拥有的权限
	_, err := c.Enforcer.DeletePermissionForUser(role)
	if err != nil {
		return fmt.Errorf("(c *CasbinService) DeletePermForAdminOrUser -> 删除角色拥有的权限异常 -> %s", err)
	}
	return nil
}

// GetIdForRole
// @Description: 查询角色对应的ID
// @Author wangyulong 2025-02-23 15:31:11
// @receiver     *CasbinService
// @param        roles []string
func (*CasbinService) GetIdForRole(roles []string) ([]uint, error) {
	var resIds []uint
	err := globals.DB.Model(&models.Role{}).Where("name IN ?", roles).Pluck("id", &resIds).Error
	if err != nil {
		return nil, fmt.Errorf("(*CasbinService) GetIdForRole -> 查询角色对应的ID 异常 -> %v", err)
	}
	return resIds, nil
}

// GetRoleForId
// @Description: 查询角色对应的Name
// @Author wangyulong 2025-02-23 15:36:48
// @receiver     *CasbinService
// @param        ids []uint
func (*CasbinService) GetRoleForId(ids []uint) ([]string, error) {
	var resRoles []string
	err := globals.DB.Model(&models.Role{}).Where("id IN ?", ids).Pluck("name", &resRoles).Error
	if err != nil {
		return nil, fmt.Errorf("(*CasbinService) GetRoleForId -> 查询角色对应的Nam 异常 -> %v", err)
	}
	return resRoles, nil
}
