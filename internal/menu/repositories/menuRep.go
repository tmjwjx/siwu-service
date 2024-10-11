package repositories

import (
	"fmt"
	"forum/internal/menu/requests"
	"forum/internal/models"
	"gorm.io/gorm"
)

// MenuSearchRep 检索获取所有菜单列表
func MenuSearchRep(db *gorm.DB, req *requests.MenuSearchReq) (*requests.MenuSearchRes, error) {

	var menuSearchRes *requests.MenuSearchRes
	var menus []models.Menu
	var menusRes []*requests.Menus
	// 建立表关联
	query := db.Model(models.Menu{})

	// 添加查询条件
	if req.Icon != "" {
		query = query.Where("icon = ?", req.Icon)
	}
	if req.Name != "" {
		query = query.Where("name = ?", req.Name)
	}
	if req.Type == 1 || req.Type == 2 || req.Type == 3 {
		query = query.Where("type = ?", req.Type)
	}
	if req.RouteName != "" {
		query = query.Where("route_name = ?", req.RouteName)
	}
	if req.RoutePath != "" {
		query = query.Where("route_path = ?", req.RoutePath)
	}
	if req.Visible == 1 || req.Visible == 2 || req.Visible == 3 {
		query = query.Where("name = ?", req.Name)
	}
	if req.Code != "" {
		query = query.Where("code = ?", req.Code)
	}
	if req.ComponentPath != "" {
		query = query.Where("component_path = ?", req.ComponentPath)
	}
	if req.ParentId != nil {
		query = query.Where("pid = ?", req.ParentId)
	}

	// 查询数据
	err := query.Limit(req.Limit).Offset(req.Page).Find(&menus).Error
	if err != nil {
		return nil, fmt.Errorf("MenuSearchRep -> 检索获取所有菜单列表失败 -> %s", err)
	}

	total := len(menus)

	for _, menu := range menus {
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
		}

		menusRes = append(menusRes, m)
	}

	menuSearchRes = &requests.MenuSearchRes{
		Menus: menusRes,
		Total: total,
	}

	return menuSearchRes, nil

}

// GetMenuIconRep 获取所有菜单图标
func GetMenuIconRep(db *gorm.DB) (*requests.GetMenuIconRes, error) {

	var icons []string
	var iconList []*requests.IconRes
	err := db.Model(models.Menu{}).Select("Icon").Find(&icons).Error
	if err != nil {
		return nil, fmt.Errorf("GetMenuIconRep -> 获取所有菜单图标 -> %s", err)
	}

	for _, icon := range icons {
		i := &requests.IconRes{
			Icon: icon,
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

	err := db.First(&menu).Error
	if err == nil {
		return fmt.Errorf("该菜单已经存在")
	}

	// 开启事务
	tx := db.Begin()
	if tx.Error != nil {
		return fmt.Errorf("CreateMenuRep -> 开启事务失败 -> %s", tx.Error)
	}

	// 添加菜单
	err = tx.Create(&menu).Error
	if err != nil {
		tx.Rollback() // 回滚事务
		return fmt.Errorf("CreateMenuRep -> 添加菜单失败 -> %s", tx.Error)
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
		return fmt.Errorf("CreateMenuRep -> 开启事务失败 -> %s", tx.Error)
	}

	// 删除菜单
	err := tx.Model(models.Menu{}).Where("id IN ?", req.IDs).Delete(nil).Error
	if err != nil {
		tx.Rollback() // 回滚事务
		return fmt.Errorf("CreateMenuRep -> 删除菜单失败 -> %s", err)
	}

	//提交事务
	err = tx.Commit().Error
	if err != nil {
		return fmt.Errorf("CreateMenuRep -> 提交事务失败 -> %s", err)
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
		return fmt.Errorf("UpdateMenuRep -> 修改菜单失败 -> %s", err)
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

	return res, nil

}
