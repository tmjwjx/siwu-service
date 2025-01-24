package internalUtils

import (
	"fmt"
	"forum/internal/models"
	"forum/pkg/globals"
	"gorm.io/gorm"
	"os"
	"path/filepath"
	"strings"
)

/*
以下是取图片路径和存图片路径的示例(以用户头像图片为例)

//查询用户的头像路径
userImages, err := internalUtils.GetImages(u.DB, globals.UserHome, id)
if err != nil {
	return nil, fmt.Errorf("UserReqContext.GetInfo() %v", err)
}
// 没有图片
if userImages == nil {
	return nil, fmt.Errorf("UserReqContext.GetInfo() err = 无法找到id为%d的用户图片", id)
}
avatarPath := (*userImages)[0]

// 存储用户头像路径
if err := internalUtils.StoreUrl(&internalUtils.UrlParam{
	UrlPath: []string{req.AvatarPath},
	Home:    globals.UserHome,
	HomeID:  req.UserId,
	DB:      u.DB,
}); err != nil {
	return fmt.Errorf("UserReqContext.Edit() -> %v", err)
}

*/

// GetImages
// @Description: 从数据库中将图片路径取出
// @Author wangyulong 2024-10-17 16:20:39
// @param        db *gorm.DB
// @param        home globals.Home
// @param        homeID uint
// @return       *[]string
// @return       error
func GetImages(db *gorm.DB, home globals.Home, homeID uint) (*[]string, error) {
	var images []string
	// 查询出数据库中相应的所有图片路径
	err := db.Model(models.Attachment{}).Where("home = ? and home_id = ?", home, homeID).Pluck("path", &images).Error
	if err != nil {
		return nil, fmt.Errorf("GetImages -> 查询出数据库中相应的所有图片路径失败 -> %s", err)
	} else if len(images) == 0 {
		// 数据库中没有图片，直接使用默认的图片
		AssignDefaultValue(home, &images)
	}
	return &images, nil
}

// AdverUrl
// @Description:  存储广告的URL时需要传的参数
// @Author wangyulong 2024-10-17 16:07:20
type AdverUrl struct {
	Home   globals.Home // 广告图片属于哪个页面
	Status int          // 该广告是否被展示
	Paths  []string     // 广告图片路径
	DB     *gorm.DB
}

// StoreAdvertisementUrl
// @Description: 存储广告图片的url
// @Author wangyulong 2024-10-17 16:16:27
// @param        au *AdverUrl
// @return       error
func StoreAdvertisementUrl(au *AdverUrl) error {

	// 开启事务
	tx := au.DB.Begin()
	if tx.Error != nil {
		return fmt.Errorf("StoreAdvertisementUrl -> 开启事务失败 -> %s", tx.Error)
	}

	for _, path := range au.Paths {
		advertisement := &models.Advertisement{
			Home:   au.Home,
			Status: au.Status,
			Path:   path,
		}

		err := tx.Model(models.Advertisement{}).Create(advertisement).Error
		if err != nil {
			tx.Rollback() // 回滚事务
			return fmt.Errorf("StoreAdvertisementUrl - > 存储广告图片失败 -> %s", err)
		}
	}

	//提交事务
	err := tx.Commit().Error
	if err != nil {
		return fmt.Errorf("StoreAdvertisementUrl -> 提交事务失败 -> %s", err)
	}

	return nil
}

// GetAdvertisementImage
// @Description: 专门用于取数据库中的广告图片
// @Author wangyulong 2024-10-17 16:20:22
// @param        db *gorm.DB
// @param        home globals.Home
// @param        status int
// @return       *[]string
// @return       error
func GetAdvertisementImage(db *gorm.DB, home globals.Home, status int) (*[]string, error) {
	var images []string
	// 查询出数据库中相应的所有图片路径
	err := db.Model(models.Advertisement{}).Where("home = ? and status = ?", home, status).Pluck("path", &images).Error
	if err != nil {
		return nil, fmt.Errorf("GetAdvertisementImage -> 查询出数据库中相应的所有图片路径失败 -> %s", err)
	} else if len(images) == 0 {
		// 数据库中没有图片，直接使用默认的图片
		images = append(images, AdvertisementDefaultImage)
	}
	return &images, nil
}

// UrlParam
// @Description: 存储URL时需要传的参数
// @Author wangyulong 2024-10-17 16:20:06
type UrlParam struct {
	UrlPath []string     // url路径
	Home    globals.Home // home:图片所属单位，即属于文章图片还是用户图片或是资源表图片
	HomeID  uint         // homeID:图片对应的具体文章或用户的ID
	DB      *gorm.DB
}

