package repositories

import (
	"forum/internal/message/requests"
	"forum/internal/models"
	"gorm.io/gorm"
)

// LikeRep
// @Description: 点赞消息
// @param        db *gorm.DB
// @param        req requests.MessageReq
// @param        id uint
// @return       rep
// @return       err
func LikeRep(db *gorm.DB, req requests.MessageReq, id uint) (res []requests.LikeMessageRes, err error) {

	query := db.Model(&models.Article{}).
		Select("sw_users.id as user_id, sw_users.nickname, sw_articles.id as article_id, sw_articles.title, sw_article_likes.created_at, sw_attachments.path").
		Joins("left join sw_article_likes on sw_article_likes.article_id = sw_articles.id").
		Joins("left join sw_users on sw_users.id = sw_article_likes.user_id").
		Where("sw_articles.user_id = ?", id).
		Limit(req.Limit).
		Offset((req.Page - 1) * req.Limit).
		Order("sw_article_likes.created_at DESC")

	// 查询用户头像
	query = query.Joins("left join sw_attachments on sw_attachments.home_id = sw_users.id AND sw_attachments.home = ?", "user")
	//Where("sw_attachments.home = ?", "user")

	if err = query.Find(&res).Error; err != nil {
		return nil, err
	}

	return res, nil
}