package repositories

import (
	"context"
	"fmt"
	"forum/internal/article/requests"
	"forum/internal/internalPkg/internalUtils"
	"forum/internal/internalPkg/redisUtils"
	"forum/internal/models"
	"forum/pkg/globals"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"gorm.io/gorm"
	"strconv"
	"time"
)

// GetTodayViewsRep
// @Description: 获取今日浏览量
// @param        db *gorm.DB
// @return       todayViews
// @return       err
// @Author tianjiajie 2025-01-16 15:10:42
func GetTodayViewsRep(db *gorm.DB) (todayViews int64, err error) {

	rdb := globals.RDB
	ctx := context.Background()
	// 查询今日浏览量
	nowTime := string(time.Now().Format("2006-01-02"))
	todayViewsStr, err := redisUtils.HGet(rdb, ctx, "todayViews", nowTime)
	if err != nil {
		return 0, err
	}
	todayViews, _ = strconv.ParseInt(todayViewsStr, 10, 64)

	return todayViews, nil
}

// GetTodayCommentsRep
// @Description: 获取今日评论数量
// @param        db *gorm.DB
// @return       todayComments
// @return       err
// @Author tianjiajie 2025-01-16 15:00:59
func GetTodayCommentsRep(db *gorm.DB) (todayComments int64, err error) {

	// 查询今日评论数量
	if err = db.Model(&models.ArticleComment{}).
		Where("DATE(created_at) = ?", time.Now().Format("2006-01-02")).
		Count(&todayComments).Error; err != nil {
		return 0, err
	}
	return todayComments, nil
}

// GetNewAddArticleRep
// @Description: 获取今日新增文章数量
// @param        db *gorm.DB
// @return       newArticle
// @return       err
// @Author tianjiajie 2025-01-16 14:58:07
func GetNewAddArticleRep(db *gorm.DB) (newArticle int64, err error) {

	// 查询今日文章数量
	if err = db.Model(&models.Article{}).
		Where("DATE(created_at) = ?", time.Now().Format("2006-01-02")).
		Count(&newArticle).Error; err != nil {
		return 0, err
	}
	return newArticle, nil

}

// GetArticleCountRep
// @Description: 获取文章数量
// @param        db *gorm.DB
// @return       articleCount
// @return       err
// @Author tianjiajie 2025-01-16 14:49:16
func GetArticleCountRep(db *gorm.DB) (articleCount int64, err error) {
	// 查询文章数量
	if err = db.Model(&models.Article{}).Count(&articleCount).Error; err != nil {
		return 0, err
	}
	return articleCount, nil
}

// GetHotTagsRep
// @Description: 查询前 n 的热门标签
// @param        db *gorm.DB
// @param        n int
// @return       tags
// @return       err
// @Author tianjiajie 2025-01-15 20:36:55
func GetHotTagsRep(db *gorm.DB, n int) (data interface{}, err error) {
	// 查询热门标签
	var hotTags []struct {
		ID    uint   `json:"id"`    // 标签ID
		Name  string `json:"name"`  // 标签名称
		Count int    `json:"count"` // 文章数量
	}
	err = db.Model(&models.Tag{}).
		Select("sw_tags.id, sw_tags.name, COUNT(sw_article_tags.article_id) AS count").
		Joins("LEFT JOIN sw_article_tags ON sw_tags.id = sw_article_tags.tag_id").
		Group("sw_tags.id").
		Order("count DESC").
		Limit(n).
		Scan(&hotTags).Error // 执行查询并将结果存入 hotTags 切片
	if err != nil {
		return nil, err
	}

	// 所有tag下的文章总数
	var total int64
	err = db.Table("sw_article_tags").
		Select("COUNT(DISTINCT article_id, tag_id) AS total_count"). // 计算不同的 article_id 和 tag_id 组合
		Scan(&total).Error
	if err != nil {
		return nil, err
	}

	var hotTagsRes []struct {
		Type  string `json:"type"`  // 标签名称
		Value string `json:"value"` // 占比
	}

	other := 100.00
	// 将查询结果存入 tags 切片
	for _, tag := range hotTags {
		// 计算占比
		temp := float64(tag.Count) / float64(total) * 100
		other -= temp
		value := fmt.Sprintf("%.2f", temp)
		hotTagsRes = append(hotTagsRes, struct {
			Type  string `json:"type"`
			Value string `json:"value"`
		}{tag.Name, value})
	}

	value := fmt.Sprintf("%.2f", other)
	hotTagsRes = append(hotTagsRes, struct {
		Type  string `json:"type"`
		Value string `json:"value"`
	}{"其他", value})

	data = hotTagsRes

	return
}

