package repositories

import (
	"fmt"
	"forum/internal/article/requests"
	"forum/internal/internalPkg/internalUtils"
	"forum/internal/models"
	"forum/pkg/globals"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"gorm.io/gorm"
	"strconv"
	"time"
)

// InsertArticlesRep
// @Description: 新建文章
// @param        db *gorm.DB
// @param        req requests.ReqPublish
// @return       int
// @return       error
func InsertArticlesRep(db *gorm.DB, req requests.ReqPublish, userId uint) (int, error) {

	newArticle := models.Article{
		// UserID:     req.UserId,
		// Title:      req.Title,
		// Status:     req.Status,
		// CategoryID: req.CategoryID,
		// Summary:    req.Summary,
		// Content:    req.Content,
		// ImageUrl:   req.ImageUrl,
	}

	// 设置文章ID
	if req.ArticleId != 0 {
		if userId != req.UserId {
			return 0, nil // 如果当前用户不是文章作者
		}
		newArticle.ID = uint(req.ArticleId)

		// 获取数据库文章记录
		db.Find(&newArticle)
	}
	{
		newArticle.UserID = req.UserId
		newArticle.Title = req.Title
		newArticle.Status = req.Status
		newArticle.CategoryID = req.CategoryID
		newArticle.Summary = req.Summary
		newArticle.Content = req.Content
		newArticle.ImageUrl = req.ImageUrl
	}

	// 设置标签
	// 查找传递过来的所有标签
	var tags []models.Tag
	if err := db.Where("id IN ?", req.Tags).Find(&tags).Error; err != nil {
		return 0, err // 如果标签不存在，返回错误
	}
	// 更新文章的标签关联
	// 手动清空当前文章的标签关联，确保旧数据被清除
	if err := db.Model(&newArticle).Association("Tags").Clear(); err != nil {
		return 0, err // 如果清除失败，返回错误
	}
	newArticle.Tags = tags

	// 设置发布时间
	if req.Status == "public" {
		now := time.Now()
		newArticle.PublishedAt = &now
	}
	result := db.Save(&newArticle)
	if result.Error != nil {
		globals.Log.Fatal("创建文章失败:", result.Error)
		return 0, result.Error
	}
	id := int(newArticle.ID)

	return id, nil
}

// UpdatePublishRep 设置文章的发布时间
func UpdatePublishRep(db *gorm.DB, article *models.Article) error {

	*article.PublishedAt = time.Now()

	// 更新数据库中的文章记录
	if err := db.Save(article).Error; err != nil {
		return err // 返回错误
	}

	return nil
}

// QueryTag 返回全部标签
func QueryTag(db *gorm.DB) (tags []models.Tag, err error) {
	// 查询 Tag 表中的所有数据
	if err = db.Find(&tags).Error; err != nil {
		globals.Log.Errorf("err = %s", err)
		return nil, err
	}

	return tags, nil
}

// QueryCategory 返回全部类目
func QueryCategory(db *gorm.DB) (categories []models.Category, err error) {

	// 查询 Category 表中的所有数据
	if err = db.Find(&categories).Error; err != nil {
		globals.Log.Errorf("err = %s", err)
		return nil, err
	}

	return categories, nil

}

// // ArticlesOrder
// // @Description: 选择排序方式 0热度 1时间
// // @param        kind int
// // @return       string
// func ArticlesOrder(kind int) string {
//	var condition string
//	if kind == 0 { // 0 代表按照热度排序
//		condition = "heat DESC"
//	} else if kind == 1 { // 1 代表按照发布时间排序
//		condition = "published_at DESC"
//	}
//	return condition
// }
//
// // emptyFilled
// // @Description: 空接口填充
// // @param        data interface{}
// func emptyFilled(data interface{}) interface{} {
//	if data == nil {
//		data = gin.H{}
//	}
//	return data
// }

