package logics

import (
	"errors"
	"forum/internal/administrator/repositories"
	"forum/internal/administrator/requests"
	"forum/internal/internalPkg/internalUtils"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// UpdateAdministratorLogic
// @Description: 编辑管理员信息
// @param        db *gorm.DB
// @param        req requests.UpdateAdministratorReq
// @return       data
// @return       err
// @Author tianjiajie 2025-02-21 21:56:20
func UpdateAdministratorLogic(db *gorm.DB, req requests.UpdateAdministratorReq, id any) (err error) {
	//// 身份验证 todo
	//if id != req.ID {
	//	return errors.New("不是本人操作")
	//}

	// 编辑管理员信息
	err = repositories.UpdateAdministratorRep(db, req)
	if err != nil {
		return err
	}
	return nil
}

// ResetAdministratorPasswordLogic
// @Description: 重置管理员密码
// @param        db *gorm.DB
// @param        req requests.ResetAdministratorPasswordReq
// @return       err
// @Author tianjiajie 2025-02-21 21:31:18
func ResetAdministratorPasswordLogic(db *gorm.DB, req requests.AdministratorReq) (err error) {
	// 重置管理员密码
	err = repositories.ResetAdministratorPasswordRep(db, req)
	if err != nil {
		return err
	}
	return nil
}

// GetAdministratorInfoLogic
// @Description: 查询管理员详情
// @param        db *gorm.DB
// @param        id string
// @return       data
// @return       err
// @Author tianjiajie 2025-02-21 21:21:48
func GetAdministratorInfoLogic(db *gorm.DB, id string) (data interface{}, err error) {
	// 查询管理员详情
	administrator, err := repositories.GetAdministratorInfoRep(db, id)
	if err != nil {
		return nil, err
	}

	// 获取头像
	str, _ := internalUtils.GetImages(db, "administrator", administrator.ID)
	if len(*str) >= 1 {
		administrator.Avatar = (*str)[0]
	}

	data = administrator
	return data, nil
}

// GetAdministratorListLogic
// @Description: 查询管理员列表
// @param        db *gorm.DB
// @param        req requests.GetAdministratorListReq
// @return       data
// @return       err
// @Author tianjiajie 2025-02-21 20:33:41
func GetAdministratorListLogic(db *gorm.DB, req requests.GetAdministratorListReq) (data interface{}, err error) {
	// 查询管理员列表
	administratorList, err := repositories.GetAdministratorListRep(db, req)
	if err != nil {
		return nil, err
	}

	// 获取头像
	for i := range administratorList {
		str, _ := internalUtils.GetImages(db, "administrator", administratorList[i].ID)
		if len(*str) >= 1 {
			administratorList[i].Avatar = (*str)[0]
		}
	}

	data = gin.H{"manager_list": administratorList}
	return data, nil
}

// DeleteAdministratorLogic
// @Description: 删除管理员
// @param        db *gorm.DB
// @param        req requests.AdministratorReq
// @return       err
// @Author tianjiajie 2025-02-21 20:33:31
func DeleteAdministratorLogic(db *gorm.DB, req requests.AdministratorReq) (err error) {
	// 删除管理员
	err = repositories.DeleteAdministratorRep(db, req)
	if err != nil {
		return err
	}
	return nil
}

// AddAdministratorLogic
// @Description: 添加管理员
// @param        db *gorm.DB
// @param        req *requests.AddAdministratorReq
// @return       articleList
// @return       err
// @Author tianjiajie 2025-02-20 20:46:16
func AddAdministratorLogic(db *gorm.DB, req requests.AddAdministratorReq) (data interface{}, err error) {

	// 邮箱查重
	b := repositories.CheckAdministratorEmail(db, req)
	if b { // 如果存在，返回错误信息
		return nil, errors.New("邮箱已存在")
	}

	// 添加管理员
	administratorId, err := repositories.AddAdministratorRep(db, req)
	if err != nil {
		return nil, err
	}
	data = gin.H{"id": administratorId}
	return data, nil
}
