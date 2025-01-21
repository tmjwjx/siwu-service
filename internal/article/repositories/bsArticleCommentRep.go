package repositories

import (
	"fmt"
	"forum/internal/article/requests"
	"forum/internal/internalPkg/internalUtils"
	"forum/internal/models"
	"forum/pkg/globals"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// BatchReviewRep 批量审核
func BatchReviewRep(db *gorm.DB, req *requests.BatchReviewReq) (*requests.BatchReviewRes, error) {

	var comments []models.ArticleComment

	// 通过id查找对应的评论
	err := db.Where("id IN ?", req.IDs).Find(&comments).Error
	if err != nil {
		return nil, fmt.Errorf("BatchReviewRep -> 批量审核查询对应评论失败 -> %s", err)
	}
	if len(comments) == 0 {
		return nil, fmt.Errorf("BatchReviewRep -> 该评论不存在 -> %s", err)
	}

	// 改变查到的评论的审核状态
	for _, comment := range comments {
		err = db.Model(&comment).Update("Examine", 1).Error
		if err != nil {
			return nil, fmt.Errorf("BatchReviewRep -> 变查到的评论的审核状态失败 -> %s", err)
		}
	}

	res := &requests.BatchReviewRes{}
	res.ComList2 = append(res.ComList2, struct{}{})

	return res, nil

}

// ShowCommentsListRep 展示评论列表(获取评论列表)
func ShowCommentsListRep(db *gorm.DB, req *requests.CommentsListReq) (*requests.CommentsListRes, error) {

	if req.Limit <= 0 || req.Offset < 0 {
		return nil, fmt.Errorf("ShowCommentsListRep -> Limit的值不能小于等于0 或者 Offset的值不能小于0")
	}

	var commentsListRes *requests.CommentsListRes
	var comments []models.ArticleComment
	var comList []requests.ComList
	var user models.User
	var article models.Article
	examines := make([]int, 2)

	query := db.Model(&models.ArticleComment{}).Joins("left join sw_users on sw_users.id = sw_article_comments.user_id").
		Joins("left join sw_articles on sw_articles.id = sw_article_comments.article_id")

	if req.Type == 1 {
		examines = append(examines, 0, 1)
	} else if req.Type == 2 {
		examines = append(examines, 0)
	} else if req.Type == 3 {
		examines = append(examines, 1)
	} else {
		return nil, fmt.Errorf("ShowCommentsListRep -> type的值不是规定值 1, 2, 3")
	}

	// 添加查询条件
	if req.Email != "" {
		query = query.Where("sw_users.email = ?", req.Email)
	}
	if req.Nickname != "" {
		query = query.Where("sw_users.nickname LIKE ?", "%"+req.Nickname+"%")
	}
	if req.Title != "" {
		query = query.Where("sw_articles.title LIKE ?", "%"+req.Title+"%")
	}
	if req.ParentEmail != "" {
		query = query.Where("sw_article_comments.parent_email = ?", req.ParentEmail)
	}

	//else {
	//	// 查询评论信息
	//	err := db.Where("examine = ?", examine).Limit(req.Limit).Offset(req.Offset).Find(&comments).Error
	//	if err != nil {
	//		return nil, fmt.Errorf("ShowCommentsListRep -> 查询评论信息失败 -> %s", err)
	//	}
	//}

	//if req.Email != "" || req.Nickname != "" || req.Title != "" || req.ParentEmail != "" {
	//	err := query.Where("examine = ?", examine).Limit(req.Limit).Offset(req.Offset).Find(&comments).Error
	//	if err != nil {
	//		return nil, fmt.Errorf("ShowCommentsListRep -> 查询评论信息失败 -> %s", err)
	//	}
	//}

	err := query.Where("examine IN ?", examines).Scan(&comments).Error
	if err != nil {
		return nil, fmt.Errorf("ShowCommentsListRep -> 查询评论信息失败 -> %s", err)
		//commentsListRes = &requests.CommentsListRes{
		//	Comlist: make([]requests.ComList, 0),
		//	Total:   0,
		//}
		//
		//return commentsListRes, nil
	}

	// 获取返回评论的总数目
	length := len(comments)

	for _, comment := range comments {

		// 查询用户信息
		// 这里要把 user 结构体中存储的上次的查询结果，清空一下，否则会影响下次的查询
		user = models.User{}

		err = db.Where("id = ?", comment.UserID).First(&user).Error
		if err != nil {
			return nil, fmt.Errorf("ShowCommentsListRep -> 查询用户信息失败 -> %s", err)
			//commentsListRes = &requests.CommentsListRes{
			//	Comlist: make([]requests.ComList, 0),
			//	Total:   0,
			//}
			//
			//return commentsListRes, nil
		}

		// 查询文章信息
		// 这里要把 article 结构体中存储的上次的查询结果，清空一下，否则会影响下次的查询
		article = models.Article{}

		err = db.Where("id = ?", comment.ArticleID).First(&article).Error
		if err != nil {
			return nil, fmt.Errorf("ShowCommentsListRep -> 查询文章信息失败 -> %s", err)
			//commentsListRes = &requests.CommentsListRes{
			//	Comlist: make([]requests.ComList, 0),
			//	Total:   0,
			//}
			//
			//return commentsListRes, nil
		}

		commentRes := requests.ComList{
			ID:        comment.ID,
			Nickname:  user.Nickname,
			Email:     user.Email,
			ArticleID: comment.ArticleID,
			Content:   comment.Content,
			Title:     article.Title,
			Summary:   article.Summary,
		}

		// 这里要把 user 结构体中存储的上次的查询结果，清空一下，否则会影响下次的查询
		user = models.User{}

		// 查询用户回复对象信息
		d := db.Where("id = ?", comment.ParentUserID).First(&user)
		if d.Error != nil {
			return nil, fmt.Errorf("ShowCommentsListRep -> 查询用户回复对象信息失败 -> %s", err)
			//commentsListRes = &requests.CommentsListRes{
			//	Comlist: make([]requests.ComList, 0),
			//	Total:   0,
			//}
			//
			//return commentsListRes, nil
		}

		commentRes.ParentNickname = user.Nickname

		// 查询用户头像
		images, err := internalUtils.GetImages(db, globals.UserHome, comment.UserID)
		if err != nil {
			return nil, fmt.Errorf("ShowCommentsListRep -> %s", err)
		} else {
			for _, path := range *images {
				commentRes.Path = path
			}
		}

		// 查询用户发的评论图片
		images2, err := internalUtils.GetImages(db, globals.CommentHome, comment.ID)
		if err != nil {
			return nil, fmt.Errorf("ShowCommentsListRep -> %s", err)
		} else {
			for _, path := range *images2 {
				commentRes.CommentPath = path
			}
		}

		comList = append(comList, commentRes)

	}

	// 分页返回数据
	page := (req.Offset - 1) * req.Limit
	if length > page {
		end := page + req.Limit
		if end > length {
			end = length
		}
		comList = comList[page:end]
		commentsListRes = &requests.CommentsListRes{
			Comlist: comList,
			Total:   length,
		}
		return commentsListRes, nil
	}

	commentsListRes = &requests.CommentsListRes{
		Comlist: make([]requests.ComList, 0),
		Total:   0,
	}

	return commentsListRes, nil
}

// AddCommentRep 添加评论
func AddCommentRep(c *gin.Context, db *gorm.DB, req *requests.AddCommentReq) (error, int) {

	// 开启事务
	tx := db.Begin()
	if tx.Error != nil {
		return fmt.Errorf("AddCommentRep -> 开启事务失败 -> %s", tx.Error), 500
	}

	// 创建要插入的数据模型
	comment := &models.ArticleComment{
		ArticleID:    req.ArticleID,
		UserID:       req.UserID,
		HighestID:    req.HighestID,
		ParentID:     req.ParentID,
		ParentUserID: req.ParentUserID,
		Content:      req.Content,
		LikesCount:   req.LikesCount,
	}

	// 插入数据
	err := tx.Create(comment).Error
	if err != nil {
		tx.Rollback() // 回滚事务
		return fmt.Errorf("AddCommentRep -> 添加评论失败 -> %s", err), 500
	}

	// 获取新插入的评论的ID
	err = tx.First(comment).Error
	if err != nil {
		tx.Rollback() // 回滚事务
		return fmt.Errorf("AddCommentRep -> 获取新插入的评论的ID失败 -> %s", err), 500
	}

	//// 将评论的图片存到文件系统中
	//err, status := controllers.UploadImagesControllers(c, "评论", comment.ID)
	//if err != nil {
	//	return err, status
	//}

	u := &internalUtils.UrlParam{
		UrlPath: req.CommentPath,
		Home:    globals.CommentHome,
		HomeID:  comment.ID,
		DB:      db,
	}
	err = internalUtils.StoreUrl(u)
	if err != nil {
		return fmt.Errorf("AddTagRep -> 存储图片的相关信息失败 -> %s", err), 500
	}

	// 提交事务
	err = tx.Commit().Error
	if err != nil {
		return fmt.Errorf("AddCommentRep -> 提交事务失败 -> %s", err), 500
	}

	return nil, 200

}

// BsDeleteCommentRep 删除评论
func BsDeleteCommentRep(db *gorm.DB, req *requests.DelCommentReq) error {

	// 开启事务
	tx := db.Begin()
	if tx.Error != nil {
		return fmt.Errorf("BsDeleteCommentRep -> 开启事务失败 -> %s", tx.Error)
	}

	// 删除评论
	result := tx.Model(&models.ArticleComment{}).Where("id = ?", req.ID).Delete(nil)
	if result.Error != nil {
		tx.Rollback() // 回滚事务
		return fmt.Errorf("BsDeleteCommentRep -> 删除评论失败 -> %s", result.Error)
	} else if result.RowsAffected == 0 {
		tx.Rollback() // 回滚事务
		return fmt.Errorf("没有找到匹配的记录或记录已经被删除")
	}

	// 提交事务
	err := tx.Commit().Error
	if err != nil {
		return fmt.Errorf("BsDeleteCommentRep -> 提交事务失败 -> %s", err)
	}

	return nil
}

// BatchDelCommentRep 批量删除
func BatchDelCommentRep(db *gorm.DB, req *requests.BsBatchDelCommentReq) error {

	// 开启事务
	tx := db.Begin()
	if tx.Error != nil {
		return fmt.Errorf("BatchDelCommentRep1 -> 开启事务失败 -> %s", tx.Error)
	}

	// 删除评论
	for _, id := range req.ID {

		// 之后如果没删除成功，就将没删除成功的标签名返回给后台

		result := tx.Model(&models.ArticleComment{}).Where("id = ?", id).Delete(nil)
		if result.Error != nil {
			// 回滚事务
			tx.Rollback()
			return fmt.Errorf("BatchDelCommentRep2 -> 评论批量删除失败 -> %s", result.Error)
		} else if result.RowsAffected == 0 {
			// 回滚事务
			tx.Rollback()
			return fmt.Errorf("BatchDelCommentRep3 -> 没有找到匹配的记录或记录已经被删除")
		}

		// 删除文件系统中的图片
		err := internalUtils.DeleteFile("评论", id)
		if err != nil {
			return fmt.Errorf("BatchDelCommentRep4 -> 删除文件系统中的图片失败 -> %s", err)
		}

	}

	// 提交事务
	err := tx.Commit().Error
	if err != nil {
		return fmt.Errorf("BatchDelCommentRep5 -> 提交事务失败 -> %s", err)
	}

	return nil

}

// UpdateCommentRep 更新评论
func UpdateCommentRep(c *gin.Context, db *gorm.DB, req *requests.UpdateCommentReq) (error, int) {

	// 开启事务
	tx := db.Begin()
	if tx.Error != nil {
		return fmt.Errorf("UpdateCommentRep -> 开启事务失败 -> %s", tx.Error), 500
	}

	// 更新评论内容
	err := tx.Model(models.ArticleComment{}).Where("id = ?", req.ID).Update("Content", req.Content).Error
	if err != nil {
		// 回滚事务
		tx.Rollback()
		return fmt.Errorf("UpdateCommentRep -> 更新评论内容失败 -> %s", err), 500
	}

	//// 更新评论图片
	//err, status := controllers.UploadImagesControllers(c, "评论", req.ID)
	//if err != nil {
	//	return fmt.Errorf("UpdateCommentRep -> 更新评论图片失败 -> %s", err), status
	//}
	u := &internalUtils.UrlParam{
		UrlPath: req.CommentPath,
		Home:    globals.CommentHome,
		HomeID:  req.ID,
		DB:      db,
	}
	err = internalUtils.StoreUrl(u)
	if err != nil {
		return fmt.Errorf("AddTagRep -> 存储图片的相关信息失败 -> %s", err), 500
	}

	// 提交事务
	err = tx.Commit().Error
	if err != nil {
		return fmt.Errorf("UpdateCommentRep -> 提交事务失败 -> %s", err), 500
	}

	return nil, 200
}

// QueryCommentRep 查询某个用户的全部评论
func QueryCommentRep(db *gorm.DB, req *requests.QueryCommentReq) (*[]*requests.QueryCommentRes, error) {

	var queryCommentRes []*requests.QueryCommentRes
	var comments []models.ArticleComment
	var user models.User
	var article models.Article

	// 查询评论信息
	err := db.Where("user_id = ?", req.UserID).Limit(req.Limit).Offset(req.Offset).Find(&comments).Error
	if err != nil {
		return nil, fmt.Errorf("QueryCommentRep -> %s", err)
	}

	for _, comment := range comments {

		// 查询用户信息
		// 这里要把 user 结构体中存储的上次的查询结果，清空一下，否则会影响下次的查询
		user = models.User{}

		err := db.Where("id = ?", comment.UserID).First(&user).Error
		if err != nil {
			return nil, fmt.Errorf("QueryCommentRep -> 查询用户信息失败 -> %s", err)
		}

		// 查询文章信息
		// 这里要把 article 结构体中存储的上次的查询结果，清空一下，否则会影响下次的查询
		article = models.Article{}

		err = db.Where("id = ?", comment.ArticleID).First(&article).Error
		if err != nil {
			return nil, fmt.Errorf("QueryCommentRep -> 查询文章信息失败 -> %s", err)
		}

		commentRes := &requests.QueryCommentRes{
			ID:        comment.ID,
			Nickname:  user.Nickname,
			Email:     user.Email,
			ArticleID: comment.ArticleID,
			Content:   comment.Content,
			Title:     article.Title,
			Summary:   article.Summary,
		}

		// 这里要把 user 结构体中存储的上次的查询结果，清空一下，否则会影响下次的查询
		user = models.User{}

		// 查询用户回复对象信息
		err = db.Where("id = ?", comment.ParentUserID).First(&user).Error
		if err != nil {
			return nil, fmt.Errorf("QueryCommentRep -> 查询用户回复对象信息失败 -> %s", err)
		}

		commentRes.ParentNickname = user.Nickname

		// 查询用户头像
		images, err := internalUtils.GetImages(db, globals.UserHome, comment.UserID)
		if err != nil {
			return nil, fmt.Errorf("QueryCommentRep -> %s", err)
		} else {
			for _, path := range *images {
				commentRes.Path = path
			}
		}

		// 查询用户发的评论图片
		images2, err := internalUtils.GetImages(db, globals.CommentHome, comment.ID)
		if err != nil {
			return nil, fmt.Errorf("QueryCommentRep -> %s", err)
		} else {
			for _, path := range *images2 {
				commentRes.CommentPath = path
			}
		}

		queryCommentRes = append(queryCommentRes, commentRes)

	}

	return &queryCommentRes, nil

}
