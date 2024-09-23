package logics

import (
	"fmt"
	"forum/internal/image/repositories"
	"forum/internal/image/requests"
	"forum/internal/internal_pkg/internal_utils"
	"forum/internal/models"
	"forum/pkg/globals"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"path/filepath"
)

// UploadHandlerLogic 图片文件的逻辑处理
// home:图片所属单位，即属于文章图片还是用户图片或是资源表图片
// homeID:图片对应的具体文章或用户的ID
func UploadHandlerLogic(c *gin.Context, home string, homeID uint) error {
	// 使用 MultipartForm 提取所有字段
	form, _ := c.MultipartForm()
	// 提取文件
	files := form.File["uploads"]
	for _, file := range files {

		// 如果文件中没有图片，直接返回nil。
		if file.Size == 0 {
			return nil
		}

		// 生成唯一的文件名
		uniqueFilename := generateUniqueFilename(file.Filename)

		// 删除文件系统中的图片
		err := internal_utils.DeleteFile(home, homeID)
		if err != nil {
			return err
		}

		// 删除 attachments 表中的图片路径

		// 将文件内容写入目标文件
		err = c.SaveUploadedFile(file, globals.SConfig.Path+"/"+uniqueFilename)
		if err != nil {
			return err
		}

		// 将文件路径及其相关信息存入数据库中
		attachment := &requests.Attachment{
			Home:   home,
			HomeID: homeID,
			Name:   file.Filename,
			Type:   "images/" + filepath.Ext(file.Filename), // filepath.Ext(filename) 获得文件的扩展名
			Size:   file.Size,
			Path:   globals.SConfig.Prefix + "/" + uniqueFilename,
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
	/*if err != nil {
		return nil, fmt.Errorf("GetImagesLogic -> %s", err)
	}
	return images, nil*/
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

// generateUniqueFilename 生成唯一文件名
func generateUniqueFilename(filename string) string {
	// 生成一个唯一的 UUID
	uniqueID := uuid.New().String()

	// 分离文件名和扩展名
	base := filename[:len(filename)-len(filepath.Ext(filename))]
	ext := filepath.Ext(filename)

	// 创建一个新的唯一文件名
	return fmt.Sprintf("%s_%s%s", base, uniqueID, ext)
}

/*// DeleteFile 从文件系统中删除图片
func DeleteFile(home string, homeID uint) error {
	var path string
	// 查询要删除的图片文件路径
	err := globals.DB.Model(models.Attachment{}).Where("home = ? and home_id = ?", home, homeID).Select("path").First(&path).Error
	if err != nil {
		//return fmt.Errorf("deleteFile -> %s", err)
		// 没有查到说明文件系统中没有该图片，直接添加进入文件系统即可
		return nil
	}
	path = "./static" + path
	// 删除图片
	err = os.Remove(path)
	if err != nil {
		return fmt.Errorf("deleteFile -> 文件系统中的图片删除失败 -> %s", err)
	}
	return nil
}*/
