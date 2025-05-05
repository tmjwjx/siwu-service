package controllers

import (
	"fmt"
	"forum/internal/image/logics"
	"forum/internal/image/requests"
	"forum/pkg/globals"
	"forum/pkg/response"
	"github.com/gin-gonic/gin"
)

// ProduceUrlCtrl
// @Description: 根据上传的图片，生成url路径，并将图片存到静态文件中，url路径存到数据库中。
// @Author wangyulong 2025-02-13 15:38:20
// @param        c *gin.Context
func ProduceUrlCtrl(c *gin.Context) {

	// 解析multipart/form-data
	if err := c.Request.ParseMultipartForm(globals.SConfig.Limit << 20); err != nil { // 设置最大大小为 32MB
		e := response.NewAppErr(globals.StatusBadRequest, fmt.Errorf("ProduceUrlCtrl -> 图片超过规定的最大允许大小 32MB"), nil)
		response.Failed(c, 400, e)
	}

	watermarkParam := requests.WatermarkParam{
		Watermark: globals.SConfig.Watermark,
		Size:      globals.SConfig.Size,
		Scale:     globals.SConfig.Scale,
		Margin:    globals.SConfig.Margin,
	}
	// 具体逻辑实现
	imageUrl, err := logics.ProduceUrlLogic(c, watermarkParam)

	// 返回响应
	if err != nil {
		if err.Error() == "上传的图片比例不对" {
			urls := make([]*requests.UrlPath, 0)
			res := &requests.ImageUrl{
				Message: "上传的图片比例不对",
				Data:    urls,
				Errno:   0,
			}
			c.JSON(200, res)
			return
		}
		fmt.Println(err)
		// 如果文件中没有图片，直接返回空。
		urls := make([]*requests.UrlPath, 0)
		res := &requests.ImageUrl{
			Message: "上传的文件中没有图片",
			Data:    urls,
			Errno:   1,
		}
		c.JSON(500, res)
	} else {
		c.JSON(200, imageUrl)
	}

}