// StoreUrl
// @Description: 存储URL
// @Author wangyulong 2024-10-17 08:35:07
// @param        u *UrlParam
// @return       error
func StoreUrl(u *UrlParam) error {

	var outFile string       // 图片在服务器的存储路径
	var fileExtension string // 图片扩展名

	for _, url := range u.UrlPath {

		if url == "" {
			return nil
		}

		originalURL := url

		// 获取最后一个 '/' 的位置
		lastSlashIndex := strings.LastIndex(originalURL, "/")
		// 如果找到了 '/', 进行截取
		if lastSlashIndex != -1 {

			// 获取最后一个 '/' 前面的部分（不包含最后一个 '/'）
			prefix := originalURL[:lastSlashIndex]
			// 替换前缀为 './static/images'
			outFile = strings.Replace(originalURL, prefix, "./static/images", 1)

		} else {
			return fmt.Errorf("StoreUrlCtrl -> 获取最后一个 '/' 的位置失败")
		}

		// 获取文件名
		fileName := originalURL[strings.LastIndex(originalURL, "/")+1:]
		// 获取扩展名（包括点）
		ext := filepath.Ext(fileName)
		// 去掉前面的点
		if len(ext) > 1 {
			fileExtension = ext[1:]
		} else {
			return fmt.Errorf("UploadHandlerLogic -> 未找到扩展名")
		}

		// 使用 os.Stat 获取文件信息
		fileInfo, err := os.Stat(outFile)
		if err != nil {
			return fmt.Errorf("UploadHandlerLogic -> 使用 os.Stat 获取文件信息失败 -> %s", err)
		}

		// // 获取文件大小（字节）
		fileSize := fileInfo.Size()

		// 将文件路径及其相关信息存入数据库中
		attachment := &models.Attachment{
			Home:   u.Home,
			HomeID: u.HomeID,
			Name:   fileName,
			Type:   "images/" + fileExtension, // filepath.Ext(filename) 获得文件的扩展名
			Size:   fileSize,
			Path:   url,
		}

		// 将文件插入数据库中
		err = InsertFile(u.DB, attachment)
		if err != nil {
			return err
		}
	}

	return nil
}

// InsertFile
// @Description: 将图片文件的路径相关的信息存入数据库中
// @Author wangyulong 2024-10-17 16:19:51
// @param        db *gorm.DB
// @param        attachment *models.Attachment
// @return       error
func InsertFile(db *gorm.DB, attachment *models.Attachment) error {

	// 开启事务
	tx := db.Begin()
	if tx.Error != nil {
		return fmt.Errorf("InsertFile -> 开启事务失败 -> %s", tx.Error)
	}

	// 删除旧的图片记录
	//d := tx.Model(&models.Attachment{}).Where("home = ? and home_id = ?", attachment.Home, attachment.HomeID).Delete(nil)
	//if d.Error != nil {
	//	tx.Rollback() // 回滚事务
	//	return fmt.Errorf("InsertFile -> 删除旧的图片记录异常 -> %s", d.Error)
	//} else if d.RowsAffected == 0 {
	//	tx.Rollback() // 回滚事务
	//	return fmt.Errorf("InsertFile -> 没有找到匹配的记录或记录已经被删除")
	//}

	// 查询记录是否存在
	var existingAttachment models.Attachment
	if tx.Model(&models.Attachment{}).Where("home = ? AND home_id = ?", attachment.Home, attachment.HomeID).First(&existingAttachment).Error == nil {
		// 如果记录存在，设置 ID
		attachment.ID = existingAttachment.ID
	}

	// 向数据库中存入文件数据
	result := tx.Model(&models.Attachment{}).Where("id = ?", attachment.ID).Omit("created_at").Save(attachment).Debug()
	if result.Error != nil {
		tx.Rollback()
		return fmt.Errorf("InsertFile -> 向数据库中存入文件数据 -> %s", result.Error)
	}

	//result := tx.Model(&models.Attachment{}).Where("home = ? and home_id = ?", attachment.Home, attachment.HomeID).Save(attachment)
	//if result.Error != nil {
	//	tx.Rollback() // 回滚事务
	//	return fmt.Errorf("InsertFile -> 向数据库中存入文件数据 -> %s", result.Error)
	//}

	//提交事务
	err := tx.Commit().Error
	if err != nil {
		return fmt.Errorf("InsertFile -> 提交事务失败 -> %s", err)
	}

	return nil
}

// AssignDefaultValue
// @Description: 分配图片的默认值
// @Author wangyulong 2024-10-17 15:17:44
// @param        home globals.Home
// @param        images *[]string
func AssignDefaultValue(home globals.Home, images *[]string) {
	if home == globals.UserHome {
		*images = append(*images, UserDefaultImage)
	} else if home == globals.ArticleHome {
		*images = append(*images, UserDefaultImage)
	} else if home == globals.TagHome {
		*images = append(*images, TagDefaultImage)
	} else if home == globals.CommentHome {
		*images = append(*images, CommentDefaultImage)
	} else if home == globals.CategoryHome {
		*images = append(*images, CategoryDefaultImage)
	}
}
