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
		//e := response.NewAppErr(globals.StatusInternalServerError, err, nil)
		//response.Failed(c, 500, e)
		//return
		url := &requests.UrlPath{
			Url: "",
		}
		res := &requests.ImageUrl{
			Errno: 0,
			Data:  url,
		}
		c.JSON(500, res)
	}

	//d := response.NewAppData(globals.StatusOK, "图片保存成功且图片url生成成功", imageUrl)
	//response.Success(c, 200, d)
	c.JSON(200, imageUrl)
}
