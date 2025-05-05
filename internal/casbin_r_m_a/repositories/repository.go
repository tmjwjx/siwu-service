package repositories

import (
	"fmt"
	"forum/internal/casbin_r_m_a/requests"
	"forum/internal/models"
	"forum/pkg/casbin"
	casbin2 "github.com/casbin/casbin/v2"
	"gorm.io/gorm"
)

// AssignMenuPermRep 为角色分配菜单权限
func AssignMenuPermRep(db *gorm.DB, req *requests.AssignMenuPermReq) error {

	// 查询表中的所有旧数据
	var rm []models.RoleMenu
	err := db.Model(models.RoleMenu{}).Where("role_id = ?", req.ID).Find(&rm).Error
	if err != nil {
		return fmt.Errorf("AssignMenuPermRep -> 查询 RoleMenu 表失败 -> %s", err)
	}

	// 删除表中的所有旧数据
	for _, rm2 := range rm {
		err = db.Model(&models.RoleMenu{}).Where("role_id = ? and menu_id = ?", rm2.RoleId, rm2.MenuId).Delete(nil).Error
		if err != nil {
			return fmt.Errorf("AssignMenuPermRep -> 删除 RoleMenu 表中旧数据失败 -> %s", err)
		}
	}

	for _, menuId := range req.PermIds {

		roleMenu := &models.RoleMenu{
			RoleId: req.ID,
			MenuId: menuId,
		}

		err := db.Model(models.RoleMenu{}).Where("role_id = ? and menu_id = ?", req.ID, menuId).First(roleMenu).Error
		if err == nil {
			// 如果该角色已经拥有了该权限，就不用再给他分配该权限了，直接跳过，接着分配下一个权限。
			continue
		}

		err = db.Model(&models.RoleMenu{}).Create(roleMenu).Error
		if err != nil {
			return fmt.Errorf("AssignMenuPermRep -> 为角色分配菜单权限失败 -> %s", err)
		}
	}

	return nil
}

// GetMenuPermRep 获取当前角色的菜单权限(用于渲染侧边栏，只要type1和2)
func GetMenuPermRep(db *gorm.DB, id string) (*requests.GetMenuPermRes, error) {

	var menuId []uint
	var superAdminId uint

	// 查询超级管理员的id
	err := db.Model(&models.Role{}).Select("id").Where("name = ?", "超级管理员").Scan(&superAdminId).Error
	if err != nil {
		return nil, fmt.Errorf("GetMenuPermRep -> 查询超级管理员的id异常 -> %s", err)
	}

	if superAdminId == 0 {
		return nil, fmt.Errorf("GetMenuPermRep -> 不存在超级管理员这个角色 -> %s", err)
	}

	if id == fmt.Sprintf("%v", superAdminId) {
		// 查询菜单id
		err = db.Model(&models.Menu{}).Pluck("id", &menuId).Error
		if err != nil {
			return nil, fmt.Errorf("GetMenuPermRep -> 1查询菜单id异常 -> %s", err)
		}
	} else {
		// 查询菜单id
		err = db.Model(&models.RoleMenu{}).Select("MenuId").Where("role_id = ?", id).Scan(&menuId).Error
		if err != nil {
			return nil, fmt.Errorf("GetMenuPermRep -> 2查询菜单id异常 -> %s", err)
		}
	}

	// 查询权限
	var menuPerm []requests.MenuPerm
	err = db.Model(models.Menu{}).Where("id IN ? and type IN ?", menuId, []uint{1, 2}).Scan(&menuPerm).Error
	if err != nil {
		return nil, fmt.Errorf("GetMenuPermRep -> 查询权限失败 -> %s", err)
	}

	res := &requests.GetMenuPermRes{
		Perm: menuPerm,
	}

	return res, nil
}

// GetApiPermRep 获取当前角色的api权限
func GetApiPermRep(casbinService *casbin.CasbinService, id string) (req *requests.GetApiPermRes, err error) {

	apiIds, err := casbinService.GetApiPermForRole(id)
	if err != nil {
		return nil, fmt.Errorf("GetApiPermRep -> 获取当前角色的api权限失败 -> %s", err)
	}

	//apiGroup := make(map[uint]uint) // 用来存放apiId和其对应的groupId
	//var groups []*requests.Group
	//
	//for _, id := range apiIds {
	//	flag := 0 // 用于判断新查询到的groupId是否，已在groups切片中有相同的值了
	//
	//	var groupId uint
	//	// 查询分组id
	//	err = db.Model(models.ApiDictItemGroup{}).Select("GroupId").Where("api_id = ?", id).Scan(&groupId).Error
	//	if err != nil {
	//		return nil, fmt.Errorf("GetApiPermRep -> 获取当前角色的api权限失败 -> %s", err)
	//	}
	//	apiGroup[id] = groupId
	//
	//	// 将groups中的group中的重复的groupId只保留一个
	//	for _, group := range groups {
	//		if groupId == group.GroupId {
	//			flag = 1
	//			break
	//		}
	//	}
	//	if flag == 0 {
	//		var children []*requests.Child
	//		group := &requests.Group{
	//			GroupId:  groupId,
	//			Children: children,
	//		}
	//		groups = append(groups, group)
	//	}
	//
	//}
	//
	//for _, group := range groups {
	//	for apiId, groupId := range apiGroup {
	//		if groupId == group.GroupId {
	//			child := &requests.Child{
	//				ID: apiId,
	//			}
	//			group.Children = append(group.Children, child)
	//			// 将已经添加的键值对删除，防止之后，再次重复遍历
	//		}
	//	}
	//}
	//
	//res := &requests.GetApiPermRes{
	//	Group: groups,
	//}

	res := &requests.GetApiPermRes{
		Group: &apiIds,
	}
	return res, nil

}

