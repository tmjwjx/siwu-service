package logics

import (
	"forum/internal/article/repositories"
	"forum/internal/article/requests"
	"forum/pkg/globals"
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

// BanArticlesLocal
// @Description: 封禁文章
// @param        db *gorm.DB
// @param        id string
// @return       error
func BanArticlesLocal(db *gorm.DB, id string) error {
	err := repositories.BanArticlesRep(db, id)
	if err != nil {
		globals.Log.Errorf("封禁文章失败 err = %s", err)
		return err
	}
	return nil
}

// DeleteArticlesLocal
// @Description: 删除文章
// @param        db *gorm.DB
// @param        id string
// @return       error
func DeleteArticlesLocal(db *gorm.DB, id string) error {
	err := repositories.DeleteArticlesRep(db, id)
	if err != nil {
		globals.Log.Errorf("封禁文章失败 err = %s", err)
		return err
	}
	return nil
}