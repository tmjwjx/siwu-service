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
	res, err := repositories.CommentRep(db, req, userId)
	if err != nil {
		return nil, err
	}
	data = gin.H{"comment_list": res}
	return data, err
}
