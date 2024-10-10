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
func InitUserInfoLogic(db *gorm.DB, qid string, gid string) (*requests.InitUserInfoRes, error) {
	req, err := repositories.InitUserInfoRep(db, qid, gid)
	return req, err
}
