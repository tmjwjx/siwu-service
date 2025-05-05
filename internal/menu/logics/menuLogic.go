package logics

import (
	"forum/internal/menu/repositories"
	"forum/internal/menu/requests"
	"gorm.io/gorm"
)

// MenuSearchLogic 检索获取所有菜单列表
func MenuSearchLogic(db *gorm.DB, req *requests.MenuSearchReq) (*requests.MenuSearchRes, error) {
	menuSearchRes, err := repositories.MenuSearchRep(db, req)
	return menuSearchRes, err
}

// GetMenuIconLogic 获取所有菜单图标
func GetMenuIconLogic(db *gorm.DB) (*requests.GetMenuIconRes, error) {
	res, err := repositories.GetMenuIconRep(db)
	return res, err
}

// CreateMenuLogic 新建菜单
func CreateMenuLogic(db *gorm.DB, req *requests.CreateMenuReq) error {
	err := repositories.CreateMenuRep(db, req)
	return err
}

// DeleteMenuLogic 删除菜单
func DeleteMenuLogic(db *gorm.DB, req *requests.DeleteMenuReq) error {
	err := repositories.DeleteMenuRep(db, req)
	return err
}

// UpdateMenuLogic 修改菜单
func UpdateMenuLogic(db *gorm.DB, req *requests.UpdateMenuReq) error {
	err := repositories.UpdateMenuRep(db, req)
	return err
}

// GetMenuDetailLogic 获取当前菜单详情
func GetMenuDetailLogic(db *gorm.DB, id string) (*requests.GetMenuDetailRes, error) {
	res, err := repositories.GetMenuDetailRep(db, id)
	return res, err
}

// GetSpecificMenuLogic
// @Description: 查询所有type为1和2的菜单
// @Author wangyulong 2024-10-15 11:25:24
// @param        db *gorm.DB
// @return       *requests.GetSpecificMenuRes
// @return       error
func GetSpecificMenuLogic(db *gorm.DB) (*requests.GetSpecificMenuRes, error) {
	res, err := repositories.GetSpecificMenuRep(db)
	return res, err
}
