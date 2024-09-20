package logics

import (
	"forum/internal/article/repositories"
	"forum/internal/article/requests"
	"forum/internal/models"
	"forum/pkg/globals"
	"gorm.io/gorm"
)

// ArticleSearchLogic 搜索文章
func ArticleSearchLogic(db *gorm.DB, req *requests.ReqSearch) ([]models.Article, error) {

	articles, err := repositories.SearchArticlesRep(db, req)
	if err != nil {
		return nil, err
	}
	return articles, nil
}

// ArticleDetailLogic
// @Description: 文章详情
// @param        db *gorm.DB
// @param        id string
// @return       data
// @return       err
func ArticleDetailLogic(db *gorm.DB, id string) (data interface{}, err error) {

	data, err = repositories.ArticleDetailRep(db, id)
	if err != nil {
		globals.Log.Errorf("err = %s", err)
		return nil, err
	}
	return nil, nil

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