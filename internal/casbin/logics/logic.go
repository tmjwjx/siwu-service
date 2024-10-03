package logics

import (
	"forum/internal/casbin/repositories"
	"forum/internal/casbin/requests"
	"gorm.io/gorm"
)

// AssignMenuPermLogic 为角色分配菜单权限
func AssignMenuPermLogic(db *gorm.DB, req *requests.AssignMenuPermReq) error {
	err := repositories.AssignMenuPermRep(db, req)
	return err
}

// GetMenuPermLogic 获取当前角色的菜单权限(用于渲染侧边栏，只要type1和2)
func GetMenuPermLogic(db *gorm.DB, id string) (*requests.GetMenuPermRes, error) {
	res, err := repositories.GetMenuPermRep(db, id)
	return res, err
}

// GetApiPermLogic 获取当前角色的api权限
func GetApiPermLogic(db *gorm.DB, id string) (req *requests.GetApiPermRes, err error) {
	res, err := repositories.GetApiPermRep(db, id)
	return res, err
}

// AssignApiPermLogic 为角色分配api权限
func AssignApiPermLogic(db *gorm.DB, req *requests.AssignApiPermReq) error {
	err := repositories.AssignApiPermRep(db, req)
	return err
}

// GetPermCodeLogic 获取当前角色的所有权限标识
func GetPermCodeLogic(db *gorm.DB, id string) (*requests.GetPermCodeRes, error) {
	res, err := repositories.GetPermCodeRep(db, id)
	return res, err
}
