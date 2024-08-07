package repositories

import (
	"fmt"
	"forum/internal/models"
	"forum/pkg/globals"
)

// UpdateTagUserCountReq 更新数据库中标签的关注人数
func UpdateTagUserCountReq(tagID uint) (string, error) {
	var tag models.Tag
	var fansCount string // 统计现在的人数
	// 查询该标签是否存在
	if err := globals.DB.Take(&tag, "id = ?", tagID).Error; err != nil {
		return "", fmt.Errorf("UpdateTagUserCountReq -> %s", err)
	}
	// 查询该标签的关注人数
	if err := globals.DB.Model(&tag).Update("fans_count", tag.FansCount+1).Error; err != nil {
		return "", fmt.Errorf("UpdateTagUserCountReq -> %s", err)
	}
	if tag.FansCount > 1000 {
		fansCount = fmt.Sprintf("%.1fk", tag.FansCount/1000)
	} else {
		fansCount = fmt.Sprintf("%d", tag.FansCount)
	}
	return fansCount, nil
}

// UpdateTagArticleCountReq 更新前端页面中标签的文章数量
func UpdateTagArticleCountReq() (string, error) {

	var tags []models.Tag
	var articleCount string // 要传递的参数，文章数量
	// 统计文章总数量时，用到的变量，int64 是一个 64 位整数类型，
	//它能表示非常大的整数范围（-2^63 到 2^63-1）。选择 int64
	//类型的原因是为了确保在处理大量数据时不会发生整数溢出
	var totalArticles int64
	// 查询所有标签及其文章
	err := globals.DB.Preload("Articles").Find(&tags).Error
	if err != nil {
		return "", fmt.Errorf("UpdateTagArticleCountReq -> %s", err)
	}
	// 统计所有标签下的文章总数
	for _, tag := range tags {
		totalArticles += int64(len(tag.Articles))
	}
	// 如果标签中文章的数量超过1000，就用 k 来表示，否者用原型
	if totalArticles > 1000 {
		articleCount = fmt.Sprintf("%.1fk", totalArticles/1000)
	} else {
		articleCount = fmt.Sprintf("%d", totalArticles)
	}
	return articleCount, nil
}

// UpdateTagHeat 更新前端页面中标签的热度
func UpdateTagHeat() (string, error) {
	var tags []models.Tag
	var heat int64
	var totalHeat string
	err := globals.DB.Preload("Articles").Find(&tags).Error
	if err != nil {
		return "", fmt.Errorf("UpdateTagHeat -> %s", err)
	}
	// 统计所有标签下的文章总数
	for _, tag := range tags {
		//totalArticles += int64(len(tag.Articles))
		for _, article := range tag.Articles {
			heat += int64(article.LikesCount/2 + article.CollectionsCount + article.CommentsCount)
		}
	}
	// 如果标签的热度超过1000，就用 k 来表示，否者用原型
	if heat > 1000 {
		totalHeat = fmt.Sprintf("%.1fk", heat/1000)
	} else {
		totalHeat = fmt.Sprintf("%d", heat)
	}
	return totalHeat, nil
}
