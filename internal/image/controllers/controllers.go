package controllers

import (
	"fmt"
	"forum/internal/image/logics"
	"forum/internal/models"
)

func StoreUrlCtrl(u *logics.UrlParam) error {

	// 逻辑处理
	err := logics.StoreUrlLogic(u)
	return err
}

// GetImagesControllers 从数据库中将图片路径取出
func GetImagesControllers(home string, homeID uint) (*[]models.Attachment, error) {
	images, err := logics.GetImagesLogic(home, homeID)
	return images, err
}

// GetAdvertisementImageCtrl 专门用于取数据库中的广告图片
func GetAdvertisementImageCtrl(home string, status int) (*[]models.Advertisement, error) {
	images, err := logics.GetAdvertisementImageLogic(home, status)
	if err != nil {
		return nil, fmt.Errorf("GetAdvertisementImageCtrl -> %s", err)
	}
	return images, nil
}
