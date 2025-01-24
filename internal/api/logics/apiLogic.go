package logics

import (
	"forum/internal/api/repositories"
	"forum/internal/api/requests"
	casbin2 "github.com/casbin/casbin/v2"
	"gorm.io/gorm"
)

//func ApiInitLogic(db *gorm.DB) (*[]*requests.ApiInitRes, error) {
//	res, err := repositories.ApiInitRep(db)
//	return res, err
//}

// GetAllApiLogic
// @Description: 获取所有api列表
// @Author wangyulong 2024-10-15 08:55:08
// @param        db *gorm.DB
// @return       *requests.GetAllApiRes
// @return       error
func GetAllApiLogic(db *gorm.DB) (*requests.GetAllApiRes, error) {
	res, err := repositories.GetAllApiRep(db)
	return res, err
}

// GetApiDetailsLogic 获取当前api详情
func GetApiDetailsLogic(db *gorm.DB, id uint) (*requests.ApiDetailsRes, error) {
	apiDetailsRes, err := repositories.GetApiDetailsRep(db, id)
	return apiDetailsRes, err
}

// GetGroupListLogic 获取所有api分组列表
func GetGroupListLogic(db *gorm.DB) (*requests.ApiGroupRes, error) {
	apiGroupRes, err := repositories.GetGroupListRep(db)
	return apiGroupRes, err
}

// GetRequestMethodLogic 获取所有请求方法
func GetRequestMethodLogic(db *gorm.DB) (*requests.ApiReqMethodRes, error) {
	apiReqMethodRes, err := repositories.GetRequestMethodRep(db)
	return apiReqMethodRes, err
}

// DeleteApiLogic 删除api
func DeleteApiLogic(e *casbin2.Enforcer, req *requests.DeleteApiReq, db *gorm.DB) error {
	err := repositories.DeleteApiRep(e, req, db)
	return err
}

// UpdateApiLogic 编辑api
func UpdateApiLogic(db *gorm.DB, req *requests.UpdateApiReq) error {
	err := repositories.UpdateApiRep(db, req)
	return err
}

// CreateApiLogic 添加api
func CreateApiLogic(db *gorm.DB, req *requests.CreateApiReq) error {
	err := repositories.CreateApiRep(db, req)
	return err
}

// SearchApiListLogic 检索api列表
func SearchApiListLogic(db *gorm.DB, req *requests.SearchApiListReq) (*requests.SearchApiListRes, error) {
	res, err := repositories.SearchApiListRep(db, req)
	return res, err
}
