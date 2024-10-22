package repositories

import (
	"fmt"
	"forum/internal/internalPkg/internalUtils"
	"forum/internal/message/requests"
	"forum/internal/models"
	"forum/pkg/globals"
	"gorm.io/gorm"
)

// CommentMesRep 评论消息(前台)
func CommentMesRep(db *gorm.DB, req *requests.MessageReq, id any) (*requests.CommentMesRes, error) {

	var commentList []*requests.CommentObj
	var result1 []requests.Result1
	err := db.Model(&models.ArticleComment{}).Joins("join sw_user on sw_user.id = sw_article_comment.use_id").
		Joins("join sw_article on sw_article.id = sw_article_comment.article_id").Where("parent_user_id = ?", id).
		Select("sw_user.nickname, sw_user.id as user_id, sw_article.id As article_id, sw_article.title, sw_article_comment.id As comment_id, sw_article_comment.content, sw_article_comment.parent_id, sw_article_comment.created_at, sw_article_comment.likes_count").
		Limit(req.Limit).Offset(req.Page).Scan(&result1).Error
	if err != nil {
		return nil, fmt.Errorf(" CommentMesRep -> 评论消息(前台) -> %s", err)
	}

	for _, res1 := range result1 {

		// 转换一下时间格式
		pastTime := internalUtils.TimeAgo(res1.CreatedAt)

		res := &requests.CommentObj{
			UserID:     res1.UserID,
			CommentID:  res1.CommentID,
			ArticleID:  res1.ArticleID,
			Nickname:   res1.Nickname,
			Title:      res1.Title,
			Content:    res1.Content,
			CreatedAt:  pastTime,
			LikesCount: res1.LikesCount,
		}

		var commentLike models.CommentLike
		err = db.Model(models.CommentLike{}).Where("comment_id = ? and user_id = ?", res1.CommentID, id).First(&commentLike).Error
		if err != nil {
			res.Status = 0
		} else {
			res.Status = 1
		}

		// 判断回复是不是对作者文章的评论
		var count int64
		err = db.Model(models.Article{}).Where("id = ? and user_id = ?", res1.ArticleID, id).Count(&count).Error
		if count > 0 {
			res.Comment = ""
		} else {
			var result2 requests.Result2
			err = db.Model(models.ArticleComment{}).Where("id = ?", res1.ParentID).Scan(&result2).Error
			if err != nil {
				return nil, fmt.Errorf(" CommentMesRep -> 评论消息(前台) -> %s", err)
			}

			res.Comment = result2.Content
		}

		// 查询用户头像
		images, err := internalUtils.GetImages(db, globals.UserHome, res1.UserID)
		if err != nil {
			return nil, fmt.Errorf(" CommentMesRep -> %s", err)
		} else {
			for _, path := range *images {
				res.Path = path
			}
		}

		commentList = append(commentList, res)

	}

	commentMesRes := &requests.CommentMesRes{
		CommentList: &commentList,
	}

	return commentMesRes, nil

}
