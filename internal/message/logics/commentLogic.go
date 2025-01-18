package logics

import (
	"forum/internal/message/repositories"
	"forum/internal/message/requests"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// CommentMesLogic
// @Description: 评论消息
// @Author tianjiajie 2025-01-18 11:11:34
func CommentMesLogic(db *gorm.DB, req *requests.MessageReq, userId uint) (data interface{}, err error) {
	// 评论消息
	res, err := repositories.CommentRep(db, req, userId)
	if err != nil {
		return nil, err
	}
	// 查询是否点赞
	for i, v := range res {
		like, err := repositories.IsCommentLikeRep(db, userId, v.CommentId)
		if err != nil {
			return nil, err
		}
		if like {
			res[i].Status = 1
		} else {
			res[i].Status = 0
		}
	}
	data = gin.H{"comment_list": res}
	return data, err
}
