package repositories

import (
	"forum/internal/internalPkg/internalUtils"
	"forum/internal/message/requests"
	"forum/internal/models"
	"forum/pkg/globals"
	"gorm.io/gorm"
)

// CommentLikeRead
// @Description: 评论点赞消息已读
// @param        db *gorm.DB
// @param        id uint
// @return       err
// @Author tianjiajie 2025-02-12 21:13:37
func CommentLikeRead(db *gorm.DB, id uint) (err error) {
	var a []int64
	err = db.Model(&models.ArticleComment{}).
		Select("id").
		Where("user_id = ?", id).
		Pluck("id", &a).
		Error
	if err != nil {
		globals.Log.Errorf("Failed to get article id: %v", err)
		return err
	}

	err = db.Model(&models.CommentLike{}).
		Where("comment_id in ?", a).
		Update("is_read", 1).
		Error

	if err != nil {
		globals.Log.Errorf("Failed to mark likes as read: %v", err)
		return err
	}

	return nil
}

// CommentLikeRep
// @Description: 评论点赞消息
// @param        db *gorm.DB
// @param        req requests.MessageReq
// @param        userId uint
// @return       res
// @return       err
// @Author tianjiajie 2025-02-06 09:12:15
func CommentLikeRep(db *gorm.DB, req requests.MessageReq, userId uint) (res []requests.CommentLikeMessageRes, err error) {

	// 查询评论点赞消息
	query := db.Model(&models.ArticleComment{}).
		Where("sw_article_comments.user_id = ? AND sw_comment_likes.deleted_at IS NULL", userId).
		Joins("left join sw_comment_likes on sw_comment_likes.comment_id = sw_article_comments.id").
		Joins("left join sw_users on sw_users.id = sw_comment_likes.user_id").
		Joins("left join sw_articles on sw_articles.id = sw_article_comments.article_id").
		Limit(req.Limit).
		Offset((req.Page - 1) * req.Limit).
		Order("sw_comment_likes.created_at DESC")

	// 查询用户头像
	query = query.Joins("left join sw_attachments on sw_attachments.home_id = sw_users.id AND sw_attachments.home = ?", "user")

	// 查询字段
	query = query.Select("sw_users.id AS user_id, " +
		"sw_users.nickname, " +
		"sw_attachments.path, " +
		"sw_articles.id AS article_id, " +
		"sw_articles.title, " +
		"sw_article_comments.id AS comment_id, " +
		"sw_article_comments.content, " +
		"sw_comment_likes.created_at")

	if err = query.Find(&res).Error; err != nil {
		return nil, err
	}

	// 格式化时间
	for i := range res {
		res[i].FormatTime = internalUtils.TimeFormat(res[i].CreatedAt)
		res[i].DailyTime = internalUtils.TimeFormatDaily(res[i].CreatedAt)
	}

	// 评论点赞消息已读
	err = CommentLikeRead(db, userId)
	if err != nil {
		return nil, err
	}

	return res, nil
}
