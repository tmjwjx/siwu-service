package logics

import (
	"fmt"
	"forum/internal/article/repositories"
	"forum/internal/article/requests"
	"forum/pkg/globals"
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
	if req.Offset <= 0 || req.Limit <= 0 {
		return nil, fmt.Errorf("GetTopLevelCommentsLogic -> Offset 或 Limit 的值不能小于或等于0")
	}

	topCommentsRes, err := repositories.GetTopLevelCommentsRep(userId, db, req)
	if err != nil {
		return nil, err
	}

	// 查询点赞状态
	for i := 0; i < len(topCommentsRes.FirstCommentsList); i++ {
		status, err := repositories.GetCommentStatusRep(db, userId, topCommentsRes.FirstCommentsList[i].ID)
		if err != nil {
			return topCommentsRes, err
		}
		globals.Log.Info("id: ", topCommentsRes.FirstCommentsList[i].ID)
		globals.Log.Info("status: ", topCommentsRes.FirstCommentsList[i].Status)
		globals.Log.Info("status: ", status)
		topCommentsRes.FirstCommentsList[i].Status = status
		globals.Log.Info("status: ", topCommentsRes.FirstCommentsList[i].Status)
	}

	return topCommentsRes, nil
}

// GetRepliesRep2Logic 返回评论回复
func GetRepliesRep2Logic(userId uint, db *gorm.DB, req *requests.RepliesReq2) (*requests.RepliesRes, error) {
	if req.Offset <= 0 || req.Limit <= 0 {
		return nil, fmt.Errorf("GetRepliesRep2Logic -> Offset 或 Limit 的值不能小于或等于0")
	}

	repliesRes, err := repositories.GetRepliesRep2Rep(userId, db, req)

	if err != nil {
		return nil, err
	}

	// 查询点赞状态
	for i := 0; i < len(repliesRes.SecondCommentsList); i++ {
		status, err := repositories.GetCommentStatusRep(db, userId, repliesRes.SecondCommentsList[i].ID)
		if err != nil {
			return repliesRes, err
		}
		globals.Log.Info("id: ", repliesRes.SecondCommentsList[i].ID)
		globals.Log.Info("status: ", repliesRes.SecondCommentsList[i].Status)
		globals.Log.Info("status: ", status)
		repliesRes.SecondCommentsList[i].Status = status
		globals.Log.Info("status: ", repliesRes.SecondCommentsList[i].Status)
	}

	return repliesRes, err
}

// DeleteCommentLogic 删除评论
func DeleteCommentLogic(req *requests.DelComment, db *gorm.DB) error {
	err := repositories.DeleteCommentRep(req, db)
	return err
}

// UpdatePraiseCountLogic 更新点赞的数量
func UpdatePraiseCountLogic(req *requests.PraiseCount, db *gorm.DB, userId uint) error {
	err := repositories.UpdatePraiseCountRep(req, db, userId)
	return err
}
