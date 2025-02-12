package repositories

import (
	"forum/internal/internalPkg/internalUtils"
	"forum/internal/message/requests"
	"forum/internal/models"
	"forum/pkg/globals"
	"gorm.io/gorm"
)

// CommentRead
// @Description: 评论消息已读
// @param        db *gorm.DB
// @param        id uint
// @return       err
// @Author tianjiajie 2025-02-12 20:56:54
func CommentRead(db *gorm.DB, id uint) (err error) {
	// 获取文章id
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

	// 查询评论id
	var b []int64
	err = db.Model(&models.ArticleComment{}).
		Select("id").
		Where("user_id = ?", id).
		Pluck("id", &b).
		Error
	if err != nil {
		globals.Log.Errorf("Failed to get article id: %v", err)
		return err
	}

	err = db.Model(&models.ArticleComment{}).
		Where("article_id in ? or parent_id in ?", a, b).
		Update("is_read", 1).
		Error

	if err != nil {
		globals.Log.Errorf("Failed to mark likes as read: %v", err)
		return err
	}

	return nil
}

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
		Where("deleted_at IS NULL").
		Debug().
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
		Joins("left join sw_articles on sw_article_comments.article_id = sw_articles.id").                                                   // 文章
		Joins("left join sw_users on sw_article_comments.user_id = sw_users.id").                                                            // 用户
		Joins("left join sw_comment_likes on sw_comment_likes.comment_id = sw_article_comments.id AND sw_comment_likes.deleted_at IS NULL"). // 点赞
		Joins("left join sw_article_comments sac on sac.id = sw_article_comments.parent_id").                                                // 上一级评论
		Joins("left join sw_attachments on sw_attachments.home_id = sw_article_comments.user_id AND sw_attachments.home = ?", "user").       // 查询用户头像
		Order("sw_article_comments.created_at DESC").
		Debug().
		Limit(req.Limit).
		Offset((req.Page - 1) * req.Limit)

	// 选择
	query = query.Select("DISTINCT sw_article_comments.user_id, "+
		"sw_users.nickname, "+
		"sw_articles.title, "+
		"sw_articles.id as article_id, "+
		"sw_article_comments.content, "+
		"sw_article_comments.created_at, "+
		"sw_article_comments.id AS comment_id, "+
		"sw_article_comments.likes_count, "+
		"sw_article_comments.parent_id, "+
		"sw_article_comments.created_at, "+
		"sac.content AS comment, "+
		"sw_attachments.path,"+
		"IF(sw_comment_likes.user_id = ?, 1, 0) AS status", userId)

	// 查询条件 文章下的评论
	query = query.Where("sw_articles.user_id = ? or sw_article_comments.parent_id IN (SELECT id FROM sw_article_comments WHERE user_id = ?)", userId, userId).
		Where("sw_article_comments.user_id != ?", userId)

	// 查询条件 评论下的评论
	//query = query.Where("sw_article_comments.parent_id IN (SELECT id FROM sw_article_comments WHERE user_id = ?)", userId)

	if err = query.
		Find(&res).Error; err != nil {
		return res, err
	}

	// 格式化时间
	for i := range res {
		res[i].FormatTime = internalUtils.TimeFormat(res[i].CreatedAt)
		res[i].DailyTime = internalUtils.TimeFormatDaily(res[i].CreatedAt)
	}

	// 评论消息已读
	err = CommentRead(db, userId)
	if err != nil {
		return res, err
	}

	return res, err
}
