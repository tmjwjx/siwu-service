package logics

import (
	"forum/internal/message/repositories"
	"forum/internal/message/requests"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// CommentLikeUnreadCountLogic
// @Description: 未读评论点赞消息数量
// @param        db *gorm.DB
// @param        userId uint
// @return       data
// @return       err
// @Author tianjiajie 2025-02-12 22:14:01
func CommentLikeUnreadCountLogic(db *gorm.DB, userId uint) (data interface{}, err error) {
	count, err := repositories.CommentLikeUnreadCount(db, userId)
	if err != nil {
		return 0, err
	}
	data = gin.H{"count": count}
	return data, nil
}

// CommentLikeMessageLogic
// @Description: 评论点赞消息
// @param        db *gorm.DB
// @param        req requests.MessageReq
// @param        userId uint
// @return       data
// @return       err
// @Author tianjiajie 2025-02-12 22:13:52
func CommentLikeMessageLogic(db *gorm.DB, req requests.MessageReq, userId uint) (data interface{}, err error) {
	// 查询评论点赞消息
	likeList, err := repositories.CommentLikeRep(db, req, userId)
	if err != nil {
		return nil, err
	}

	data = gin.H{"like_list": likeList}

	// 返回查询结果
	return data, nil
}
