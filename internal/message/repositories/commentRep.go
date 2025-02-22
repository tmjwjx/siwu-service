package repositories

import (
	"forum/internal/internalPkg/internalUtils"
	"forum/internal/message/requests"
	"forum/internal/models"
	"forum/pkg/globals"
	"gorm.io/gorm"
)

// CommentUnreadCount
// @Description: 评论消息未读数量
// @param        db *gorm.DB
// @param        id uint
// @return       count
// @return       err
// @Author tianjiajie 2025-02-12 21:24:39
func CommentUnreadCount(db *gorm.DB, id uint) (count int64, err error) {
	var a []int64
	err = db.Model(&models.Article{}).
		Select("id").
		Where("user_id = ?", id).
		Pluck("id", &a).
		Error
	if err != nil {
		globals.Log.Errorf("Failed to get article id: %v", err)
		return 0, err
	}

	err = db.Model(&models.ArticleComment{}).
		Where("article_id in ?", a).
		Where("is_read = ?", 0).
		Count(&count).
		Error

	if err != nil {
		globals.Log.Errorf("Failed to count unread likes: %v", err)
		return 0, err
	}

	return count, nil
}

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
		// 添加必要连接
		Joins("LEFT JOIN sw_articles ON sw_articles.id = sw_article_comments.article_id").
		Joins("LEFT JOIN sw_users ON sw_users.id = sw_article_comments.user_id").
		Joins("LEFT JOIN sw_attachments ON sw_attachments.home_id = sw_users.id AND sw_attachments.home = 'user'").
		Joins("LEFT JOIN sw_article_comments parent_comment ON parent_comment.id = sw_article_comments.parent_id").
		Order("sw_article_comments.created_at DESC").
		Limit(req.Limit).
		Offset((req.Page - 1) * req.Limit)

	query = query.Select(`
        DISTINCT sw_article_comments.id,
        sw_article_comments.user_id,
        sw_users.nickname,
        sw_articles.title,
        sw_articles.id AS article_id,
        sw_article_comments.content,
        sw_article_comments.likes_count,
        sw_article_comments.parent_id,
        sw_article_comments.created_at,
        parent_comment.content AS parent_comment,
        sw_attachments.path,
        EXISTS(
            SELECT 1 FROM sw_comment_likes
            WHERE comment_id = sw_article_comments.id
            AND user_id = ?
            AND deleted_at IS NULL
        ) AS status
    `, userId)

	query = query.Where("((sw_articles.user_id = ? AND sw_article_comments.user_id != ?)OR(parent_comment.user_id = ? AND sw_article_comments.user_id != ?))", userId, userId, userId, userId)

	if err = query.Find(&res).Error; err != nil {
		return res, err
	}

	// 时间格式化
	for i := range res {
		res[i].FormatTime = internalUtils.TimeFormat(res[i].CreatedAt)
		res[i].DailyTime = internalUtils.TimeFormatDaily(res[i].CreatedAt)
	}

	// 标记已读
	if err = CommentRead(db, userId); err != nil {
		return res, err
	}

	return res, nil
}
