package logics

import (
	"forum/internal/article/repositories"
	"forum/internal/article/requests"
	"gorm.io/gorm"
)

// InsertCommentLogic 将评论存入数据库中
func InsertCommentLogic(articleCommentReq *requests.ArticleCommentReq, db *gorm.DB) error {
	err := repositories.InsertCommentRep(articleCommentReq, db)
	return err
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
func GetTopLevelCommentsLogic(req *requests.TopCommentsReq) (*[]*requests.TopCommentsRes, error) {
	topCommentsRes, err := repositories.GetTopLevelCommentsRep(req)
	return topCommentsRes, err
}

// GetRepliesRep2Logic 返回评论回复
func GetRepliesRep2Logic(req *requests.RepliesReq2) (*[]*requests.RepliesRes, error) {
	repliesRes, err := repositories.GetRepliesRep2Rep(req)
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
