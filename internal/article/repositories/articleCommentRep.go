package repositories

import (
	"fmt"
	"forum/internal/article/requests"
	"forum/internal/internalPkg/internalUtils"
	"forum/internal/models"
	"forum/pkg/globals"
	"gorm.io/gorm"
)

// InsertCommentRep 将评论存入数据库中
func InsertCommentRep(userId uint, articleCommentReq *requests.ArticleCommentReq, db *gorm.DB) (error, int) {

	// 开启事务
	tx := db.Begin()
	if tx.Error != nil {
		return fmt.Errorf("InsertCommentRep -> 开启事务失败 -> %s", tx.Error), 500
	}

	// 构建要插入的结构体
	articleComment := &models.ArticleComment{
		ArticleID: articleCommentReq.ArticleID,
		UserID:    userId,
		HighestID: articleCommentReq.HighestID,
		ParentID:  articleCommentReq.ParentID,
		Content:   articleCommentReq.Content,
	}

	// 插入评论
	err := tx.Create(&articleComment).Error
	if err != nil {
		tx.Rollback() // 回滚事务
		return fmt.Errorf("InsertCommentRep -> 插入评论失败 ->  %s", err), 500
	}

	// 获取新插入的评论的ID
	err = tx.First(articleComment).Error
	if err != nil {
		tx.Rollback() // 回滚事务
		return fmt.Errorf("InsertCommentRep -> 获取新插入的评论的ID失败 -> %s", err), 500
	}

	//// 将评论的图片存到文件系统中
	//err, status := controllers.UploadImagesControllers(c, "评论", articleComment.ID)
	//if err != nil {
	//	return err, status
	//}
	u := &internalUtils.UrlParam{
		UrlPath: articleCommentReq.Path,
		Home:    globals.CommentHome,
		HomeID:  articleComment.ID,
		DB:      db,
	}
	err = internalUtils.StoreUrl(u)
	if err != nil {
		return fmt.Errorf("AddTagRep -> 存储图片的相关信息失败 -> %s", err), 500
	}

	// 提交事务
	err = tx.Commit().Error
	if err != nil {
		return fmt.Errorf("InsertCommentRep -> 提交事务失败 -> %s", err), 500
	}
	return nil, 200

}

/*func GetCommentByArticleRep(articleID string) (*[]requests.ArticleCommentRes, error) {
	var articleComments []models.ArticleComment
	err := globals.DB.Where("article_id = ?", articleID).Order("created_at asc").Find(&articleComments).Error
	if err != nil {
		return nil, fmt.Errorf("GetCommentByArticleRep -> %s", err)
	}
	commentMap := make(map[uint][]requests.ArticleCommentRes)
	var topLevelComments []requests.ArticleCommentRes

	// 构建评论的Map, 方便后续嵌套
	for _, comment := range articleComments {
		response := requests.ArticleCommentRes{
			ID:        comment.ID,
			ArticleID: comment.ArticleID,
			UserID:    comment.UserID,
			Content:   comment.Content,
			ParentID:  comment.ParentID,
			CreateAT:  comment.CreatedAt,
		}
		if comment.ParentID == nil {
			topLevelComments = append(topLevelComments, response)
		} else {
			commentMap[*comment.ParentID] = append(commentMap[*comment.ParentID], response)
		}
	}

	// 递归地构建评论树
	for i := range topLevelComments {
		buildCommentTree(&topLevelComments[i], commentMap)
	}
	return &topLevelComments, nil
}

func buildCommentTree(comment *requests.ArticleCommentRes, commentMap map[uint][]requests.ArticleCommentRes) {
	if replies, ok := commentMap[comment.ID]; ok {
		comment.Replies = replies
		for i := range comment.Replies {
			buildCommentTree(&comment.Replies[i], commentMap)
		}
	}
}*/

