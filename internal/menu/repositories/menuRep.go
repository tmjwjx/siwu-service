package repositories

import (
	"fmt"
	"forum/internal/menu/requests"
	"forum/internal/models"
	"gorm.io/gorm"
	"strconv"
)

// MenuSearchRep 检索获取所有菜单列表
func MenuSearchRep(db *gorm.DB, req *requests.MenuSearchReq) (*requests.MenuSearchRes, error) {

	//if req.Limit <= 0 && req.Page < 0 {
	//	return nil, fmt.Errorf("MenuSearchRep -> Limit的值不能小于等于0, 或者Page的值不能小于0")
	//}

	var menus []models.Menu
	var menusRes []*requests.Menus
	// 建立表关联
	query := db.Model(models.Menu{})

	// 添加查询条件
	if req.Icon != "" {
		query = query.Where("icon LIKE ?", "%"+req.Icon+"%")
	}
	if req.Name != "" {
		query = query.Where("name LIKE ?", "%"+req.Name+"%")
	}
	if req.Type == 1 || req.Type == 2 || req.Type == 3 {
		query = query.Where("type LIKE ?", "%"+strconv.Itoa(req.Type)+"%")
	}
	if req.RouteName != "" {
		query = query.Where("route_name LIKE ?", "%"+req.RouteName+"%")
	}
	if req.RoutePath != "" {
		query = query.Where("route_path LIKE ?", "%"+req.RoutePath+"%")
	}
	if req.Visible == 1 || req.Visible == 2 {
		query = query.Where("isVisible LIKE ?", "%"+req.Name+"%")
	}
	if req.Code != "" {
		query = query.Where("code LIKE ?", "%"+req.Code+"%")
	}
	if req.ComponentPath != "" {
		query = query.Where("component_path LIKE ?", "%"+req.ComponentPath+"%")
	}
	if req.ParentId > 0 {
		query = query.Where("parent_id LIKE ?", "%"+fmt.Sprintf("%v", req.ParentId)+"%")
	}

	// 查询数据
	err := query.Find(&menus).Error
	if err != nil {
		return nil, fmt.Errorf("MenuSearchRep -> 检索获取所有菜单列表失败 -> %s", err)
	}

	length := len(menus)
	if length == 0 {
		//return nil, fmt.Errorf("没有符合条件的菜单")
		res := &requests.MenuSearchRes{
			Menus: make([]*requests.Menus, 0),
			Total: 0,
		}
		return res, nil
	}

	var apiIds []uint
	for _, menu := range menus {

		err = db.Model(&models.MenuApi{}).Where("menu_id = ?", menu.ID).Pluck("api_id", &apiIds).Error
		if err != nil {
			return nil, fmt.Errorf("MenuSearchRep -> 查询 menuApi 表失败 -> %s", err)
		}

		m := &requests.Menus{
			ID:            menu.ID,
			ParentId:      menu.ParentId,
			Name:          menu.Name,
			Code:          menu.Code,
			Icon:          menu.Icon,
			Type:          menu.Type,
			RouteName:     menu.RouteName,
			RoutePath:     menu.RoutePath,
			ComponentPath: menu.ComponentPath,
			Visible:       menu.Visible,
			Sort:          menu.Sort,
			Desc:          menu.Desc,
			ApiId:         apiIds,
		}

		menusRes = append(menusRes, m)
	}

	// 分页返回数据
	//offset := (req.Page - 1) * req.Limit
	//if length > offset {
	//	end := offset + req.Limit
	//	if end > length {
	//		end = length
	//	}
	//	menusRes = menusRes[offset:end]
	//	res := &requests.MenuSearchRes{
	//		Menus: menusRes,
	//		Total: length,
	//	}
	//	return res, nil
	//}
	//
	//res := &requests.MenuSearchRes{
	//	Menus: make([]*requests.Menus, 0),
	//	Total: 0,
	//}

	res := &requests.MenuSearchRes{
		Menus: menusRes,
		Total: length,
	}

	return res, nil
}

