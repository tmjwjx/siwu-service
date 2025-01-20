package logics

import (
	"forum/internal/tag/repositories"
	"forum/internal/tag/requests"
	"gorm.io/gorm"
)

// UpdateTagUserCountLogic 更新数据库中标签的关注人数
func UpdateTagUserCountLogic(userId uint, db *gorm.DB, tagID uint) (*requests.TagFansCountRes, error) {
	// 更新数据库中标签的关注人数
	fansCount, err := repositories.UpdateTagUserCountReq(userId, db, tagID)
	return fansCount, err
}

// UpdateTagArticleCountLogic 更新前端的标签页
func UpdateTagArticleCountLogic(db *gorm.DB) (*requests.TagRes, error) {
	tagRes, err := repositories.UpdateTagArticleCountReq(db)
	return tagRes, err
}

// StorageTagLogic
// @Description: 存储新用户选择的标签
// @Author wangyulong 2024-10-14 21:36:29
// @param        db *gorm.DB
// @param        req *requests.StorageTagReq
// @param        userId uint
// @return       error
func StorageTagLogic(db *gorm.DB, req *requests.StorageTagReq, userId uint) error {
	err := repositories.StorageTagRep(db, req, userId)
	return err
}

// GetAllTagLogic
// @Description: 获取所有标签的id和name
// @Author wangyulong 2024-10-14 21:59:18
// @param        db *gorm.DB
// @return       *requests.GetAllTagRes
// @return       error
func GetAllTagLogic(db *gorm.DB) (*requests.GetAllTagRes, error) {
	res, err := repositories.GetAllTagRep(db)
	return res, err
}