// GetTopLevelCommentsRep 返回顶级评论
func GetTopLevelCommentsRep(userId uint, db *gorm.DB, req *requests.TopCommentsReq) (*requests.TopCommentsRes, error) {
	var firstCommentsList []*requests.FirstComment
	var count int64
	var articleComments []models.ArticleComment
	var commentId []uint

	// 查询顶级评论
	err := db.Where("article_id = ? AND highest_id = ?", req.ArticleID, 0).
		Order("created_at asc").Find(&articleComments).Error
	if err != nil {
		return nil, fmt.Errorf("GetTopLevelCommentsRep -> %s", err)
	}

	length := len(articleComments)

	// 查询用户对该篇文章中的评论的点赞情况
	err = db.Model(models.ArticleComment{}).Joins("join sw_comment_likes on sw_comment_likes.comment_id = sw_article_comments.id").
		Select("sw_article_comments.id").Where("sw_article_comments.article_id = ? and sw_comment_likes.user_id = ?", req.ArticleID, userId).Find(&commentId).Error
	if err != nil {
		return nil, fmt.Errorf("GetTopLevelCommentsRep -> 查询用户对该篇文章中的评论的点赞情况失败 -> %s", err)
	}

	for _, comment := range articleComments {

		// 转换一下时间格式
		pastTime := internalUtils.TimeAgo(comment.CreatedAt)

		topComment := &requests.FirstComment{
			ID:           comment.ID,
			CreateAT:     pastTime,
			ArticleID:    comment.ArticleID,
			UserID:       comment.UserID,
			HighestID:    comment.HighestID,
			ParentID:     comment.ParentID,
			ParentUserID: comment.ParentUserID,
			Content:      comment.Content,
			LikesCount:   comment.LikesCount,
		}
		// 判断用户对评论是否已经点过赞了
		for _, id := range commentId {
			if comment.ID == id {
				// 1 表示对该评论该用户已经点过赞了
				topComment.Status = 1
				break
			}
		}
		// 2 表示对该评论该用户从没有点过赞
		topComment.Status = 2

		firstCommentsList = append(firstCommentsList, topComment)
	}

	for _, comment := range firstCommentsList {

		// 统计每个顶级评论的回复数量
		err = db.Model(&models.ArticleComment{}).Where("highest_id = ?", comment.ID).Count(&count).Error
		if err != nil {
			return nil, fmt.Errorf("GetTopLevelCommentsRep -> %s", err)
		}
		comment.RepliesCount = count

		// 查询用户名字
		err = db.Model(&models.User{}).Select("Nickname").Where("id = ?", comment.UserID).First(&comment.Nickname).Error
		if err != nil {
			return nil, fmt.Errorf("GetTopLevelCommentsRep -> %s", err)
		}

		// 查询用户头像
		images, err := internalUtils.GetImages(db, globals.UserHome, comment.ID)
		if err != nil {
			return nil, fmt.Errorf("GetTopLevelCommentsRep -> %s", err)
		} else {
			for _, path := range *images {
				comment.Path = path
			}
		}

		// 用户发的评论中的图片
		CommentImages, err := internalUtils.GetImages(db, globals.CommentHome, comment.ID)
		if err != nil {
			return nil, fmt.Errorf("GetTopLevelCommentsRep -> %s", err)
		} else {
			for _, path := range *CommentImages {
				comment.CommentPath = path
			}
		}
	}

	// 分页返回数据
	page := (req.Offset - 1) * req.Limit
	if length > page {
		end := page + req.Limit
		if end > length {
			end = length
		}

		firstCommentsList = firstCommentsList[page:end]
		res := &requests.TopCommentsRes{
			FirstCommentsList: firstCommentsList,
			CommentsTotal:     length,
		}

		return res, nil
	}

	topCommentsRes := &requests.TopCommentsRes{
		FirstCommentsList: make([]*requests.FirstComment, 0),
		CommentsTotal:     len(firstCommentsList),
	}

	return topCommentsRes, nil
}

/*func GetCommentsRep(req *requests.CommentReq) (*[]models.ArticleComment, error) {
	var articleComments []models.ArticleComment
	limit := 10   // 分批加载顶级评论的数量
	maxDepth := 2 // 预加载子评论的最大层级

	// 使用 GORM 的 Preload 方法与预加载顶级评论的子评论
	query := globals.DB.Where("article_id = ? AND parent_id IS NULL", req.ArticleID).Order("created_at asc").Limit(limit).Offset(req.Offset)

	// 根据最大深度，逐层预加载子评论
	for i := 0; i < maxDepth; i++ {
		preloadQuery := fmt.Sprintf("Replies.%s", strings.Repeat("Replies.", i))
		query = query.Preload(preloadQuery)
	}
	err := query.Find(&articleComments).Error
	if err != nil {
		return nil, fmt.Errorf("GetCommentsRep -> %s", err)
	}
	return &articleComments, nil
}

func GetRepliesRep(req *requests.RepliesReq) (*[]models.ArticleComment, error) {
	var articleComments []models.ArticleComment
	err := globals.DB.Where("parent_id = ?", req.ParentID).Order("create_at asc").Limit(req.Limit).Offset(req.Offset).Find(&articleComments).Error
	if err != nil {
		return nil, fmt.Errorf("GetRepliesRep -> %s", err)
	}
	return &articleComments, nil
}*/

