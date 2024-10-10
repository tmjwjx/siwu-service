package logics

import (
	"forum/internal/message/repositories"
	"forum/internal/message/requests"
	"gorm.io/gorm"
)

// CommentMesLogic 评论消息(前台)
func CommentMesLogic(db *gorm.DB, req *requests.MessageReq, id any) (*requests.CommentMesRes, error) {
	res, err := repositories.CommentMesRep(db, req, id)
	return res, err
}