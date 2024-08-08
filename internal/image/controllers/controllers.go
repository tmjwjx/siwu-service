package controllers

import (
	"fmt"
	"forum/internal/image/logics"
	"forum/internal/models"
	"forum/pkg/globals"
	"github.com/gin-gonic/gin"
)

// UploadHandlerControllers 将前端传过来的图片文件存到文件系统中
// kind 用来判断图片是用户的还是文章的, ID 是用户或文章的ID

func UploadImagesControllers(c *gin.Context, home string, homeID uint) (error, int) { // (error, int): error表示错误,int表示响应的状态码
	// 将文件系统中的目录映射到 URL 路径
	globals.Router.Static("/images", "./static/images")

	// 解析multipart/form-data
	if err := c.Request.ParseMultipartForm(32 << 20); err != nil { // 设置最大大小为 32MB
		return err, 400
	}
	// 具体逻辑实现
	err := logics.UploadHandlerLogic(c, home, homeID)
	if err != nil {
		return err, 500
	}
	return nil, 200
}

// GetImagesControllers 从数据库中将图片路径取出
func GetImagesControllers(home string, homeID uint) (*[]models.Attachment, error) {
	images, err := logics.GetImagesLogic(home, homeID)
	if err != nil {
		return nil, fmt.Errorf("GetImagesControllers -> %s", err)
	}
	return images, nil
}

// GetAdvertisementImageCtrl 专门用于取数据库中的广告图片
func GetAdvertisementImageCtrl(home string, status int) (*[]models.Advertisement, error) {
	images, err := logics.GetAdvertisementImageLogic(home, status)
	if err != nil {
		return nil, fmt.Errorf("GetAdvertisementImageCtrl -> %s", err)
	}
	return images, nil
}
