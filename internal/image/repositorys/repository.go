package repositorys

import (
	"fmt"
	"forum/internal/image/requests"
	"forum/internal/models"
	"forum/pkg/globals"
)

// InsertFile 将图片文件的路径相关的信息存入数据库中
func InsertFile(attachment *requests.Attachment) error {
	// 向数据库中存入文件数据
	result := globals.DB.Create(attachment)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

// GetImages 从数据库中将图片路径取出
func GetImages(home string, homeID uint) (*[]models.Attachment, error) {
	var images []models.Attachment
	// 查询出数据库中相应的所有图片路径
	err := globals.DB.Select("path").Where("home = ? and home_id = ?", home, homeID).Find(&images).Error
	if err != nil {
		return nil, fmt.Errorf("GetImages -> %s", err)
	}
	return &images, nil
}

// GetAdvertisementImage 专门用于取数据库中的广告图片
func GetAdvertisementImage(home string, status int) (*[]models.Advertisement, error) {
	var images []models.Advertisement
	// 查询出数据库中相应的所有图片路径
	err := globals.DB.Select("path").Where("home = ? and status = ?", home, status).Find(&images).Error
	if err != nil {
		return nil, fmt.Errorf("GetAdvertisementImage -> %s", err)
	}
	return &images, nil
}
