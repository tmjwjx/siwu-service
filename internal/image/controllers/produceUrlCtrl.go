package controllers

import (
	"fmt"
	"forum/internal/image/logics"
	"forum/internal/image/requests"
	"forum/pkg/globals"
	"forum/pkg/response"
	"github.com/gin-gonic/gin"
)

func ProduceUrlCtrl(c *gin.Context) {

	// 解析multipart/form-data
	if err := c.Request.ParseMultipartForm(32 << 20); err != nil { // 设置最大大小为 32MB
		e := response.NewAppErr(globals.StatusBadRequest, fmt.Errorf("ProduceUrlCtrl -> 图片超过规定的最大允许大小 32MB"), nil)
		response.Failed(c, 400, e)
	}
	// 具体逻辑实现
	imageUrl, err := logics.ProduceUrlLogic(c)

	// 返回响应
	if err != nil {
		//var urls []*requests.UrlPath
		//res := &requests.ImageUrl{
		//	Data:  urls,
		//	Errno: 1,
		//}
		// 如果文件中没有图片，直接返回空。
		urls := make([]*requests.UrlPath, 0)
		res := &requests.ImageUrl{
			Data:  urls,
			Errno: 1,
		}
		c.JSON(500, res)
	} else {
		c.JSON(200, imageUrl)
	}

}
