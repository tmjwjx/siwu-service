package repositories

import (
	"forum/internal/internalPkg/internalUtils"
	"forum/internal/message/requests"
	"forum/internal/models"
	"forum/pkg/globals"
	"gorm.io/gorm"
)

// CollectionRead
// @Description: 收藏消息已读
// @param        db *gorm.DB
// @param        id uint
// @return       err
// @Author tianjiajie 2025-02-12 20:56:29
func CollectionRead(db *gorm.DB, id uint) (err error) {
	var a []int64
	err = db.Model(&models.Article{}).
		Select("id").
		Where("user_id = ?", id).
		Pluck("id", &a).
		Error
	if err != nil {
		globals.Log.Errorf("Failed to get article id: %v", err)
		return err
	}

	err = db.Model(&models.ArticleCollection{}).
		Where("article_id in ?", a).
		Update("is_read", 1).
		Error

	if err != nil {
		globals.Log.Errorf("Failed to mark likes as read: %v", err)
		return err
	}

	return nil
}

// CollectionRep
// @Description: 收藏消息
// @param        db *gorm.DB
// @param        req requests.MessageReq
// @param        id uint
// @return       res
// @return       err
// @Author tianjiajie 2024-10-05 17:49:26
func CollectionRep(db *gorm.DB, req requests.MessageReq, id uint) (res []requests.LikeAndCollectionMessageRes, err error) {

	query := db.Model(&models.Article{}).
		Select("DISTINCT sw_users.id as user_id, sw_users.nickname, sw_articles.id as article_id, sw_articles.title, sw_article_collections.created_at, sw_attachments.path").
		// 这里是 right join 不是 left join
		Joins("right join sw_article_collections on sw_article_collections.article_id = sw_articles.id").
		Joins("left join sw_users on sw_users.id = sw_article_collections.user_id").
		Where("sw_articles.user_id = ? AND sw_article_collections.deleted_at IS NULL", id).
		Limit(req.Limit).
		Offset((req.Page - 1) * req.Limit).
		Order("sw_article_collections.created_at DESC")

	// 查询用户头像
	query = query.Joins("left join sw_attachments on sw_attachments.home_id = sw_users.id AND sw_attachments.home = ?", "user")
	//Where("sw_attachments.home = ?", "user")

	if err = query.Find(&res).Error; err != nil {
		return nil, err
	}

	// 格式化时间
	for i := range res {
		res[i].FormatTime = internalUtils.TimeFormat(res[i].CreatedAt)
		res[i].DailyTime = internalUtils.TimeFormatDaily(res[i].CreatedAt)
	}

	// 收藏消息已读
	err = CollectionRead(db, id)
	if err != nil {
		return nil, err
	}

	return res, nil
}
