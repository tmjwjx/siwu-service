package logics

import (
	"forum/internal/article/repositories"
	"forum/internal/article/requests"
	"gorm.io/gorm"
)

// GetArticleList
// @Description: 检索获取已经发布的文章列表
// @param        db *gorm.DB
// @return       data
// @return       err
func GetArticleList(db *gorm.DB, req *requests.ArticleListReq) (data interface{}, err error) {

	data, err = repositories.SearchArticlesListRep(db, req)
	//if err != nil {
	//	return nil, err
	//}
	return data, nil
}