// GetHotArticleRep
// @Description: 查询前 n 篇热门文章数据
// @param        db *gorm.DB
// @param        limit int
// @return       articleList
// @return       err
// @Author tianjiajie 2025-01-15 14:54:37
func GetHotArticleRep(db *gorm.DB, n int) (articleList []requests.HotArticleRes, err error) {

	var articles []struct {
		ID              uint    `json:"id"`                // 文章ID
		Title           string  `json:"title"`             // 文章标题
		LikesCount      int     `json:"likes_count"`       // 点赞数
		Increase        float64 `json:"increase"`          // 点赞涨幅
		TodayLikesCount int     `json:"today_likes_count"` // 今日点赞数
	}
	// 查询前n篇文章数据
	err = db.Table("sw_articles a").
		Select("a.id, a.title, a.likes_count, " +
			"IF(a.likes_count > 0, IFNULL(COUNT(b.article_id), 0) / a.likes_count, 0) AS increase," +
			"IFNULL(COUNT(b.article_id), 0) AS today_likes_count").
		Joins("LEFT JOIN sw_article_likes b ON a.id = b.article_id AND b.deleted_at IS NULL AND DATE(b.updated_at) = CURDATE()").
		Group("a.id").
		Order("a.heat DESC").
		Limit(n).
		//Debug().              // 添加这一行
		Scan(&articles).Error // 执行查询并将结果存入 articles 切片
	if err != nil {
		return nil, err
	}
	// 将查询结果存入 articleList 切片
	for _, article := range articles {
		// 处理点赞量格式
		var likesCount string
		if article.LikesCount >= 10000 {
			likesCount = fmt.Sprintf("%.2fw", float64(article.LikesCount)/10000)
		} else if article.LikesCount >= 1000 {
			likesCount = fmt.Sprintf("%.2fk", float64(article.LikesCount)/1000)
		} else {
			likesCount = strconv.Itoa(article.LikesCount)
		}

		// 处理点赞涨幅格式
		increase := fmt.Sprintf("%.2f", article.Increase*100)

		// 将处理后的数据存入 articleList
		articleList = append(articleList, requests.HotArticleRes{
			ID:         article.ID,
			Title:      article.Title,
			LikesCount: likesCount,
			Increase:   increase,
		})
	}
	return
}

// GetArticleNumRep
// @Description: 查询 某天发布的文章
// @param        db *gorm.DB
// @param        now time.Time
// @param        ago time.Time
// @return       articleSum
// @return       err
// @Author tianjiajie 2025-01-15 09:48:50
func GetArticleNumRep(db *gorm.DB, date time.Time) (articleCount int64, err error) {
	// 查询指定日期发布的文章数量
	if err = db.Model(&models.Article{}).
		Where("DATE(published_at) = ?", date.Format("2006-01-02")).
		Count(&articleCount).Error; err != nil {
		return 0, err
	}

	return articleCount, nil
}

