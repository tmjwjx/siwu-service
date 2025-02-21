package logics

import (
	"errors"
	"forum/internal/administrator/repositories"
	"forum/internal/administrator/requests"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func DeleteAdministratorLogic(db *gorm.DB, req requests.DeleteAdministratorReq) (err error) {
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
