package logics

import (
	"errors"
	"forum/internal/article/repositories"
	"forum/internal/article/requests"
	"forum/pkg/globals"
	"gorm.io/gorm"
	"strconv"
)

// ArticleBanLocal
// @Description: 封禁文章
// @param        db *gorm.DB
// @param        id string
// @return       error
func ArticleBanLocal(db *gorm.DB, idList requests.ArticleOperationListReq) error {
	err := repositories.BanArticlesRep(db, idList)
	if err != nil {
		globals.Log.Errorf("封禁文章失败 err = %s", err)
		return err
	}
	return nil
}

// ArticleUnblockLocal
// @Description: 解封文章
// @param        db *gorm.DB
// @param        idList requests.ArticleOperationListReq
// @return       error
// @Author tianjiajie 2025-01-21 11:30:50
func ArticleUnblockLocal(db *gorm.DB, idList requests.ArticleOperationListReq) error {
	err := repositories.UnblockArticlesRep(db, idList)
	if err != nil {
		globals.Log.Errorf("解封文章失败 err = %s", err)
		return err
	}
	return nil
}

// ArticleDeleteLocal
// @Description: 删除文章
// @param        db *gorm.DB
// @param        id string
// @return       error
func ArticleDeleteLocal(db *gorm.DB, idList requests.ArticleOperationListReq) error {
	for _, id := range idList.IdList {
		idStr := strconv.Itoa(int(id))
		err := repositories.DeleteArticlesRep(db, idStr)
		if err != nil {
			globals.Log.Errorf("封禁文章失败 err = %s", err)
			return err
		}
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
