package logics

import (
	"fmt"
	"forum/internal/image/repositorys"
	"forum/internal/image/requests"
	"forum/internal/models"
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
	files := form.File["upload[]"]
	for _, file := range files {
		// 将文件内容写入目标文件
		err := c.SaveUploadedFile(file, "./static/images/"+file.Filename)
		if err != nil {
			//e := response.NewAppErr(globals.StatusInternalServerError, err, nil)
			//response.Failed(c, e, 5000)
			return err
		}

		// 生成唯一的文件名
		uniqueFilename := generateUniqueFilename(file.Filename)
		// 将文件路径及其相关信息存入数据库中
		attachment := &requests.Attachment{
			Home:   home,
			HomeID: homeID,
			Name:   file.Filename,
			Type:   "images/" + filepath.Ext(file.Filename), // filepath.Ext(filename) 获得文件的扩展名
			Size:   file.Size,
			Path:   "/images/" + uniqueFilename,
		}

		// 将文件插入数据库中
		err = repositorys.InsertFile(attachment)
		if err != nil {
			return err
		}
	}

	return nil
}

// GetImagesLogic 从数据库中将图片路径取出
func GetImagesLogic(home string, homeID uint) (*[]models.Attachment, error) {
	images, err := repositorys.GetImages(home, homeID)
	if err != nil {
		return nil, fmt.Errorf("GetImagesLogic -> %s", err)
	}
	return images, nil
}

// GetAdvertisementImageLogic 专门用于取数据库中的广告图片
func GetAdvertisementImageLogic(home string, status int) (*[]models.Advertisement, error) {
	// 查询出数据库中相应的所有图片路径
	images, err := repositorys.GetAdvertisementImage(home, status)
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
