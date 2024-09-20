package logics

import (
	"forum/internal/article/repositories"
	"forum/pkg/globals"
	"gorm.io/gorm"
)

// ArticleBanLocal
// @Description: 封禁文章
// @param        db *gorm.DB
// @param        id string
// @return       error
func ArticleBanLocal(db *gorm.DB, id string) error {
	err := repositories.BanArticlesRep(db, id)
	if err != nil {
		globals.Log.Errorf("封禁文章失败 err = %s", err)
		return err
	}
	return nil
}

// ArticleDeleteLocal
// @Description: 删除文章
// @param        db *gorm.DB
// @param        id string
// @return       error
func ArticleDeleteLocal(db *gorm.DB, id string) error {
	err := repositories.DeleteArticlesRep(db, id)
	if err != nil {
		globals.Log.Errorf("封禁文章失败 err = %s", err)
		return err
	}
	return nil
}