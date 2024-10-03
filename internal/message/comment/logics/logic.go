package logics

import (
	"forum/internal/message/comment/repositories"
	"forum/internal/message/comment/requests"
	"gorm.io/gorm"
)

// CommentMesLogic 评论消息(前台)
func CommentMesLogic(db *gorm.DB, req *requests.CommentMesReq, id any) (*requests.CommentMesRes, error) {
	res, err := repositories.CommentMesRep(db, req, id)
	return res, err
}
