package repositories

import (
	"errors"
	"fmt"
	"forum/internal/internalPkg/internalUtils"
	"forum/internal/models"
	"forum/internal/tag/requests"
	"forum/pkg/globals"
	"gorm.io/gorm"
)

// UpdateTagUserCountReq 更新数据库中标签的关注人数
func UpdateTagUserCountReq(userId uint, db *gorm.DB, tagID uint) (*requests.TagFansCountRes, error) {
	var tag models.Tag
	var fansCount string // 统计现在的人数

	// 开启事务
	tx := db.Begin()
	if tx.Error != nil {
		return nil, fmt.Errorf("UpdateTagUserCountReq -> 开启事务失败 -> %s", tx.Error)
	}

	var userTag models.UserTag
	err := tx.Model(&models.UserTag{}).First(&userTag).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// 如果没有找到，说明该用户还没有对该标签点过赞，进行下一步点赞就可以了
		} else {
			return nil, fmt.Errorf("查询用户是否点赞过该标签异常 -> %s", err)
		}
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
	err = tx.Commit().Error
	if err != nil {
		return nil, fmt.Errorf("UpdateTagUserCountReq -> 提交事务失败 -> %s", err)
	}

	return tagFansCountRes, nil

}

// UpdateTagArticleCountReq 更新前端的标签页
func UpdateTagArticleCountReq(db *gorm.DB) (*requests.TagRes, error) {

	// 查询数据时要用到的结构体
	var tags []models.Tag
	// 用于存储要响应给前端的数据
	var tagList []*requests.Tag
	// 查询所有标签及其文章
	err := db.Preload("Articles").Find(&tags).Error
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
		images, err := internalUtils.GetImages(db, globals.TagHome, tag.ID)
		if err != nil {
			return nil, fmt.Errorf("UpdateTagArticleCountReq -> %s", err)
		} else {
			for _, path := range *images {
				t.Path = path
			}
		}

		var userTag models.UserTag
		err = db.Model(&models.UserTag{}).First(&userTag).Error
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				t.Status = 2
			} else {
				return nil, fmt.Errorf("UpdateTagArticleCountReq -> 查询用户是否关注该标签异常 -> %s", err)
			}
		} else {
			t.Status = 1
		}

		tagList = append(tagList, t)

	}

	tagRes := &requests.TagRes{
		TagList: tagList,
	}

	return tagRes, nil

}

// StorageTagRep
// @Description: 存储新用户选择的标签
// @Author wangyulong 2024-10-14 21:32:41
// @param        db *gorm.DB
// @param        req *requests.StorageTagReq
// @param        userId uint
// @return       error
func StorageTagRep(db *gorm.DB, req *requests.StorageTagReq, userId uint) error {

	// 开启事务
	tx := db.Begin()
	if tx.Error != nil {
		return fmt.Errorf("StorageTagRep -> 开启事务失败 -> %s", tx.Error)
	}

	for _, tagID := range req.TagIDs {
		userTag := &models.UserTag{
			UserID: userId,
			TagID:  tagID,
		}
		err := tx.Create(userTag).Error
		if err != nil {
			tx.Rollback() // 回滚事务
			return fmt.Errorf("StorageTagRep -> 存储新用户选择的标签失败 -> %s", err)
		}
	}

	// 提交事务
	err := tx.Commit().Error
	if err != nil {
		return fmt.Errorf("StorageTagRep -> 提交事务失败 -> %s", err)
	}

	return nil
}

// GetAllTagRep
// @Description: 获取所有标签的id和name
// @Author wangyulong 2024-10-14 21:56:00
// @param        db gorm.DB
// @return       *requests.GetAllTagRes
// @return       error
func GetAllTagRep(db *gorm.DB) (*requests.GetAllTagRes, error) {

	var tags []requests.T
	err := db.Model(models.Tag{}).Select("id, name").Scan(&tags).Error
	if err != nil {
		return nil, fmt.Errorf("GetAllTagRep -> 获取所有标签的id和name失败 -> %s", err)
	}
	res := &requests.GetAllTagRes{
		Tags: tags,
	}
	return res, nil
}
