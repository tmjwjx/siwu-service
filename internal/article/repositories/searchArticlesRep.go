package repositories

import (
	"fmt"
	"forum/internal/article/requests"
	"forum/internal/models"
	"forum/pkg/globals"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// SearchArticlesRep 搜索文章
func SearchArticlesRep(db *gorm.DB, req *requests.ReqSearch) (articles []models.Article, err error) {
	var condition string
	if req.Kind == 0 { // 0 代表按照热度排序
		condition = "heat DESC"
	} else if req.Kind == 1 { // 1 代表按照发布时间排序
		condition = "published_at DESC"
	}

	query := db.Model(&models.Article{})

	if req.UserId != 0 { // 按用户ID筛选
		query = query.Joins("JOIN user_follows ON articles.user_id = user_follows.followed_id").
			Where("user_follows.follower_id = ?", req.UserId)
		condition = "published_at DESC" // 按照 关注 查询只能按照 时间 排序
	} else if req.Query != "" { // 按关键词搜索
		query = query.Where("title LIKE ? OR summary LIKE ?", "%"+req.Query+"%", "%"+req.Query+"%")
	} else if req.Tag != "" { //按标签筛选
		var tag models.Tag
		// 根据标签名称查找 tag_id
		if err = db.Where("name = ?", req.Tag).First(&tag).Error; err != nil {
			globals.Log.Errorf("err = %s", err)
			return nil, err // 返回 nil 和错误信息
		}
		fmt.Println(tag)
		query = query.Joins("JOIN article_tags ON articles.id = article_tags.article_id").
			Where("article_tags.tag_id = ?", tag.ID)
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
	//var articleList []requests.ArcList
	var articleList []requests.SearchArticleListRes

	query := db.Model(&models.Article{}).Preload("Tags").
		Joins("LEFT JOIN sw_users ON sw_users.id = sw_articles.user_id").
		Joins("LEFT JOIN sw_article_tags ON sw_article_tags.article_id = sw_articles.id").
		Select("DISTINCT sw_articles.*, sw_users.nickname")
	//query.Select(
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
	//Joins("LEFT JOIN sw_users ON sw_users.id = sw_articles.user_id")
	//Joins("LEFT JOIN sw_article_tags ON sw_article_tags.article_id = sw_articles.id").
	//Joins("LEFT JOIN sw_tags ON sw_tags.id = sw_article_tags.tag_id")

	//状态
	query = query.Where("article_condition = ?", req.ArticleCondition)

	//时间
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
	if req.Title != "" {
		query = query.Where("title LIKE ?", "%"+req.Title+"%")
	}

	//标签id
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

	//articleMap := make(map[uint]requests.SearchArticleListRes)
	//for _, article := range articleList {
	//	_, exists := articleMap[article.ID]
	//	if !exists {
	//		var tags []requests.TagsRes
	//
	//		db.Model(&models.ArticleTag{}).
	//			Select("sw_tags.id AS tag_id, "+
	//				"sw_tags.name AS tag_name").
	//			Joins("LEFT JOIN sw_tags ON sw_tags.id = sw_article_tags.tag_id").
	//			Where("article_id = ?", article.ID).Find(&tags)
	//
	//		article.Tags = append(article.Tags, tags...)
	//
	//		articleMap[article.ID] = article
	//	}
	//}

	data = gin.H{"article_list": articleList}

	return data, err
}