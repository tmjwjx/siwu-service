package repositories

import (
	"forum/internal/message/requests"
	"gorm.io/gorm"
)

// CommentRep
// @Description: 评论消息
// @param        db *gorm.DB
// @param        req *requests.MessageReq
// @param        userId uint
// @return       res
// @return       err
// @Author tianjiajie 2025-01-18 11:25:24
func CommentRep(db *gorm.DB, req *requests.MessageReq, userId uint) (res []requests.CommentMessageRes, err error) {
	query := db.Table("sw_article_comments").
		Joins("left join sw_users on sw_article_comments.user_id = sw_users.id").
		Joins("left join sw_articles on sw_article_comments.article_id = sw_articles.id").
		Joins("left join sw_comment_likes on sw_comment_likes.article_id = sw_articles.id").
		Limit(req.Limit).
		Offset((req.Page - 1) * req.Limit).
		Order("sw_article_collections.created_at DESC")

	// 查询用户头像
	query = query.
		Joins("left join sw_attachments on sw_attachments.home_id = sw_article_comments.user_id AND sw_attachments.home = ?", "user")

	// 选择
	query = query.Select("sw_article_comments.user_id, " +
		"sw_users.nickname, " +
		"sw_articles.title, " +
		"sw_article_comments.content, " +
		"sw_article_comments.created_at, " +
		"sw_article_comments.article_id" +
		"sw_article_comments.id" +
		"sw_attachments.path")

	if err = query.
		Debug().
		Find(&res).Error; err != nil {
		return res, err
	}
	return res, err

}
