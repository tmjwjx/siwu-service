package repositories

import (
	"fmt"
	"forum/internal/image/requests"
	"forum/internal/models"
	"forum/pkg/globals"
)

// InsertFile 将图片文件的路径相关的信息存入数据库中
func InsertFile(attachment *requests.Attachment) error {

	// 开启事务
	tx := globals.DB.Begin()
	if tx.Error != nil {
		return fmt.Errorf("InsertFile -> 开启事务失败 -> %s", tx.Error)
	}

	// 向数据库中存入文件数据
	result := tx.Model(&models.Attachment{}).Where("home = ? and home_id = ?", attachment.Home, attachment.HomeID).Save(attachment)
	if result.Error != nil {
		tx.Rollback() // 回滚事务
		return fmt.Errorf("InsertFile -> 向数据库中存入文件数据 -> %s", result.Error)
	}

	//提交事务
	err := tx.Commit().Error
	if err != nil {
		return fmt.Errorf("InsertFile -> 提交事务失败 -> %s", err)
	}

	return nil
}

// GetImages 从数据库中将图片路径取出
func GetImages(home string, homeID uint) (*[]models.Attachment, error) {
	var images []models.Attachment
	// 查询出数据库中相应的所有图片路径
	err := globals.DB.Select("Path").Where("home = ? and home_id = ?", home, homeID).Find(&images).Error
	if err != nil {
		return nil, fmt.Errorf("GetImages -> 在数据库中没有查到对应图片 -> %s", err)
	} else if len(images) == 0 {
		return nil, fmt.Errorf("GetImages -> 在数据库中没有查到对应图片")
	}
	return &images, nil
}

// GetAdvertisementImage 专门用于取数据库中的广告图片
func GetAdvertisementImage(home string, status int) (*[]models.Advertisement, error) {
	var images []models.Advertisement
	// 查询出数据库中相应的所有图片路径
	err := globals.DB.Select("path").Where("home = ? and status = ?", home, status).Find(&images).Error
	if err != nil {
		return nil, fmt.Errorf("GetAdvertisementImage -> 在数据库中没有查到对应的广告图片 -> %s", err)
	}
	return &images, nil
}
