package logics

import (
	"forum/internal/article/repositories"
	"forum/internal/article/requests"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// ArticleCreateLogic 发布文章
func ArticleCreateLogic(db *gorm.DB, req requests.ReqPublish, userId uint) (data interface{}, err error) {
	id, err := repositories.InsertArticlesRep(db, req, userId)
	if err != nil {
		return nil, err
	}
	data = gin.H{"id": id}
	return data, nil
}
