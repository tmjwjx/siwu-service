package repositories

import (
	"fmt"
	"forum/internal/image/controllers"
	"forum/internal/internalPkg/internalUtils"
	"forum/internal/models"
	"forum/internal/tag/requests"
	"forum/pkg/globals"
	"gorm.io/gorm"
)

// UpdateTagUserCountReq 更新数据库中标签的关注人数
func UpdateTagUserCountReq(db *gorm.DB, tagID uint) (*requests.TagFansCountRes, error) {
	var tag models.Tag
	var fansCount string // 统计现在的人数

	// 开启事务
	tx := db.Begin()
	if tx.Error != nil {
		return nil, fmt.Errorf("UpdateTagUserCountReq -> 开启事务失败 -> %s", tx.Error)
	}

	// 查询该标签是否存在
	if err := tx.Take(&tag, "id = ?", tagID).Error; err != nil {
		tx.Rollback() // 回滚事务
		return nil, fmt.Errorf("UpdateTagUserCountReq -> 查询该标签是否存在失败 -> %s", err)
	}
	// 更新该标签的关注人数
	if err := tx.Model(&tag).Update("fans_count", tag.FansCount+1).Error; err != nil {
		tx.Rollback() // 回滚事务
		return nil, fmt.Errorf("UpdateTagUserCountReq -> 更新该标签的关注人数失败 -> %s", err)
	}
	if tag.FansCount > 1000 {
		fansCount = fmt.Sprintf("标签人数: %.1fk", float64(tag.FansCount/1000))
	} else {
		fansCount = fmt.Sprintf("标签人数: %d", tag.FansCount)
	}

	tagFansCount := &requests.TagFansCount{
		FansCount: fansCount,
	}

	tagFansCountRes := &requests.TagFansCountRes{
		TagFansCount: tagFansCount,
	}

	// 提交事务
	err := tx.Commit().Error
	if err != nil {
		return nil, fmt.Errorf("UpdateTagUserCountReq -> 提交事务失败 -> %s", err)
	}

	return tagFansCountRes, nil

}

// UpdateTagArticleCountReq 更新前端的标签页
func UpdateTagArticleCountReq() (*requests.TagRes, error) {

	// 查询数据时要用到的结构体
	var tags []models.Tag
	// 用于存储要响应给前端的数据
	var tagList []*requests.Tag
	// 查询所有标签及其文章
	err := globals.DB.Preload("Articles").Find(&tags).Error
	if err != nil {
		return nil, fmt.Errorf("UpdateTagArticleCountReq -> %s", err)
	}
	// 将数据库中的数据，写到要响应的结构体中
	for _, tag := range tags {
		t := &requests.Tag{}
		t.ID = tag.ID
		t.Name = tag.Name
		t.Description = tag.Description
		tag.ArticleCount = len(tag.Articles)
		// 如果标签中文章的数量超过1000，就用 k 来表示，否者用原型
		if tag.ArticleCount > 1000 {
			t.ArticleCount = fmt.Sprintf("%.1fk", float64(tag.ArticleCount/1000))
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
			t.FansCount = fmt.Sprintf("%.1fk", float64(tag.FansCount/1000))
		} else {
			t.FansCount = fmt.Sprintf("%d", tag.FansCount)
		}
		// 将图片存入结构体 t 中
		images, err := controllers.GetImagesControllers("标签", tag.ID)
		if err != nil {
			// return nil, fmt.Errorf("UpdateTagArticleCountReq -> %s", err)
			// 如果没有找到就使用默认标签头像图片
			t.Path = internalUtils.TagDefaultImage
		} else {
			for _, image := range *images {
				t.Path = image.Path
			}
		}

		tagList = append(tagList, t)

	}

	tagRes := &requests.TagRes{
		TagList: tagList,
	}

	return tagRes, nil

}