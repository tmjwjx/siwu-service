package repositories

import (
	"fmt"
	"forum/internal/image/controllers"
	"forum/internal/models"
	"forum/internal/tag/requests"
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

// UpdateTagArticleCountReq 更新前端的标签页
func UpdateTagArticleCountReq() ([]*requests.TagRes, error) {

	// 查询数据时要用到的结构体
	var tags []models.Tag
	// 用于存储要响应给前端的数据
	var tagRes []*requests.TagRes
	// 查询所有标签及其文章
	err := globals.DB.Preload("Articles").Find(&tags).Error
	if err != nil {
		return nil, fmt.Errorf("UpdateTagArticleCountReq -> %s", err)
	}
	// 将数据库中的数据，写到要响应的结构体中
	for _, tag := range tags {
		t := &requests.TagRes{}
		t.ID = tag.ID
		t.Name = tag.Name
		t.Description = tag.Description
		tag.ArticleCount = len(tag.Articles)
		// 如果标签中文章的数量超过1000，就用 k 来表示，否者用原型
		if tag.ArticleCount > 1000 {
			t.ArticleCount = fmt.Sprintf("%dk", tag.ArticleCount/1000)
		} else {
			t.ArticleCount = fmt.Sprintf("%d", tag.ArticleCount)
		}
		for _, article := range tag.Articles {
			tag.Heat = article.LikesCount/2 + article.CollectionsCount + article.CommentsCount
		}
		// 如果标签的热度超过1000，就用 k 来表示，否者用原型
		if tag.Heat > 1000 {
			t.Heat = fmt.Sprintf("%dk", tag.Heat/1000)
		} else {
			t.Heat = fmt.Sprintf("%d", tag.Heat)
		}
		if tag.FansCount > 1000 {
			t.FansCount = fmt.Sprintf("%.1fk", tag.FansCount/1000)
		} else {
			t.FansCount = fmt.Sprintf("%d", tag.FansCount)
		}
		// 将图片存入结构体 t 中
		images, err := controllers.GetImagesControllers("标签", tag.ID)
		if err != nil {
			return nil, fmt.Errorf("UpdateTagArticleCountReq -> %s", err)
		}
		for _, image := range *images {
			t.Path = append(t.Path, image.Path)
		}
		tagRes = append(tagRes, t)
	}

	return tagRes, nil
}
