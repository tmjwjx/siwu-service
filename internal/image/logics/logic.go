package logics

import (
	"fmt"
	"forum/internal/image/repositorys"
	"forum/internal/image/requests"
	"forum/pkg/globals"
	"forum/pkg/response"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"path/filepath"
)

// UploadHandlerLogic 图片文件的逻辑处理
func UploadHandlerLogic(c *gin.Context, kind int, ID uint) {
	// 使用 MultipartForm 提取所有字段
	form, _ := c.MultipartForm()
	// 提取文件
	files := form.File["upload[]"]
	for _, file := range files {
		// 将文件内容写入目标文件
		err := c.SaveUploadedFile(file, "./static/images"+file.Filename)
		if err != nil {
			e := response.NewAppErr(globals.StatusInternalServerError, err, nil)
			response.Failed(c, e, 5000)
			return
		}

		// 生成唯一的文件名
		uniqueFilename := generateUniqueFilename(file.Filename)
		// 将文件路径及其相关信息存入数据库中
		attachment := &requests.Attachment{
			Name: file.Filename,
			Type: "images/" + filepath.Ext(file.Filename), // filepath.Ext(filename) 获得文件的扩展名
			Size: file.Size,
			Path: "/images" + uniqueFilename,
		}

		// 将文件插入数据库中
		repositorys.InsertFile(c, attachment, kind, ID)
	}
	// 所有文件上传成功，返回响应
	d := response.NewAppData(globals.StatusOK, "用户信息更新成功", nil)
	response.Success(c, d, 2000)

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
