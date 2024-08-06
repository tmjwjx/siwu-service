package repositorys

import (
	"forum/internal/image/requests"
	"forum/internal/models"
	"forum/pkg/globals"
	"github.com/gin-gonic/gin"
)

// InsertFile 将图片文件的路径相关的信息存入数据库中
func InsertFile(c *gin.Context, attachment *requests.Attachment, kind int, ID uint) error {
	// 向数据库中存入文件数据
	result := globals.DB.Create(attachment)
	if result.Error != nil {
		return result.Error
	}

	// 判断图片是用户图片，还是文章图片
	if kind == 1 {
		err := InsertPathToUser(c, ID, attachment.Path)
		if err != nil {
			//e := response.NewAppErr(globals.StatusInternalServerError, err, nil)
			//response.Failed(c, e, 5000)
			return err
		}
	} else if kind == 2 {
		err := InsertPathToArticle(c, ID, attachment.Path)
		if err != nil {
			//e := response.NewAppErr(globals.StatusInternalServerError, err, nil)
			//response.Failed(c, e, 5000)
			return err
		}
	}
	return nil
}

// InsertPathToUser 将图片的url路径存进User 表中
func InsertPathToUser(c *gin.Context, userID uint, path string) error {
	var user models.User
	if err := globals.DB.Take(&user, "id = ?", userID).Error; err != nil {

		return err
	}
	if err := globals.DB.Model(&user).Update("path", path).Error; err != nil {

		return err
	}

	return nil
}

// InsertPathToArticle 将图片的url路径存进Article 表中
func InsertPathToArticle(c *gin.Context, articleID uint, path string) error {
	var article models.Article
	if err := globals.DB.Take(&article, "id = ?", articleID).Error; err != nil {
		return err
	}
	if err := globals.DB.Model(&article).Update("path", path).Error; err != nil {
		return err
	}

	return nil
}
