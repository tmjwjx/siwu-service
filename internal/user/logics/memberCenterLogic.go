package logics

import (
	"forum/internal/user/repositories"
	"forum/internal/user/requests"
	"gorm.io/gorm"
)

// InitUserInfoLogic
// @Description: 初始化用户信息
// @Author wangyulong 2024-10-10 12:40:58
// @param        db *gorm.DB
// @param        qid string
// @param        gid string
// @return       *requests.InitUserInfoRes
// @return       error
func InitUserInfoLogic(db *gorm.DB, qid string, gid string, tag string) (*requests.InitUserInfoRes, error) {
	req, err := repositories.InitUserInfoRep(db, qid, gid, tag)
	return req, err
}

// EditSignatureLogic
// @Description: 编辑个签
// @Author wangyulong 2024-10-11 16:10:42
// @param        db *gorm.DB
// @param        req *requests.EditSignatureReq
// @param        id uint
// @return       error
func EditSignatureLogic(db *gorm.DB, req *requests.EditSignatureReq, id uint) error {
	err := repositories.EditSignatureRep(db, req, id)
	return err
}
