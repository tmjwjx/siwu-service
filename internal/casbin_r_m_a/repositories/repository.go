package repositories

import (
	"fmt"
	"forum/internal/casbin_r_m_a/requests"
	"forum/internal/models"
	"forum/pkg/casbin"
	"gorm.io/gorm"
)

// AssignMenuPermRep 为角色分配菜单权限
func AssignMenuPermRep(db *gorm.DB, req *requests.AssignMenuPermReq) error {

	for _, menuId := range req.PermIds {

		roleMenu := &models.RoleMenu{
			RoleId: req.ID,
			MenuId: menuId,
		}

		err := db.First(roleMenu).Error
		if err == nil {
			// 如果该角色已经拥有了该权限，就不用再给他分配该权限了，直接跳过，接着分配下一个权限。
			continue
		}

		err = db.Create(roleMenu).Error
		if err != nil {
			return fmt.Errorf("AssignMenuPermRep -> 为角色分配菜单权限失败 -> %s", err)
		}
	}

	return nil
}

// GetMenuPermRep 获取当前角色的菜单权限(用于渲染侧边栏，只要type1和2)
func GetMenuPermRep(db *gorm.DB, id string) (*requests.GetMenuPermRes, error) {

	var menuId uint
	// 查询菜单id
	err := db.Model(&models.RoleMenu{}).Select("MenuId").Where("role_id = ?", id).Scan(&menuId).Error
	if err != nil {
		return nil, fmt.Errorf("GetMenuPermRep -> 该角色没有任何权限 -> %s", err)
	}
	// 查询权限
	var menuPerm []requests.MenuPerm
	err = db.Model(models.Menu{}).Where("id = ? and type IN (?)", menuId, []uint{1, 2}).Scan(&menuPerm).Error
	if err != nil {
		return nil, fmt.Errorf("GetMenuPermRep -> 查询权限失败 -> %s", err)
	}

	res := &requests.GetMenuPermRes{
		Perm: menuPerm,
	}

	return res, nil
}

// GetApiPermRep 获取当前角色的api权限
func GetApiPermRep(db *gorm.DB, id string) (req *requests.GetApiPermRes, err error) {

	casbinService, err := casbin.NewCasbinService(db)
	if err != nil {
		return nil, fmt.Errorf("GetApiPermRep -> 获取当前角色的api权限失败 -> %s", err)
	}

	apiIds, err := casbinService.GetApiPerm(id)
	if err != nil {
		return nil, fmt.Errorf("GetApiPermRep -> 获取当前角色的api权限失败 -> %s", err)
	}

	apiGroup := make(map[uint]uint) // 用来存放apiId和其对应的groupId
	var groups []*requests.Group

	for _, id := range apiIds {
		flag := 0 // 用于判断新查询到的groupId是否，已在groups切片中有相同的值了

		var groupId uint
		// 查询分组id
		err = db.Model(models.ApiGroup{}).Select("GroupId").Where("api_id = ?", id).Scan(&groupId).Error
		if err != nil {
			return nil, fmt.Errorf("GetApiPermRep -> 获取当前角色的api权限失败 -> %s", err)
		}
		apiGroup[id] = groupId

		// 将groups中的group中的重复的groupId只保留一个
		for _, group := range groups {
			if groupId == group.GroupId {
				flag = 1
				break
			}
		}
		if flag == 0 {
			var children []*requests.Child
			group := &requests.Group{
				GroupId:  groupId,
				Children: children,
			}
			groups = append(groups, group)
		}

	}

	for _, group := range groups {
		for apiId, groupId := range apiGroup {
			if groupId == group.GroupId {
				child := &requests.Child{
					ID: apiId,
				}
				group.Children = append(group.Children, child)
				// 将已经添加的键值对删除，防止之后，再次重复遍历
			}
		}
	}

	res := &requests.GetApiPermRes{
		Group: groups,
	}

	return res, nil

}

// AssignApiPermRep 为角色分配api权限
func AssignApiPermRep(db *gorm.DB, req *requests.AssignApiPermReq) error {

	casbinService, err := casbin.NewCasbinService(db)
	if err != nil {
		return fmt.Errorf("AssignApiPermRep -> 为角色分配api权限失败 -> %s", err)
	}
	// 为角色分配api权限
	err = casbinService.ModifyRolePolicy(req.ID, req.Apis)
	if err != nil {
		return fmt.Errorf("AssignApiPermRep -> 为角色分配api权限失败 -> %s", err)
	}
	return nil
}

// GetPermCodeRep 获取当前角色的所有权限标识
func GetPermCodeRep(db *gorm.DB, id string) (*requests.GetPermCodeRes, error) {

	var codeList []string
	err := db.Model(&models.RoleMenu{}).Joins("join sw_menu on sw_menu.id = sw_role_menu.menu_id").Where("role_id = ?", id).Pluck("sw_menu.Code", &codeList).Error
	if err != nil {
		return nil, fmt.Errorf("GetPermCodeRep -> 获取当前角色的所有权限标识失败 -> %s", err)
	}

	res := &requests.GetPermCodeRes{
		CodeList: codeList,
	}

	return res, nil
}
