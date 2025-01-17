package logics

import (
	"context"
	"forum/internal/article/repositories"
	"forum/internal/article/requests"
	"forum/internal/internalPkg/redisUtils"
	"forum/pkg/globals"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"time"
)

// ArticleSearchLogic 搜索文章
func ArticleSearchLogic(db *gorm.DB, req *requests.ArticleSearchReq) (data interface{}, err error) {

	articles, err := repositories.SearchArticlesRep(db, req)
	if err != nil {
		return nil, err
	}

	data = gin.H{"selectedList": articles}

	return data, nil
}

// ArticleDetailLogic
// @Description: 获取文章详情
// @param        db *gorm.DB
// @param        id string
// @return       data
// @return       err
func ArticleDetailLogic(db *gorm.DB, articleId string, userId uint) (data interface{}, err error) {

	// 查询文章详情
	article, err := repositories.ArticleDetailRep(db, articleId)
	if err != nil {
		globals.Log.Errorf("err = %s", err)
		return nil, err
	}
	// 查询用户是否点赞
	article.LikeStatus, err = repositories.LikeStatusRep(db, articleId, userId)
	// 查询用户是否收藏
	article.CollectionStatus, err = repositories.CollectionStatusRep(db, articleId, userId)

	// 查询相关文章
	about, err := repositories.AboutArticleRep(db, articleId, userId)
	if err != nil {
		globals.Log.Errorf("err = %s", err)
		return nil, err
	}

	// 增加点击量
	err = repositories.AddArticleClickRep(db, articleId)
	if err != nil {
		globals.Log.Errorf("err = %s", err)
		return nil, err
	}

	// 增加当天的访问量
	rdb := globals.RDB
	ctx := context.Background()
	nowTime := string(time.Now().Format("2006-01-02"))
	_, err = redisUtils.IncrementHash(rdb, ctx, "todayViews", nowTime, 1)
	if err != nil {
		return nil, err
	}

	data = gin.H{"article": article, "about": about}
	return data, nil

}

// ArticleListLogic
// @Description: 检索获取已经发布的文章列表
// @param        db *gorm.DB
// @return       data
// @return       err
func ArticleListLogic(db *gorm.DB, req *requests.ArticleListReq) (data interface{}, err error) {

	data, err = repositories.SearchArticlesListRep(db, req)
	if err != nil {
		return nil, err
	}
	return data, nil
}

// ArticleEditLogic
// @Description: 编辑文章的界面所需的数据
// @param        db *gorm.DB
// @return       data
// @return       err
func ArticleEditLogic(db *gorm.DB) (data interface{}, err error) {
	t := requests.TagsReq{}
	c := requests.CategorysReq{}

	// 查询标签
	tags, err := repositories.QueryTag(db)
	if err != nil {
		return nil, err
	}
	// 将标签转换为前端需要的数据格式
	for _, tag := range tags {
		t = append(t, struct {
			Value uint   `json:"value"`
			Label string `json:"label"`
		}{Value: tag.ID, Label: tag.Name})
	}

	// 查询类目
	categories, err := repositories.QueryCategory(db)
	if err != nil {
		return nil, err
	}
	// 将类目转换为前端需要的数据格式
	for _, category := range categories {
		c = append(c, struct {
			Value uint   `json:"value"`
			Label string `json:"label"`
		}{Value: category.ID, Label: category.Name})
	}

	data = gin.H{"tags": t, "categories": c}

	return data, nil
}

// GetArticlesByTagLogic
// @Description: 获取标签下的文章
// @Author tianjiajie 2024-10-15 16:59:17
func GetArticlesByTagLogic(db *gorm.DB, req *requests.GetArticleByTagReq) (data interface{}, err error) {

	articleList, err := repositories.GetArticlesByTagRep(db, req)
	if err != nil {
		return nil, err
	}
	data = gin.H{"article_list": articleList}
	return data, nil

}

// GetUserArticleOrCollectionLogic
// @Description: 获取用户文章或收藏列表
// @Author tianjiajie 2024-10-18 15:37:28
func GetUserArticleOrCollectionLogic(db *gorm.DB, req *requests.UserArticleOrCollectionReq, id int) (data interface{}, err error) {

	articleList, err := repositories.GetUserArticleOrCollectionRep(db, req, id)
	if err != nil {
		return nil, err
	}
	// 将文章状态转换为中文
	for i := 0; i < len(articleList); i++ {
		switch articleList[i].Status {
		case "public":
			articleList[i].Status = "公开"
		case "private":
			articleList[i].Status = "私有"
		case "draft":
			articleList[i].Status = "草稿"
		}

	}
	data = gin.H{"dataList": articleList}
	return data, nil
}