// AssignApiPermRep 为角色分配api权限
func AssignApiPermRep(db *gorm.DB, req *requests.AssignApiPermReq, e *casbin2.Enforcer) error {

	casbinService := casbin.CasbinService{
		Enforcer: e,
	}

	var apiIds []string

	for _, apiId := range req.Apis {
		id := fmt.Sprintf("%v", apiId)
		apiIds = append(apiIds, id)
	}
	// 为角色分配api权限
	err := casbinService.ModifyRolePolicy(fmt.Sprintf("%v", req.ID), apiIds)
	if err != nil {
		return fmt.Errorf("AssignApiPermRep -> 为角色分配api权限失败 -> %s", err)
	}
	return nil
}

// GetPermCodeRep 获取当前角色的所有权限标识
func GetPermCodeRep(casbinService *casbin.CasbinService, db *gorm.DB, id string) (*requests.GetPermCodeRes, error) {

	var apiIds []uint

	var superAdminId uint

	// 查询超级管理员的id
	err := db.Model(&models.Role{}).Select("id").Where("name = ?", "超级管理员").Scan(&superAdminId).Error
	if err != nil {
		return nil, fmt.Errorf("GetPermCodeRep -> 查询超级管理员的id异常 -> %s", err)
	}

	if superAdminId == 0 {
		return nil, fmt.Errorf("GetPermCodeRep -> 不存在超级管理员这个角色 -> %s", err)
	}

	if id == fmt.Sprintf("%v", superAdminId) {
		// 获取当前角色的api权限
		err := db.Model(&models.Api{}).Pluck("id", &apiIds).Error
		if err != nil {
			return nil, fmt.Errorf("GetPermCodeRep -> 获取当前角色的api权限失败 -> %s", err)
		}
	} else {
		// 获取当前角色的api权限
		apiId, err := casbinService.GetApiPermForRole(id)
		if err != nil {
			return nil, fmt.Errorf("GetPermCodeRep -> 获取当前角色的api权限失败 -> %s", err)
		}
		apiIds = apiId
	}

	// 获取type为3的菜单(即按钮)
	var menuIds []uint
	err = db.Model(&models.Menu{}).Where("type = ?", 3).Pluck("id", &menuIds).Error
	if err != nil {
		return nil, fmt.Errorf("GetPermCodeRep -> 获取type为3的菜单(即按钮)失败 -> %s", err)
	}

	// 通过菜单id查询其拥有的api的id
	//menuApi := make(map[uint]uint)
	var endMenuId []uint
	var apiID []uint
	for _, menuId := range menuIds {
		err = db.Model(&models.MenuApi{}).Where("menu_id = ?", menuId).Pluck("api_id", &apiID).Error
		if err != nil {
			return nil, fmt.Errorf("GetPermCodeRep -> 通过菜单id查询其拥有的api的id失败 -> %s", err)
		}

		// 通过查询到的菜单拥有的api,去角色拥有的api权限中寻找相等的api，如果相等，这该菜单(即按钮)属于该角色
		for _, mApi := range apiID {
			flag := 0
			for _, rApi := range apiIds {
				if rApi == mApi {
					endMenuId = append(endMenuId, menuId)
					flag = 1
					break
				}
			}
			if flag == 1 {
				break
			}
		}
	}

	//通过查到的角色拥有的菜单(即按钮),去查询菜单的code字段
	var codes []string
	err = db.Model(&models.Menu{}).Where("id IN ?", endMenuId).Pluck("code", &codes).Error
	if err != nil {
		return nil, fmt.Errorf("GetPermCodeRep -> 通过菜单id查询其拥有的api的id失败 -> %s", err)
	}

	//var codeList []string
	//err = db.Model(&models.RoleMenu{}).
	//	Joins("join sw_menus on sw_menus.id = sw_role_menus.menu_id").
	//	Where("sw_role_menus.role_id = ? and sw_menus.type = ?", id, 3).
	//	Distinct("sw_menus.code").
	//	Pluck("sw_menus.code", &codeList).Error
	//if err != nil {
	//	return nil, fmt.Errorf("GetPermCodeRep -> 获取当前角色的所有权限标识失败 -> %s", err)
	//}

	res := &requests.GetPermCodeRes{
		CodeList: codes,
	}

	return res, nil
}
