package logics

import (
	"forum/internal/message/repositories"
	"forum/internal/message/requests"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

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
