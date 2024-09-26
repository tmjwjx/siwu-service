package logics

import (
	"fmt"
	"forum/internal/image/requests"
	"forum/pkg/globals"
	"github.com/gin-gonic/gin"
	"strconv"
)

// ProduceUrlLogic 生成图片的url路径
func ProduceUrlLogic(c *gin.Context) (*requests.ImageUrl, error) {

	// 提取文件
	file, err := c.FormFile("file")

	// 如果文件中没有图片，直接返回空字符串。
	if err != nil || file == nil {
		url := &requests.UrlPath{
			Url: "",
		}
		res := &requests.ImageUrl{
			Errno: 0,
			Data:  url,
		}
		return res, nil
	}

	// 生成唯一的文件名
	uniqueFilename := generateUniqueFilename(file.Filename)

	// 生成url
	imageUrl := "http://" + globals.AppConfig.App.Host + ":" + strconv.Itoa(globals.AppConfig.App.Port) + globals.SConfig.Prefix + "/" + uniqueFilename

	//urlPath := &requests.ImageUrl{
	//	ImageUrl: imageUrl,
	//}
	url := &requests.UrlPath{
		Url: imageUrl,
	}
	res := &requests.ImageUrl{
		Errno: 0,
		Data:  url,
	}

	// 将文件内容写入目标文件
	err = c.SaveUploadedFile(file, globals.SConfig.Path+"/"+uniqueFilename)
	if err != nil {
		return nil, fmt.Errorf("ProduceUrlLogic -> 图片保存失败")
	}

	return res, nil
}