// SearchArticlesRep 搜索文章
func SearchArticlesRep(db *gorm.DB, req *requests.ArticleSearchReq) (articles []requests.SearchArticleListRes, err error) {

	condition := internalUtils.ArticlesOrder(req.Kind) // 选择排序方式  0热度 1时间

	query := db.Model(&models.Article{}).Preload("Tags").
		Select("sw_articles.*, sw_users.nickname").
		Joins("LEFT JOIN sw_users ON sw_users.id = sw_articles.user_id")

	// if req.UserId != 0 { // 按用户ID筛选
	//	query = query.Joins("JOIN user_follows ON articles.user_id = user_follows.followed_id").
	//		Where("user_follows.follower_id = ?", req.UserId)
	//	condition = "published_at DESC" // 按照 关注 查询只能按照 时间 排序
	// }
	if req.Keyword != "" { // 按关键词搜索
		query = query.Where("title LIKE ? OR summary LIKE ?", "%"+req.Keyword+"%", "%"+req.Keyword+"%")
	}

	if req.CategoryId != 0 { // 按类目筛选
		query = query.Where("category_id = ?", req.CategoryId)
	}

	// 选择排序方式 时间or热度
	// 先排序 后分页
	query = query.Order(condition)

	// 判断是否分页
	if req.Limit != 0 {
		offset := (req.Page - 1) * req.Limit // 计算当前页的偏移量，用于分页
		query = query.Limit(req.Limit).Offset(offset)
	}

	// 执行查询
	if err = query.Find(&articles).Error; err != nil {
		globals.Log.Errorf("err = %s", err)
		return // 结束函数执行
	}

	return articles, err
}

// SearchArticlesListRep
// @Description: 检索获取已经发布的文章列表
// @param        db *gorm.DB
// @param        req *requests.ArticleListReq
// @return       data
// @return       err
func SearchArticlesListRep(db *gorm.DB, req *requests.ArticleListReq) (data interface{}, err error) {
	// var articleList []requests.ArcList
	var articleList []requests.SearchArticleListRes

	query := db.Model(&models.Article{}).Preload("Tags").
		Joins("LEFT JOIN sw_users ON sw_users.id = sw_articles.user_id").
		Joins("LEFT JOIN sw_article_tags ON sw_article_tags.article_id = sw_articles.id").
		Select("DISTINCT sw_articles.*, sw_users.nickname")
	// query.Select(
	//	"sw_articles.id, " +
	//		"title, " +
	//		"article_condition, " +
	//		"views_count, " +
	//		"likes_count, " +
	//		"collections_count, " +
	//		"comments_count, " +
	//		"sw_articles.heat, " +
	//		"sw_users.nickname, " +
	//		//"sw_tags.id AS tag_id, " +
	//		//"sw_tags.name AS tag_name, "
	//		"sw_articles.published_at, " +
	//		"sw_articles.updated_at")
	// Joins("LEFT JOIN sw_users ON sw_users.id = sw_articles.user_id")
	// Joins("LEFT JOIN sw_article_tags ON sw_article_tags.article_id = sw_articles.id").
	// Joins("LEFT JOIN sw_tags ON sw_tags.id = sw_article_tags.tag_id")

	// 状态 0全部1公开2封禁
	if req.ArticleCondition != 0 {
		query = query.Where("article_condition = ?", req.ArticleCondition)
	}

	// 时间
	if !req.StartTime.IsZero() && !req.EndTime.IsZero() {
		query = query.Where("published_at BETWEEN ? AND ?", req.StartTime, req.EndTime)
	} else {
		if !req.StartTime.IsZero() {
			query = query.Where("published_at >= ?", req.StartTime)
		}
		if !req.EndTime.IsZero() {
			query = query.Where("published_at <= ?", req.EndTime)
		}
	}

	// 公开文章
	query = query.Where("status = ?", "public")

	// 关键字
	if req.Keyword != "" {
		query = query.Where("title LIKE ?", "%"+req.Keyword+"%")
	}

	// 标签id
	if len(req.ArticleTags) > 0 {
		query = query.Where("sw_article_tags.tag_id IN (?)", req.ArticleTags)
	}

	// 发布人用户名
	if req.Nickname != "" {
		query = query.Where("sw_users.nickname = ?", req.Nickname)
	}

	// 浏览量 点赞量 收藏量 评论数量 热度
	if req.ViewsCount != 0 {
		query = query.Where("views_count >= ?", req.ViewsCount)
	}
	if req.LikesCount != 0 {
		query = query.Where("likes_count >= ?", req.LikesCount)
	}
	if req.CollectionsCount != 0 {
		query = query.Where("collections_count >= ?", req.CollectionsCount)
	}
	if req.CommentsCount != 0 {
		query = query.Where("comments_count >= ?", req.CommentsCount)
	}
	if req.Heat != 0 {
		query = query.Where("sw_articles.heat >= ?", req.Heat)
	}

	// 分页
	if req.Limit != 0 {
		offset := (req.Page - 1) * req.Limit
		query = query.Limit(req.Limit).Offset(offset)
	}

	// 执行查询
	if err = query.Find(&articleList).Error; err != nil {
		globals.Log.Errorf("err = %s", err)
		return // 结束函数执行
	}

	data = gin.H{"article_list": articleList}

	return data, err
}

