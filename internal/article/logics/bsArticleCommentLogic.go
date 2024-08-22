package logics

import (
	"forum/internal/article/repositories"
	"forum/internal/article/requests"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// ShowCommentsListLogic 展示评论列表
func ShowCommentsListLogic(db *gorm.DB, req *requests.CommentsListReq) (*[]*requests.CommentsListRes, error) {
	commentsListRes, err := repositories.ShowCommentsListRep(db, req)
	return commentsListRes, err
}

// AddCommentLogic 添加评论
func AddCommentLogic(c *gin.Context, db *gorm.DB, req *requests.AddCommentReq) (error, int) {
	err, status := repositories.AddCommentRep(c, db, req)
	return err, status
}

// BsDeleteCommentLogic 删除评论
func BsDeleteCommentLogic(db *gorm.DB, req *requests.DelCommentReq) error {
	err := repositories.BsDeleteCommentRep(db, req)
	return err
}

// BatchDelCommentLogic 批量删除
func BatchDelCommentLogic(db *gorm.DB, req *requests.BsBatchDelCommentReq) error {
	err := repositories.BatchDelCommentRep(db, req)
	return err
}

// UpdateCommentLogic 更新评论
func UpdateCommentLogic(c *gin.Context, db *gorm.DB, req *requests.UpdateCommentReq) (error, int) {
	err, status := repositories.UpdateCommentRep(c, db, req)
	return err, status
}

// QueryCommentLogic 查询某个用户的全部评论
func QueryCommentLogic(db *gorm.DB, req *requests.QueryCommentReq) (*[]*requests.QueryCommentRes, error) {
	queryCommentRes, err := repositories.QueryCommentRep(db, req)
	return queryCommentRes, err
}
