package logics

import (
	"errors"
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

// UserArticleDeleteLocal
// @Description: 用户 删除文章
// @param        db *gorm.DB
// @param        id string
// @return       error
// @Author tianjiajie 2025-01-21 09:06:00
func UserArticleDeleteLocal(db *gorm.DB, articleId string, userId uint) error {
	authorId, err := repositories.QueryArticleAuthor(db, articleId)
	if err != nil {
		globals.Log.Errorf("查询文章作者失败 err = %s", err)
		return err
	}

	// 判断是否是作者
	if authorId != userId {
		globals.Log.Errorf("用户不是文章作者")
		err = errors.New("用户不是文章作者")
	} else {
		err = repositories.DeleteArticlesRep(db, articleId)
		if err != nil {
			globals.Log.Errorf("封禁文章失败 err = %s", err)
		}
	}
	return err
}