// GetMenuIconRep 获取所有菜单图标
func GetMenuIconRep(db *gorm.DB) (*requests.GetMenuIconRes, error) {

	var icons []string
	var iconList []*requests.IconRes

	// 查询所有菜单图标
	query := db.Table("sw_dict_types").Joins("join sw_dict_items on sw_dict_items.dict_type_code = sw_dict_types.code")

	err := query.Where("sw_dict_types.name = ?  AND sw_dict_items.deleted_at IS NULL", "图标名").Pluck("sw_dict_items.label", &icons).Error
	if err != nil {
		return nil, fmt.Errorf("GetRequestMethodRep -> 查询所有菜单图标失败 -> %s", err)
	}

	/*err := db.Model(models.Menu{}).
		Where("icon <> ?", "").
		Distinct("icon").
		Pluck("Icon", &icons).Error
	if err != nil {
		return nil, fmt.Errorf("GetMenuIconRep -> 获取所有菜单图标 -> %s", err)
	}*/

	for _, icon := range icons {
		i := &requests.IconRes{
			Label: icon,
			Value: icon,
		}
		iconList = append(iconList, i)
	}

	res := &requests.GetMenuIconRes{
		IconList: iconList,
	}

	return res, nil

}

// CreateMenuRep 新建菜单
func CreateMenuRep(db *gorm.DB, req *requests.CreateMenuReq) error {

	//var routePath string
	//if req.RouteParam != "" {
	//	routePath = routePath + req.RouteParam
	//}

	//err := db.First(&menu).Error
	//if err == nil {
	//	fmt.Println("----------------------------", menu.ID)
	//	return fmt.Errorf("该菜单已经存在")
	//}

	// 开启事务
	tx := db.Begin()
	if tx.Error != nil {
		return fmt.Errorf("CreateMenuRep -> 开启事务失败 -> %s", tx.Error)
	}

	menu := &models.Menu{
		ParentId:      req.ParentId,
		Name:          req.Name,
		Code:          req.Code,
		Icon:          req.Icon,
		Type:          req.Type,
		RouteName:     req.RouteName,
		RoutePath:     req.RoutePath,
		RouteParam:    req.RouteParam,
		ComponentPath: req.ComponentPath,
		Visible:       req.Visible,
		Sort:          req.Sort,
		Desc:          req.Desc,
	}

	// 添加菜单
	err := tx.Create(&menu).Error
	if err != nil {
		tx.Rollback() // 回滚事务
		return fmt.Errorf("CreateMenuRep -> 添加菜单失败 -> %s", tx.Error)
	}

	if req.Type == 3 {
		// 添加按钮拥有的api
		for _, apiId := range req.ApiIds {
			menuApi := &models.MenuApi{
				MenuID: menu.ID,
				ApiID:  apiId,
			}
			err = tx.Create(menuApi).Error
			if err != nil {
				tx.Rollback() // 回滚事务
				return fmt.Errorf("CreateMenuRep -> 添加按钮拥有的api失败 -> %s", tx.Error)
			}
		}
	}

	//提交事务
	err = tx.Commit().Error
	if err != nil {
		return fmt.Errorf("CreateMenuRep -> 提交事务失败 -> %s", err)
	}

	return nil

}

// DeleteMenuRep 删除菜单
func DeleteMenuRep(db *gorm.DB, req *requests.DeleteMenuReq) error {

	// 开启事务
	tx := db.Begin()
	if tx.Error != nil {
		return fmt.Errorf("DeleteMenuRep -> 开启事务失败 -> %s", tx.Error)
	}

	// 删除菜单
	err := tx.Where("id IN ?", req.IDs).Delete(&models.Menu{}).Error
	if err != nil {
		tx.Rollback() // 回滚事务
		return fmt.Errorf("DeleteMenuRep -> 删除菜单失败 -> %s", err)
	}

	//提交事务
	err = tx.Commit().Error
	if err != nil {
		return fmt.Errorf("DeleteMenuRep -> 提交事务失败 -> %s", err)
	}

	return nil

}