// InsertArticlesRep
// @Description: 新建文章
// @param        db *gorm.DB
// @param        req requests.ReqPublish
// @return       int
// @return       error
func InsertArticlesRep(db *gorm.DB, req requests.ReqPublish, userId uint) (id int, err error) {

	newArticle := models.Article{
		// UserID:     req.UserId,
		// Title:      req.Title,
		// Status:     req.Status,
		// CategoryID: req.CategoryID,
		// Summary:    req.Summary,
		// Content:    req.Content,
		// ImageUrl:   req.ImageUrl,
	}
	fmt.Println("我进来了")

	// 判断当前用户是否是文章作者
	if userId != req.UserId {
		// 如果当前用户不是文章作者
		return 0, errors.New("当前用户不是文章作者")
	}

	// 设置文章ID
	if req.ArticleId != 0 {
		newArticle.ID = uint(req.ArticleId)

		//// 获取数据库文章记录
		//db.Find(&newArticle)
		//fmt.Println("我进来了")

		if err = db.First(&newArticle, req.ArticleId).Error; err != nil {
			return 0, fmt.Errorf("未找到指定文章: %w", err)
		}

		// 更新文章的标签关联
		// 手动清空当前文章的标签关联，确保旧数据被清除
		if err = db.Model(&newArticle).Association("Tags").Clear(); err != nil {
			return 0, err // 如果清除失败，返回错误
		}
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

	fmt.Println("newArticle", newArticle)

	// 设置标签
	// 查找传递过来的所有标签
	var tags []models.Tag
	if err := db.Where("id IN ?", req.Tags).Find(&tags).Error; err != nil {
		return 0, err // 如果标签不存在，返回错误
	}

	newArticle.Tags = tags

	// 设置发布时间
	if req.Status == "public" && newArticle.PublishedAt == nil {
		now := time.Now()
		newArticle.PublishedAt = &now
	}
	result := db.Save(&newArticle)
	if result.Error != nil {
		globals.Log.Fatal("创建文章失败:", result.Error)
		return 0, result.Error
	}
	id = int(newArticle.ID)

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

	// 查询文章列表
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

	// 查询公开文章
	query = query.Where("sw_articles.status = ?", "public")

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
	var articleList []requests.SearchArticleListRes
	var totalCount int64

	// 构建基础查询
	query := db.Model(&models.Article{}).Preload("Tags").
		Joins("LEFT JOIN sw_users ON sw_users.id = sw_articles.user_id").
		Joins("LEFT JOIN sw_article_tags ON sw_article_tags.article_id = sw_articles.id")

	// 状态 0公开1全部2封禁
	if req.ArticleCondition != 1 {
		query = query.Where("article_condition = ?", req.ArticleCondition)
	}

	// 时间范围过滤
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
	query = query.Where("sw_articles.status = ?", "public")

	// 公开/封禁文章过滤
	if req.ArticleCondition != 1 {
		query = query.Where("article_condition = ?", req.ArticleCondition)
	}

	// 关键字过滤
	if req.Keyword != "" {
		query = query.Where("title LIKE ?", "%"+req.Keyword+"%")
	}

	// 标签id过滤
	if len(req.ArticleTags) > 0 {
		query = query.Where("sw_article_tags.tag_id IN (?)", req.ArticleTags)
	}

	// 发布人用户名过滤
	if req.Nickname != "" {
		query = query.Where("sw_users.nickname = ?", req.Nickname)
	}

	// 浏览量、点赞量等过滤
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

	// 获取总数（不使用 DISTINCT）
	err = query.Group("sw_articles.id").
		Count(&totalCount).Error
	if err != nil {
		globals.Log.Errorf("Error counting articles: %v", err)
		return nil, err
	}

	// 选择排序方式
	condition := internalUtils.ArticlesOrder(req.Kind) // 选择排序方式  0热度 1时间
	query = query.Order(condition)

	// 分页
	if req.Limit != 0 {
		offset := (req.Page - 1) * req.Limit
		query = query.Limit(req.Limit).Offset(offset)
	}

	// 执行查询，获取文章列表（使用 DISTINCT）
	query = query.Select("DISTINCT sw_articles.*, sw_users.nickname")
	if err = query.Find(&articleList).Error; err != nil {
		globals.Log.Errorf("err = %s", err)
		return nil, err
	}

	// 返回结果
	data = gin.H{"article_list": articleList, "total": totalCount}
	return data, nil
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
		Select("sw_articles.*, sw_users.nickname").
		Joins("LEFT JOIN sw_users ON sw_users.id = sw_articles.user_id").
		Preload("Tags").
		Omit("like_status", "collection_status").Debug()

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

// AddArticleClickRep
// @Description: 增加文章点击量
// @param        db *gorm.DB
// @param        articleId string
// @return       error
// @Author tianjiajie 2025-01-16 08:46:29
func AddArticleClickRep(db *gorm.DB, articleId string) error {
	// 增加文章点击量
	if err := db.Model(&models.Article{}).Where("id = ?", articleId).Update("views_count", gorm.Expr("views_count + ?", 1)).Error; err != nil {
		return err
	}
	return nil
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

			if errors.Is(err, gorm.ErrRecordNotFound) {
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
			} else {
				return err
			}
		}

	} else if req.LikeStatus == false {
		if err = db.Where("article_id = ? AND user_id = ?", req.ArticleId, userId).First(&like).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return nil // 用户本来就没有点赞，直接返回
			}
			return err // 其他错误，返回
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

			if errors.Is(err, gorm.ErrRecordNotFound) {
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
			} else {
				return err
			}

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

// GetArticlesByTagRep
// @Description: 获取标签下的文章
// @param        db *gorm.DB
// @param        tagId int
// @return       articles
// @return       err
// @Author tianjiajie 2024-10-15 17:07:21
func GetArticlesByTagRep(db *gorm.DB, req *requests.GetArticleByTagReq) (articles []requests.SearchArticleListRes, err error) {

	globals.Log.Infof("req = %v", req)
	condition := internalUtils.ArticlesOrder(req.Kind) // 选择排序方式  0热度 1时间

	// 查询标签下的文章
	query := db.Model(&models.Article{}).Preload("Tags").
		Select("sw_articles.*, sw_users.nickname").
		Joins("LEFT JOIN sw_users ON sw_users.id = sw_articles.user_id").
		Joins("LEFT JOIN sw_article_tags ON sw_article_tags.article_id = sw_articles.id").
		Joins("LEFT JOIN sw_tags ON sw_tags.id = sw_article_tags.tag_id").
		Where("sw_article_tags.tag_id = ?", req.Id).
		Order(condition)

	// 分页
	if req.Limit != 0 {
		offset := (req.Page - 1) * req.Limit
		query = query.Limit(req.Limit).Offset(offset)
	}

	// 只获取公开文章
	query = query.Where("sw_articles.status = ?", "public").Where("sw_articles.article_condition = ?", 0)

	// 执行查询
	if err = query.Find(&articles).Error; err != nil {
		globals.Log.Errorf("err = %s", err)
		return // 结束函数执行
	}

	return articles, nil
}

// GetUserArticleOrCollectionRep
// @Description: 获取用户文章或收藏列表
// @param        *gorm.DB *gorm.DB
// @param        *requests.UserArticleOrCollectionReq *requests.UserArticleOrCollectionReq
// @param        int int
// @return       articles
// @return       err
// @Author tianjiajie 2024-10-18 16:38:47
func GetUserArticleOrCollectionRep(db *gorm.DB, req *requests.UserArticleOrCollectionReq, id int) (articles []requests.SearchArticleListRes, total int64, err error) {

	condition := internalUtils.ArticlesOrder(1) // 选择排序方式  0热度 1时间

	query := db.Model(&models.Article{}).Preload("Tags").
		Select("sw_articles.*, sw_users.nickname").
		Joins("LEFT JOIN sw_users ON sw_users.id = sw_articles.user_id")
	// 获取未封禁的文章
	//Where("sw_articles.article_condition = ?", 0)

	// 判断是发布的文章还是收藏的文章
	switch req.Type {
	case "收藏":
		query = query.Joins("JOIN article_collections ON articles.id = article_collections.article_id").
			Where("article_collections.user_id = ?", req.Id)
	case "文章":
		query = query.Where("sw_articles.user_id = ?", req.Id)
	}

	// 按关键词搜索
	if req.Keyword != "" {
		query = query.Where("title LIKE ? OR summary LIKE ?", "%"+req.Keyword+"%", "%"+req.Keyword+"%")
	}

	// 选择排序方式 时间or热度
	query = query.Order(condition)

	// 判断是否分页
	if req.Limit != 0 {
		offset := (req.Page - 1) * req.Limit // 计算当前页的偏移量，用于分页
		query = query.Limit(req.Limit).Offset(offset)
	}

	// 只获取公开文章
	if req.Id != id {
		query = query.Where("sw_articles.status = ?", "public").
			Where("sw_articles.article_condition = ?", 0)
	}

	// 查询数量
	err = query.Count(&total).Error
	if err != nil {
		globals.Log.Errorf("Error counting articles: %v", err)
		return nil, 0, err
	}

	// 执行查询
	if err = query.Find(&articles).Error; err != nil {
		globals.Log.Errorf("err = %s", err)
		return // 结束函数执行
	}

	return articles, total, nil
}