// BanArticlesRep
// @Description: 封禁文章
// @param        db *gorm.DB
// @param        id string
// @return       error
func BanArticlesRep(db *gorm.DB, id string) error {
	// 修改文章状态
	if err := db.Model(&models.Article{}).Where("id = ?", id).Update("article_condition", 2).Error; err != nil {
		globals.Log.Errorf("err = %s", err)
		return err
	}
	return nil
}

// DeleteArticlesRep
// @Description: 删除文章
// @param        db *gorm.DB
// @param        id string
// @return       error
func DeleteArticlesRep(db *gorm.DB, id string) error {
	// 删除文章

	idInt, _ := strconv.Atoi(id)
	result := db.Delete(&models.Article{}, idInt) // 使用模型类型 + ID
	if result.Error != nil {
		return result.Error
	}

	// if err := db.Where("id = ?", id).Delete(&models.Article{}).Error; err != nil {
	//	globals.Log.Errorf("err = %s", err)
	//	return err
	// }
	return nil
}

// ArticleDetailRep
// @Description: 获取文章详情
// @param        db *gorm.DB
// @param        id string
// @return       requests.ArticleDetailRes
// @return       error
func ArticleDetailRep(db *gorm.DB, id string) (requests.ArticleDetailRes, error) {

	articleDetail := requests.ArticleDetailRes{}
	query := db.Model(&models.Article{}).
		Preload("Tags").
		Omit("like_status", "collection_status")

	query = query.Where("sw_articles.id = ?", id)

	// query = query.Joins("LEFT JOIN sw_users ON sw_users.id = sw_articles.user_id")

	query.Find(&articleDetail)

	return articleDetail, nil
}

// ArticleLikeQueryReq
// @Description: 查询文章点赞
func ArticleLikeQueryReq(db *gorm.DB, articleId int, userId int) (bool, error) {
	// var like models.Like
	// if err := db.Where("article_id = ? AND user_id = ?", articleId, userId).First(&like).Error; err != nil {
	//	return false, err
	// }
	return true, nil
}

// AboutArticleRep
// @Description: 获取相关推荐
// @param        db *gorm.DB
// @param        articleId string
// @param        userId uint
// @return       about
// @return       err
func AboutArticleRep(db *gorm.DB, articleId string, userId uint) (about []requests.AboutArticleRes, err error) {

	// articleId 和 userId 暂时搁置 未使用
	// 查询相关推荐
	query := db.Model(&models.Article{})
	query = query.Order("RAND()")
	query.Limit(4).Find(&about)

	return about, nil
}

