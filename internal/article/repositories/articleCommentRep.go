package repositories

import (
	"fmt"
	"forum/internal/article/requests"
	"forum/internal/image/controllers"
	"forum/internal/internal_pkg/internal_utils"
	"forum/internal/models"
	"forum/pkg/globals"
	"gorm.io/gorm"
)

// InsertCommentRep 将评论存入数据库中
func InsertCommentRep(articleCommentReq *requests.ArticleCommentReq, db *gorm.DB) error {

	// 开启事务
	tx := db.Begin()
	if tx.Error != nil {
		return fmt.Errorf("InsertCommentRep -> %s", tx.Error)
	}

	// 构建要插入的结构体
	articleComment := &models.ArticleComment{
		ArticleID: articleCommentReq.ArticleID,
		UserID:    articleCommentReq.UserID,
		HighestID: articleCommentReq.HighestID,
		ParentID:  articleCommentReq.ParentID,
		Content:   articleCommentReq.Content,
	}

	// 插入评论
	err := tx.Create(&articleComment).Error
	if err != nil {
		return fmt.Errorf("InsertCommentRep -> %s", err)
	}

	// 提交事务
	err = tx.Commit().Error
	if err != nil {
		return fmt.Errorf("InsertCommentRep -> %s", err)
	}
	return nil
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
func GetTopLevelCommentsRep(req *requests.TopCommentsReq) (*[]requests.TopCommentsRes, error) {
	var topCommentsRes []requests.TopCommentsRes
	var count int64

	// 查询顶级评论
	err := globals.DB.Model(&models.ArticleComment{}).Where("article_id = ? AND parent_id IS NULL", req.ArticleID).
		Order("created_at asc").Limit(req.Limit).Offset(req.Offset).Find(&topCommentsRes).Error
	if err != nil {
		return nil, fmt.Errorf("GetTopLevelCommentsRep -> %s", err)
	}

	for _, comment := range topCommentsRes {

		// 统计每个顶级评论的回复数量
		err := globals.DB.Model(&models.ArticleComment{}).Where("highest_id = ?", comment.ID).Count(&count).Error
		if err != nil {
			return nil, fmt.Errorf("GetTopLevelCommentsRep -> %s", err)
		}
		comment.RepliesCount = count

		// 查询用户名字
		err = globals.DB.Model(&models.User{}).Select("Nickname").Where("id = ?", comment.UserID).First(&comment.Nickname).Error
		if err != nil {
			return nil, fmt.Errorf("GetTopLevelCommentsRep -> %s", err)
		}

		// 查询用户头像
		images, err := controllers.GetImagesControllers("用户", comment.ID)
		if err != nil {
			comment.Path = internal_utils.UserDefaultImage
		}
		for _, image := range *images {
			comment.Path = image.Path
		}

		// 用户发的评论中的图片
		CommentImages, err := controllers.GetImagesControllers("评论", comment.ID)
		if err != nil {
			comment.CommentPath = ""
		}

		for _, image := range *CommentImages {
			comment.CommentPath = image.Path
		}

	}

	return &topCommentsRes, nil

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
func GetRepliesRep2Rep(req *requests.RepliesReq2) (*[]requests.RepliesRes, error) {
	var repliesRes []requests.RepliesRes

	err := globals.DB.Where("highest_id = ?", req.HighestID).Order("create_at asc").Limit(req.Limit).Offset(req.Offset).Find(&repliesRes).Error
	if err != nil {
		return nil, fmt.Errorf("GetRepliesRep2Rep -> %s", err)
	}

	for _, comment := range repliesRes {

		// 查询用户名字
		err := globals.DB.Model(&models.User{}).Select("Nickname").Where("id = ?", comment.UserID).First(&comment.Nickname).Error
		if err != nil {
			return nil, fmt.Errorf("GetRepliesRep2Rep -> %s", err)
		}

		// 查询用户头像
		images, err := controllers.GetImagesControllers("用户", comment.ID)
		if err != nil {
			comment.ParentPath = internal_utils.UserDefaultImage
		}

		for _, image := range *images {
			comment.Path = image.Path
		}

		// 查询用户回复对象的名字
		err = globals.DB.Model(&models.User{}).Select("Nickname").Where("id = ?", comment.ParentID).First(&comment.ParentNickname).Error
		if err != nil {
			return nil, fmt.Errorf("GetRepliesRep2Rep -> %s", err)
		}

		// 查询用户回复对象的头像
		images2, err := controllers.GetImagesControllers("用户", *comment.ParentID)
		if err != nil {
			comment.ParentPath = internal_utils.UserDefaultImage
		}

		for _, image := range *images2 {
			comment.ParentPath = image.Path
		}

		// 用户发的评论中的图片
		CommentImages, err := controllers.GetImagesControllers("评论", comment.ID)
		if err != nil {
			comment.ParentPath = ""
		}

		for _, image := range *CommentImages {
			comment.CommentPath = image.Path
		}

	}

	return &repliesRes, nil

}

// DeleteCommentRep 删除评论
func DeleteCommentRep(req *requests.DelComment, db *gorm.DB) error {
	// 开启事务
	tx := db.Begin()
	if tx.Error != nil {
		return fmt.Errorf("DeleteCommentRep -> %s", tx.Error)
	}

	// 删除评论
	err := tx.Model(&models.ArticleComment{}).Where("id = ?", req.ID).Delete(nil).Error
	if err != nil {
		tx.Rollback() // 回滚事务
		return fmt.Errorf("DeleteCommentRep -> %s", err)
	}

	//提交事务
	err = tx.Commit().Error
	if err != nil {
		return fmt.Errorf("DeleteCommentRep -> %s", err)
	}

	return nil
}

// UpdatePraiseCountRep 更新点赞的数量
func UpdatePraiseCountRep(req *requests.PraiseCount, db *gorm.DB) error {
	// 开启事务
	tx := db.Begin()
	if tx.Error != nil {
		return fmt.Errorf("UpdatePraiseCountRep -> %s", tx.Error)
	}

	// 更新评论的点赞数量
	if req.Status == 1 {
		err := tx.Model(&models.ArticleComment{}).Where("id = ?", req.ID).UpdateColumn("likes_count", gorm.Expr("likes_count + ?", 1)).Error
		if err != nil {
			return fmt.Errorf("UpdatePraiseCountRep -> %s", err)
		}
	} else if req.Status == 2 {
		err := tx.Model(&models.ArticleComment{}).Where("id = ?", req.ID).UpdateColumn("likes_count", gorm.Expr("likes_count - ?", 1)).Error
		if err != nil {
			return fmt.Errorf("UpdatePraiseCountRep -> %s", err)
		}
	}

	//提交事务
	err := tx.Commit().Error
	if err != nil {
		return fmt.Errorf("DeleteCommentRep -> %s", err)
	}
	return nil
}
