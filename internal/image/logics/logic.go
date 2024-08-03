package logics

import (
	"fmt"
	"forum/internal/image/repositorys"
	"forum/internal/image/requests"
	"github.com/gin-gonic/gin"
	"net/http"
)

func UploadHandlerLogic(c *gin.Context) {
	// 使用 MultipartForm 提取所有字段
	form, _ := c.MultipartForm()
	// 提取文件
	files := form.File["upload[]"]
	for _, file := range files {
		// 将文件内容写入目标文件
		err := c.SaveUploadedFile(file, "./static/image"+file.Filename)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"msg": file.Filename + "文件保存失败"})
		}
		// 将文件路径及其相关信息存入数据库中
		attachment := &requests.Attachment{
			Name: file.Filename,
			Type: "image/png",
			Size: file.Size,
			Path: "./static/image" + file.Filename,
		}
		// 将文件插入数据库中
		result := repositorys.InsertFile(attachment)
		if result.Error != nil {
			// 插入数据失败，返回响应
			c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
			return
		}
	}
	// 所有文件上传成功，返回响应
	c.JSON(http.StatusOK, gin.H{"msg": fmt.Sprintf("成功上传 %d 个文件", len(files))})
}
