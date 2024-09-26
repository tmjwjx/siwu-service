package logics

import (
	"forum/internal/tag/repositories"
	"forum/internal/tag/requests"
	"gorm.io/gorm"
)

// UpdateTagUserCountLogic 更新数据库中标签的关注人数
func UpdateTagUserCountLogic(db *gorm.DB, tagID uint) (*requests.TagFansCountRes, error) {
	// 更新数据库中标签的关注人数
	fansCount, err := repositories.UpdateTagUserCountReq(db, tagID)
	return fansCount, err
}

// UpdateTagArticleCountLogic 更新前端的标签页
func UpdateTagArticleCountLogic() (*requests.TagRes, error) {
	tagRes, err := repositories.UpdateTagArticleCountReq()
	return tagRes, err
}
