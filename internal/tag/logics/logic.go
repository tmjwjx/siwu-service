package logics

import "forum/internal/tag/repositories"

// UpdateTagUserCountLogic 更新数据库中标签的关注人数
func UpdateTagUserCountLogic(tagID uint) (string, error) {
	// 更新数据库中标签的关注人数
	fansCount, err := repositories.UpdateTagUserCountReq(tagID)
	return fansCount, err
}

// UpdateTagArticleCountLogic 更新前端页面中标签的文章数量
func UpdateTagArticleCountLogic() (string, error) {
	// 更新前端页面中标签的文章数量
	articleCount, err := repositories.UpdateTagArticleCountReq()
	return articleCount, err
}

// UpdateTagHeatLogic 更新前端页面中标签的热度
func UpdateTagHeatLogic() (string, error) {
	// 更新前端页面中标签的热度
	totalHeat, err := repositories.UpdateTagHeat()
	return totalHeat, err
}
