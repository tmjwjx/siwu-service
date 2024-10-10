package logics

import (
	"fmt"
	"forum/internal/image/repositories"
	"forum/internal/image/requests"
	"forum/internal/internalPkg/internalUtils"
	"forum/internal/models"
	"forum/pkg/globals"
	"os"
	"path/filepath"
	"strings"
)

// UrlParam 存储URL时许需要传的参数
type UrlParam struct {
	UrlPath []string     // url路径
	Home    globals.Home // home:图片所属单位，即属于文章图片还是用户图片或是资源表图片
	HomeID  uint         // homeID:图片对应的具体文章或用户的ID
}

func StoreUrlLogic(u *UrlParam) error {

	var outFile string       // 图片在服务器的存储路径
	var fileExtension string // 图片扩展名

	for _, url := range u.UrlPath {

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
		attachment := &requests.Attachment{
			Home:   u.Home,
			HomeID: u.HomeID,
			Name:   fileName,
			Type:   "images/" + fileExtension, // filepath.Ext(filename) 获得文件的扩展名
			Size:   fileSize,
			Path:   url,
		}

		// 将文件插入数据库中
		err = repositories.InsertFile(attachment)
		if err != nil {
			return err
		}
	}

	return nil
}

// GetImagesLogic 从数据库中将图片路径取出
func GetImagesLogic(home string, homeID uint) (*[]models.Attachment, error) {
	images, err := repositories.GetImages(home, homeID)
	return images, err
}

// GetAdvertisementImageLogic 专门用于取数据库中的广告图片
func GetAdvertisementImageLogic(home string, status int) (*[]models.Advertisement, error) {
	// 查询出数据库中相应的所有图片路径
	images, err := repositories.GetAdvertisementImage(home, status)
	if err != nil {
		return nil, fmt.Errorf("GetAdvertisementImageLogic -> %s", err)
	}
	return images, nil
}