// GetRepliesRep2Rep 返回评论回复
func GetRepliesRep2Rep(userId uint, db *gorm.DB, req *requests.RepliesReq2) (*requests.RepliesRes, error) {
	var secondCommentsList []*requests.SecondComment
	var articleComments []models.ArticleComment
	var commentId []uint

	err := db.Where("highest_id = ?", req.HighestID).Order("created_at asc").Find(&articleComments).Error
	if err != nil {
		return nil, fmt.Errorf("GetRepliesRep2Rep -> %s", err)
	}

	length := len(articleComments)

	// 查询用户对该篇文章中的评论的点赞情况
	err = db.Model(models.ArticleComment{}).Joins("join sw_comment_likes on sw_comment_likes.comment_id = sw_article_comments.id").
		Where("sw_article_comments.highest_id = ? and sw_comment_likes.user_id = ?", req.HighestID, userId).Find(&commentId).Error
	if err != nil {
		return nil, fmt.Errorf("GetTopLevelCommentsRep -> 查询用户对该篇文章中的评论的点赞情况失败 -> %s", err)
	}

	for _, comment := range articleComments {

		// 转换一下时间格式
		pastTime := internalUtils.TimeAgo(comment.CreatedAt)

		replies := &requests.SecondComment{
			ID:           comment.ID,
			CreateAT:     pastTime,
			ArticleID:    comment.ArticleID,
			UserID:       comment.UserID,
			HighestID:    comment.HighestID,
			ParentID:     comment.ParentID,
			ParentUserID: comment.ParentUserID,
			Content:      comment.Content,
			LikesCount:   comment.LikesCount,
		}

		// 判断用户对评论是否已经点过赞了
		for _, id := range commentId {
			if comment.ID == id {
				// 1 表示对该评论该用户已经点过赞了
				replies.Status = 1
				break
			}
		}
		// 2 表示对该评论该用户从没有点过赞
		replies.Status = 2

		secondCommentsList = append(secondCommentsList, replies)
	}

	for _, comment := range secondCommentsList {

		// 查询用户名字
		err := db.Model(&models.User{}).Select("Nickname").Where("id = ?", comment.UserID).First(&comment.Nickname).Error
		if err != nil {
			return nil, fmt.Errorf("GetRepliesRep2Rep -> %s", err)
		}

		// 查询用户头像
		images, err := internalUtils.GetImages(db, globals.UserHome, comment.ID)
		if err != nil {
			return nil, fmt.Errorf("GetRepliesRep2Rep -> %s", err)
		} else {
			for _, path := range *images {
				comment.Path = path
			}
		}

		// 查询用户回复对象的名字
		err = db.Model(&models.User{}).Select("Nickname").Where("id = ?", comment.ParentID).First(&comment.ParentNickname).Error
		if err != nil {
			return nil, fmt.Errorf("GetRepliesRep2Rep -> %s", err)
		}

		// 查询用户回复对象的头像
		images2, err := internalUtils.GetImages(db, globals.UserHome, *comment.ParentUserID)
		if err != nil {
			return nil, fmt.Errorf("GetRepliesRep2Rep -> %s", err)
		} else {
			for _, path := range *images2 {
				comment.ParentPath = path
			}
		}

		// 用户发的评论中的图片
		CommentImages, err := internalUtils.GetImages(db, globals.CommentHome, comment.ID)
		if err != nil {
			return nil, fmt.Errorf("GetRepliesRep2Rep -> %s", err)
		} else {
			for _, path := range *CommentImages {
				comment.CommentPath = path
			}
		}
	}

	// 分页返回数据
	page := (req.Offset - 1) * req.Limit
	if length > page {
		end := page + req.Limit
		if end > length {
			end = length
		}

		secondCommentsList = secondCommentsList[page:end]
		res := &requests.RepliesRes{
			SecondCommentsList: secondCommentsList,
		}
		return res, nil
	}

	repliesRes := &requests.RepliesRes{
		SecondCommentsList: make([]*requests.SecondComment, 0),
	}

	return repliesRes, nil
}

