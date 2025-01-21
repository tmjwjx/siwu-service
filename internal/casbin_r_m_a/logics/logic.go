package logics

import (
	"forum/internal/casbin_r_m_a/repositories"
	"forum/internal/casbin_r_m_a/requests"
	"forum/pkg/casbin"
	casbin2 "github.com/casbin/casbin/v2"
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
func GetApiPermLogic(e *casbin2.Enforcer, id string) (req *requests.GetApiPermRes, err error) {

	casbinService := &casbin.CasbinService{
		Enforcer: e,
	}

	res, err := repositories.GetApiPermRep(casbinService, id)
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
