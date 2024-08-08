package logics

import (
	"forum/internal/tag/repositories"
	"forum/internal/tag/requests"
)

// UpdateTagUserCountLogic 更新数据库中标签的关注人数
func UpdateTagUserCountLogic(tagID uint) (string, error) {
	// 更新数据库中标签的关注人数
	fansCount, err := repositories.UpdateTagUserCountReq(tagID)
	return fansCount, err
}

// UpdateTagArticleCountLogic 更新前端的标签页
func UpdateTagArticleCountLogic() ([]*requests.TagRes, error) {
	tagRes, err := repositories.UpdateTagArticleCountReq()
	return tagRes, err
}