// DeleteCommentRep 删除评论
func DeleteCommentRep(req *requests.DelComment, db *gorm.DB) error {
	// 开启事务
	tx := db.Begin()
	if tx.Error != nil {
		return fmt.Errorf("DeleteCommentRep -> 开启事务失败 -> %s", tx.Error)
	}

	// 删除评论
	result := tx.Model(&models.ArticleComment{}).Where("id = ?", req.ID).Delete(nil)
	if result.Error != nil {
		tx.Rollback() // 回滚事务
		return fmt.Errorf("DeleteCommentRep -> %s", result.Error)
	} else if result.RowsAffected == 0 {
		tx.Rollback() // 回滚事务
		return fmt.Errorf("没有找到匹配的记录或记录已经被删除")
	}

	// 提交事务
	err := tx.Commit().Error
	if err != nil {
		return fmt.Errorf("DeleteCommentRep -> 提交事务失败 -> %s", err)
	}

	return nil
}

// UpdatePraiseCountRep 更新点赞的数量
func UpdatePraiseCountRep(req *requests.PraiseCount, db *gorm.DB) error {

	// 开启事务
	tx := db.Begin()
	if tx.Error != nil {
		return fmt.Errorf("UpdatePraiseCountRep -> 开启事务失败 -> %s", tx.Error)
	}

	// 更新 comment_likes 表中的数据
	commentLike := &models.CommentLike{
		CommentID: req.ID,
		UserID:    req.UserID,
	}

	if req.Status == 1 || req.Status == 2 {

		// 更新评论的点赞数量
		if req.Status == 1 {

			err := tx.Model(&models.ArticleComment{}).Where("id = ?", req.ID).UpdateColumn("likes_count", gorm.Expr("likes_count + ?", 1)).Error
			if err != nil {
				tx.Rollback() // 回滚事务
				return fmt.Errorf("UpdatePraiseCountRep1 -> 更新评论的点赞数量 -> %s", err)
			}

			// 每个用户只能对一个评论点赞一次，所以先查询一下，该用户是否已经点赞过该评论了。
			err = tx.Where("comment_id = ? and user_id = ?", req.ID, req.UserID).First(&commentLike).Error
			if err != nil {
				// 如果没有查询到，说明该用户没对该评论点赞过，可以点赞，否则，直接跳过。
				err := tx.Create(&commentLike).Error
				if err != nil {
					tx.Rollback() // 回滚事务
					return fmt.Errorf("DeleteCommentRep1 -> 更新 comment_likes 表中的数据失败 -> %s", err)
				}
			}

		} else if req.Status == 2 {
			err := tx.Model(&models.ArticleComment{}).Where("id = ?", req.ID).UpdateColumn("likes_count", gorm.Expr("likes_count - ?", 1)).Error
			if err != nil {
				tx.Rollback() // 回滚事务
				return fmt.Errorf("UpdatePraiseCountRep2 -> 更新评论的点赞数量 -> %s", err)
			}

			// 每个用户只能对一个评论点赞一次，所以先查询一下，该用户是否已经点赞过该评论了。
			err = tx.Where("comment_id = ? and user_id = ?", req.ID, req.UserID).First(&commentLike).Error
			if err == nil {
				// 如果查询到了，说明该用户对该评论点赞过，可以删除点赞，否则，直接跳过。
				err := tx.Delete(&commentLike).Error
				if err != nil {
					tx.Rollback() // 回滚事务
					return fmt.Errorf("DeleteCommentRep2 -> 更新 comment_likes 表中的数据失败 -> %s", err)
				}
			}

		}

	} else {
		return fmt.Errorf("UpdatePraiseCountRep -> 更新点赞的数量失败 -> status 的值只能是 1 或 2, 1:代表增加点赞, 2:代表取消点赞")
	}

	// 提交事务
	err := tx.Commit().Error
	if err != nil {
		return fmt.Errorf("DeleteCommentRep -> 提交事务失败 -> %s", err)
	}

	return nil
}
