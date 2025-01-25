package logics

import (
	"forum/internal/article/repositories"
	"forum/internal/article/requests"
	"gorm.io/gorm"
)

// LikeArticleLogic
// @Description: 文章点赞逻辑
// @param        db *gorm.DB
// @param        articleId uint
// @param        userId uint
// @return       data
// @return       err
func LikeArticleLogic(db *gorm.DB, req requests.ArticleLikeReq, userId uint) (err error) {
	err = repositories.UpdateLikeRep(db, req, userId)
	if err != nil {
		return err
	}

	return nil
}

// CollectionArticleLogic
// @Description: 文章收藏逻辑
// @param        db *gorm.DB
// @param        req requests.ArticleCollectionReq
// @param        userId uint
// @return       err
func CollectionArticleLogic(db *gorm.DB, req requests.ArticleCollectionReq, userId uint) (err error) {
	err = repositories.UpdateCollectionRep(db, req, userId)
	if err != nil {
		return err
	}
	return nil
}
