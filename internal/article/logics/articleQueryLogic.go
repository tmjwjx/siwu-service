package logics

import (
	"forum/internal/article/repositories"
	"forum/internal/article/requests"
	"forum/pkg/globals"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// ArticleSearchLogic 搜索文章
func ArticleSearchLogic(db *gorm.DB, req *requests.ArticleSearchReq) ([]requests.SearchArticleListRes, error) {

	articles, err := repositories.SearchArticlesRep(db, req)
	if err != nil {
		return nil, err
	}
	return articles, nil
}

// ArticleDetailLogic
// @Description: 获取文章详情
// @param        db *gorm.DB
// @param        id string
// @return       data
// @return       err
func ArticleDetailLogic(db *gorm.DB, articleId string, userId uint) (data interface{}, err error) {

	article, err := repositories.ArticleDetailRep(db, articleId)
	if err != nil {
		globals.Log.Errorf("err = %s", err)
		return nil, err
	}

	about, err := repositories.AboutArticleRep(db, articleId, userId)
	if err != nil {
		globals.Log.Errorf("err = %s", err)
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
	//if err != nil {
	//	return nil, err
	//}
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