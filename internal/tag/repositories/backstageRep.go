package repositories

import (
	"fmt"
	"forum/internal/internalPkg/internalUtils"
	"forum/internal/models"
	"forum/internal/tag/requests"
	"forum/pkg/globals"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// AddTagRep 新增标签
func AddTagRep(c *gin.Context, db *gorm.DB, req *requests.BsAddTagReq) (error, int) {

	// 开启事务
	tx := db.Begin()
	if tx.Error != nil {
		return fmt.Errorf("AddTagRep -> 开启事务失败 -> %s", tx.Error), 500
	}

	// 查询该标签是否已经存在
	var tag models.Tag
	//err := tx.Where("name = ?", req.Name).First(&tag).Error
	//if err == nil {
	//	return fmt.Errorf("AddTagRep -> 该标签已经存在"), 500
	//}

	// 到这里说明该标签不存在，可以插入到数据库中
	// 将处图片之外的信息存到数据库中
	tag.Name = req.Name
	tag.Description = req.Description
	tag.ArticleCount = req.ArticleCount
	tag.Heat = req.Heat
	tag.FansCount = req.FansCount

	err := tx.Create(&tag).Error
	if err != nil {
		// 回滚事务
		tx.Rollback()
		return fmt.Errorf("AddTagRep -> 将除图片之外的信息存到数据库中失败 -> %s", err), 500
	}

	// 获取标签的的ID
	err = tx.Where("name = ?", req.Name).First(&tag).Error
	if err != nil {
		// 回滚事务
		tx.Rollback()
		return fmt.Errorf("AddTagRep ->  获取标签的的ID失败 -> %s", err), 500
	}

	// 存储图片的相关信息
	//err, status := controllers.UploadImagesControllers(c, "标签", tag.ID)
	//if err != nil {
	//	return fmt.Errorf("AddTagRep ->  存储图片的相关信息失败 -> %s", err), status
	//}

	u := &internalUtils.UrlParam{
		UrlPath: req.Path,
		Home:    globals.TagHome,
		HomeID:  tag.ID,
		DB:      db,
	}
	err = internalUtils.StoreUrl(u)
	if err != nil {
		return fmt.Errorf("AddTagRep -> 存储图片的相关信息失败 -> %s", err), 500
	}

	// 提交事务
	err = tx.Commit().Error
	if err != nil {
		return fmt.Errorf("AddTagRep -> 提交事务失败 -> %s", err), 500
	}

	return nil, 200

}

// DeleteTagRep 删除标签
func DeleteTagRep(db *gorm.DB, req *requests.BsDelTagReq) error {

	// 开启事务
	tx := db.Begin()
	if tx.Error != nil {
		return fmt.Errorf("DeleteTagRep -> 开启事务失败 -> %s", tx.Error)
	}

	// 删除标签
	result := tx.Model(&models.Tag{}).Where("id = ?", req.ID).Delete(nil)
	if result.Error != nil {
		tx.Rollback() // 回滚事务
		return fmt.Errorf("DeleteTagRep -> %s", result.Error)
	} else if result.RowsAffected == 0 {
		tx.Rollback() // 回滚事务
		return fmt.Errorf("DeleteTagRep -> 没有找到匹配的记录或记录已经被删除")
	}

	// 删除文件系统中的图片
	err := internalUtils.DeleteFile("标签", req.ID)
	if err != nil {
		return fmt.Errorf("DeleteTagRep -> %s", err)
	}

	// 提交事务
	err = tx.Commit().Error
	if err != nil {
		return fmt.Errorf("DeleteTagRep -> 提交事务失败 -> %s", err)
	}

	return nil

}

// BatchDelTagRep 批量删除
func BatchDelTagRep(db *gorm.DB, req *requests.BsBatchDelTagReq) error {

	// 开启事务
	tx := db.Begin()
	if tx.Error != nil {
		return fmt.Errorf("BatchDelTagRep -> 开启事务失败 -> %s", tx.Error)
	}

	// 删除标签
	for _, id := range req.ID {

		// 之后如果没删除成功，就将没删除成功的标签名返回给后台

		result := tx.Model(&models.Tag{}).Where("id = ?", id).Delete(nil)
		if result.Error != nil {
			// 回滚事务
			tx.Rollback()
			return fmt.Errorf("BatchDelTagRep1 -> 标签批量删除失败 -> %s", result.Error)
		} else if result.RowsAffected == 0 {
			// 回滚事务
			tx.Rollback()
			return fmt.Errorf("BatchDelTagRep2 -> 没有找到 id : %d 的记录或记录已经被删除", id)
		}

		// 删除文件系统中的图片
		err := internalUtils.DeleteFile("标签", id)
		if err != nil {
			return fmt.Errorf("BatchDelTagRep3 -> %s", err)
		}

	}

	// 提交事务
	err := tx.Commit().Error
	if err != nil {
		return fmt.Errorf("BatchDelTagRep4 -> 提交事务失败 -> %s", err)
	}

	return nil

}

// UpdateTagRep 更新标签
func UpdateTagRep(c *gin.Context, db *gorm.DB, req *requests.BsUpTagReq) (error, int) {

	// 开启事务
	tx := db.Begin()
	if tx.Error != nil {
		return fmt.Errorf("UpdateTagRep -> 开启事务失败 -> %s", tx.Error), 500
	}

	var tag models.Tag
	// 更新标签头像以外得信息
	err := tx.Model(&models.Tag{}).Where("id = ? and deleted_at is null", req.ID).First(&tag).Error

	if err != nil {
		// 回滚事务
		tx.Rollback()
		return fmt.Errorf("UpdateTagRep -> 该标签不存在 -> %s", err), 500
	}

	err = tx.Model(&tag).Updates(map[string]interface{}{
		"name":          req.Name,
		"description":   req.Description,
		"article_count": req.ArticleCount,
		"heat":          req.Heat,
		"fans_count":    req.FansCount,
	}).Error
	if err != nil {
		// 回滚事务
		tx.Rollback()
		return fmt.Errorf("UpdateTagRep -> 标签信息更新失败 -> %s", err), 500
	}

	//// 更新标签头像
	//err, status := controllers.UploadImagesControllers(c, "标签", req.ID)
	//if err != nil {
	//	return fmt.Errorf("UpdateTagRep ->  更新标签头像失败 -> %s", err), status
	//}
	u := &internalUtils.UrlParam{
		UrlPath: req.Path,
		Home:    globals.TagHome,
		HomeID:  tag.ID,
		DB:      db,
	}
	err = internalUtils.StoreUrl(u)
	if err != nil {
		return fmt.Errorf("AddTagRep -> 存储图片的相关信息失败 -> %s", err), 500
	}

	// 提交事务
	err = tx.Commit().Error
	if err != nil {
		return fmt.Errorf("UpdateTagRep -> 提交事务失败 -> %s", err), 500
	}

	return nil, 200
}

// QueryTagRep 查询标签
func QueryTagRep(db *gorm.DB, req *requests.BsQueTagReq) (*requests.BsQueTagRes, error) {

	// 查询标签是否存在
	var tag models.Tag
	err := db.Where("name = ?", req.Name).First(&tag).Error
	if err != nil {
		return nil, fmt.Errorf("QueryTagRep -> 该标签不存在 -> %s", err)
	}

	// 设置响应数据
	tagRes := &requests.BsQueTagRes{
		ID:           tag.ID,
		Name:         tag.Name,
		Description:  tag.Description,
		ArticleCount: tag.ArticleCount,
		Heat:         tag.Heat,
		FansCount:    tag.FansCount,
	}

	// 查询标签头像
	images, err := internalUtils.GetImages(db, globals.TagHome, tag.ID)
	if err != nil {
		return nil, fmt.Errorf("QueryTagRep -> %s", err)
	} else {
		for _, path := range *images {
			tagRes.Path = path
		}
	}

	return tagRes, nil

}

// BatchQueryTagRep 批量查询标签
func BatchQueryTagRep(db *gorm.DB, req *requests.BsBatchQueTagReq) (*[]*requests.BsQueTagRes, error) {

	var tagRes []*requests.BsQueTagRes
	var tags []models.Tag

	// 查询标签数据
	err := db.Limit(req.Limit).Offset(req.Offset).Find(&tags).Error
	if err != nil {
		return nil, fmt.Errorf("BatchQueryTagRep -> 批量查询标签失败 -> %s", err)
	}

	for _, tag := range tags {
		t := &requests.BsQueTagRes{
			ID:           tag.ID,
			Name:         tag.Name,
			Description:  tag.Description,
			ArticleCount: tag.ArticleCount,
			Heat:         tag.Heat,
			FansCount:    tag.FansCount,
		}
		tagRes = append(tagRes, t)
	}
	// 查询标签头像
	for _, tag := range tagRes {

		// 查询标签头像
		images, err := internalUtils.GetImages(db, globals.TagHome, tag.ID)
		if err != nil {
			return nil, fmt.Errorf("BatchQueryTagRep -> %s", err)
		} else {
			for _, path := range *images {
				tag.Path = path
			}
		}

	}

	return &tagRes, nil

}