// LikeStatusRep
// @Description: 查询用户是否点赞
// @param        db *gorm.DB
// @param        articleId string
// @param        userId uint
// @return       bool
// @return       error
func LikeStatusRep(db *gorm.DB, articleId string, userId uint) (bool, error) {
	var like models.ArticleLike
	// 查询用户是否点赞
	if err := db.Where("article_id = ? AND user_id = ?", articleId, userId).First(&like).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// 处理没有找到点赞记录的情况
			fmt.Println("没有找到点赞记录")
			return false, nil
		}
		return false, err
	}
	return true, nil
}

// CollectionStatusRep
// @Description: 查询用户是否收藏
// @param        db *gorm.DB
// @param        articleId string
// @param        userId uint
// @return       bool
// @return       error
func CollectionStatusRep(db *gorm.DB, articleId string, userId uint) (bool, error) {
	var collection models.ArticleCollection
	// 查询用户是否收藏
	if err := db.Where("article_id = ? AND user_id = ?", articleId, userId).First(&collection).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// 处理没有找到收藏记录的情况
			fmt.Println("没有找到收藏记录")
			return false, nil
		}
		return false, err
	}
	return true, nil
}

// UpdateLikeRep
// @Description: 更新点赞
// @param        db *gorm.DB
// @param        articleId string
// @param        userId uint
// @return       error
func UpdateLikeRep(db *gorm.DB, req requests.ArticleLikeReq, userId uint) (err error) {

	var like models.ArticleLike

	if req.LikeStatus == true {
		if err = db.Where("article_id = ? AND user_id = ?", req.ArticleId, userId).First(&like).Error; err != nil {

			// 如果没有找到点赞记录，创建点赞记录
			like = models.ArticleLike{
				ArticleID: req.ArticleId,
				UserID:    userId,
			}
			if err = db.Create(&like).Error; err != nil {
				return err
			}

			// 创建点赞记录后 通知文章作者
			var article models.Article
			if err = db.Where("id = ?", req.ArticleId).First(&article).Error; err != nil {
				return err
			}
			// 通知文章作者
			authorId := strconv.Itoa(int(article.UserID))
			internalUtils.MessagePush("like", authorId)
		}

	} else if req.LikeStatus == false {
		// 查询用户是否点赞 如果点赞了 删除点赞记录
		if err = db.Where("article_id = ? AND user_id = ?", req.ArticleId, userId).First(&like).Error; err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			return err
		}
		// 删除点赞记录
		if err = db.Delete(&like).Error; err != nil {
			return err
		}
	}

	return nil
}

// UpdateCollectionRep
// @Description: 更新收藏
// @param        db *gorm.DB
// @param        req requests.ArticleCollectionReq
// @param        userId uint
// @return       error
func UpdateCollectionRep(db *gorm.DB, req requests.ArticleCollectionReq, userId uint) (err error) {

	var collection models.ArticleCollection

	if req.CollectionStatus == true {
		if err = db.Where("article_id = ? AND user_id = ?", req.ArticleId, userId).First(&collection).Error; err != nil {
			// 如果没有找到收藏记录，创建收藏记录
			collection = models.ArticleCollection{
				ArticleID: req.ArticleId,
				UserID:    userId,
			}
			if err = db.Create(&collection).Error; err != nil {
				return err
			}

			// 创建收藏记录后 通知文章作者
			var article models.Article
			if err = db.Where("id = ?", req.ArticleId).First(&article).Error; err != nil {
				//globals.Log.Errorf("err = %s", err)
				return err
			}
			// 通知文章作者
			authorId := strconv.Itoa(int(article.UserID))
			internalUtils.MessagePush("collection", authorId)

		}
	} else if req.CollectionStatus == false {
		//// 查询用户是否收藏 如果收藏了 删除收藏记录
		//if err = db.Where("article_id = ? AND user_id = ?", req.ArticleId, userId).First(&collection).Error; err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		//	return err
		//}
		// 删除收藏记录
		if err = db.Where("article_id = ? AND user_id = ?", req.ArticleId, userId).Delete(&collection).Error; err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			return err
		}

	}

	return nil
}