// UpdateMenuRep 修改菜单
func UpdateMenuRep(db *gorm.DB, req *requests.UpdateMenuReq) error {

	var menu models.Menu

	// 查询该菜单是否存在
	err := db.Where("id = ?", req.ID).First(&menu).Error
	if err != nil {
		return fmt.Errorf("UpdateMenuRep -> 该菜单不存在")
	}

	//if req.RouteParam != "" {
	//	req.RoutePath = req.RoutePath + req.RoutePath
	//}

	menu.ParentId = req.ParentId
	menu.Type = req.Type
	menu.Icon = req.Icon
	menu.Name = req.Name
	menu.Sort = req.Sort
	menu.Visible = req.Visible
	menu.RouteName = req.RouteName
	menu.RoutePath = req.RoutePath
	menu.ComponentPath = req.ComponentPath
	menu.Desc = req.Desc
	menu.Code = req.Code
	menu.RouteParam = req.RouteParam

	// 开启事务
	tx := db.Begin()
	if tx.Error != nil {
		return fmt.Errorf("UpdateMenuRep -> 开启事务失败 -> %s", tx.Error)
	}

	// 修改菜单
	err = tx.Updates(menu).Error
	if err != nil {
		tx.Rollback() // 回滚菜单
		return fmt.Errorf("UpdateMenuRep -> 修改菜单失败 -> %s", err)
	}

	if req.Type == 3 {
		// 删除原来按钮拥有的api
		err = tx.Model(&models.MenuApi{}).Where("menu_id = ?", req.ID).Delete(nil).Error
		if err != nil {
			tx.Rollback() // 回滚菜单
			return fmt.Errorf("UpdateMenuRep -> 删除原来按钮拥有的api失败 -> %s", err)
		}

		// 添加按钮新拥有的api
		for _, apiId := range req.ApiIds {
			menuApi := &models.MenuApi{
				MenuID: menu.ID,
				ApiID:  apiId,
			}
			err = tx.Create(menuApi).Error
			if err != nil {
				tx.Rollback() // 回滚事务
				return fmt.Errorf("CreateMenuRep -> 添加按钮新拥有的api失败 -> %s", tx.Error)
			}
		}
	}

	//提交事务
	err = tx.Commit().Error
	if err != nil {
		return fmt.Errorf("UpdateMenuRep -> 提交事务失败 -> %s", err)
	}

	return nil

}

// GetMenuDetailRep 获取当前菜单详情
func GetMenuDetailRep(db *gorm.DB, id string) (*requests.GetMenuDetailRes, error) {

	var menu models.Menu

	// 查询该菜单是否存在
	err := db.Where("id = ?", id).First(&menu).Error
	if err != nil {
		return nil, fmt.Errorf("GetMenuDetailRep -> 该菜单不存在")
	}

	res := &requests.GetMenuDetailRes{
		ID:            menu.ID,
		ParentId:      menu.ParentId,
		Name:          menu.Name,
		Code:          menu.Code,
		Icon:          menu.Icon,
		Type:          menu.Type,
		RouteName:     menu.RouteName,
		RoutePath:     menu.RoutePath,
		ComponentPath: menu.ComponentPath,
		Visible:       menu.Visible,
		RouteParam:    menu.RouteParam,
		Sort:          menu.Sort,
		Desc:          menu.Desc,
	}

	if res.Type == 3 {
		// 查询按钮拥有的api
		var apiIds []uint
		err = db.Model(&models.MenuApi{}).Where("menu_id = ?", id).Pluck("api_id", &apiIds).Error
		if err != nil {
			return nil, fmt.Errorf("GetMenuDetailRep -> 查询按钮拥有的api失败 -> %s", err)
		}
		res.ApiIds = apiIds
	}

	return res, nil

}

// GetSpecificMenuRep
// @Description: 获取所有type为1和2的菜单
// @Author wangyulong 2024-10-15 11:24:12
// @param        db *gorm.DB
// @return       *requests.GetSpecificMenuRes
// @return       error
func GetSpecificMenuRep(db *gorm.DB) (*requests.GetSpecificMenuRes, error) {

	var menuTypeList []requests.List
	// 查询所有type为1和2的菜单
	err := db.Model(models.Menu{}).Select("id, name, parent_id").Where("type IN ?", []int{1, 2}).Find(&menuTypeList).Error
	if err != nil {
		return nil, fmt.Errorf("GetSpecificMenuRep -> 查询所有type为1和2的菜单失败 -> %s", err)
	}
	res := &requests.GetSpecificMenuRes{
		MenuTypeList: &menuTypeList,
	}
	return res, nil
}
