package logics

import (
	"forum/internal/article/repositories"
	"forum/internal/article/requests"
	"gorm.io/gorm"
)

// InsertCommentLogic 将评论存入数据库中
func InsertCommentLogic(userId uint, articleCommentReq *requests.ArticleCommentReq, db *gorm.DB) (error, int) {
	err, status := repositories.InsertCommentRep(userId, articleCommentReq, db)
	return err, status
}

/*func GetCommentByArticleLogic(articleID string) (*[]requests.ArticleCommentRes, error) {
	topLevelComments, err := repositories.GetCommentByArticleRep(articleID)
	return topLevelComments, err
}*/

/*func GetCommentsLogic(req *requests.CommentReq) (*[]models.ArticleComment, error) {
	articleComments, err := repositories.GetCommentsRep(req)
	return articleComments, err
}

func GetRepliesLogic(req *requests.RepliesReq) (*[]models.ArticleComment, error) {
	articleComments, err := repositories.GetRepliesRep(req)
	return articleComments, err
}*/

// GetTopLevelCommentsLogic 返回顶级评论
func GetTopLevelCommentsLogic(userId uint, db *gorm.DB, req *requests.TopCommentsReq) (*requests.TopCommentsRes, error) {
	topCommentsRes, err := repositories.GetTopLevelCommentsRep(userId, db, req)
	return topCommentsRes, err
}

// GetRepliesRep2Logic 返回评论回复
func GetRepliesRep2Logic(userId uint, db *gorm.DB, req *requests.RepliesReq2) (*requests.RepliesRes, error) {
	repliesRes, err := repositories.GetRepliesRep2Rep(userId, db, req)
	return repliesRes, err
}

// DeleteCommentLogic 删除评论
func DeleteCommentLogic(req *requests.DelComment, db *gorm.DB) error {
	err := repositories.DeleteCommentRep(req, db)
	return err
}

// UpdatePraiseCountLogic 更新点赞的数量
func UpdatePraiseCountLogic(req *requests.PraiseCount, db *gorm.DB) error {
	err := repositories.UpdatePraiseCountRep(req, db)
	return err
}
