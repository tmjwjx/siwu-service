package repositories

import (
	"forum/internal/internalPkg/internalUtils"
	"forum/internal/message/requests"
	"gorm.io/gorm"
)

// IsCommentLikeRep
// @Description: 查询是否给评论点赞
// @param        db *gorm.DB
// @param        userId uint
// @param        commentId uint
// @return       bool
// @return       error
// @Author tianjiajie 2025-01-18 20:18:42
func IsCommentLikeRep(db *gorm.DB, userId uint, commentId uint) (bool, error) {
	var count int64
	err := db.Table("sw_comment_likes").
		Where("user_id = ? AND comment_id = ?", userId, commentId).
		Count(&count).Error
	if err != nil {
		return false, err
	}
	if count > 0 {
		return true, nil
	}
	return false, nil
}

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
		Joins("left join sw_comment_likes on sw_comment_likes.comment_id = sw_article_comments.id").
		Joins("left join sw_article_comments sac on sac.id = sw_article_comments.parent_id").
		Order("sw_article_comments.created_at DESC").
		Limit(req.Limit).
		Offset((req.Page - 1) * req.Limit)

	// 查询条件
	query = query.Where("sw_article_comments.user_id = ?", userId)

	// 查询用户头像
	query = query.
		Joins("left join sw_attachments on sw_attachments.home_id = sw_article_comments.user_id AND sw_attachments.home = ?", "user")

	// 选择
	query = query.Select("sw_article_comments.user_id, " +
		"sw_users.nickname, " +
		"sw_articles.title, " +
		"sw_articles.id as article_id, " +
		"sw_article_comments.content, " +
		"sw_article_comments.created_at, " +
		//"COUNT(sw_comment_likes.user_id) AS like_count, " + // 这里使用 COUNT()
		"sw_article_comments.id AS comment_id, " +
		"sw_article_comments.likes_count, " +
		//"sac.id AS parent_id, " +
		"sw_article_comments.parent_id, " +
		"sac.content AS comment, " +
		"sw_attachments.path").
		Group("sw_article_comments.id, sw_attachments.path, sw_article_comments.created_at") // 添加这一行用于去重

	if err = query.
		Find(&res).Error; err != nil {
		return res, err
	}

	// 格式化时间
	for i := range res {
		res[i].FormatTime = internalUtils.TimeFormat(res[i].CreatedAt)
		res[i].DailyTime = internalUtils.TimeFormatDaily(res[i].CreatedAt)
	}

	return res, err

}
