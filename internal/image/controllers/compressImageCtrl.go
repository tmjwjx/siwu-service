package controllers

import (
	"forum/internal/image/logics"
	"forum/internal/image/requests"
	"forum/internal/internalPkg/internalUtils"
	"forum/pkg/globals"
	"forum/pkg/response"
	"github.com/gin-gonic/gin"
	"net/http"
)

// CompressImageCtrl 压缩图片
func CompressImageCtrl(c *gin.Context) {

	// 获取参数
	path := c.Query("path")
	width := c.Query("width")
	height := c.Query("height")
	level := c.Query("level")

	widthInt, err := internalUtils.ChangeStringToInt(width)
	if err != nil {
		e := response.NewAppErr(globals.StatusInternalServerError, err, nil)
		response.Failed(c, 500, e)
		return
	}

	heightInt, err := internalUtils.ChangeStringToInt(height)
	if err != nil {
		e := response.NewAppErr(globals.StatusInternalServerError, err, nil)
		response.Failed(c, 500, e)
		return
	}

	levelInt, err := internalUtils.ChangeStringToInt(level)
	if err != nil {
		e := response.NewAppErr(globals.StatusInternalServerError, err, nil)
		response.Failed(c, 500, e)
		return
	}

	req := requests.CompressImageReq{
		Path:   path,
		Width:  widthInt,
		Height: heightInt,
		Level:  levelInt,
	}

	// 逻辑处理
	contentType, processed, err := logics.CompressImageLogic(globals.RDB, req)
	// 返回响应
	if err != nil {
		e := response.NewAppErr(globals.StatusInternalServerError, err, nil)
		response.Failed(c, 500, e)
		return
	}

	// 设置响应的Content-Type
	c.Header("Content-Type", contentType)

	// 直接返回压缩后的图片数据
	c.Data(http.StatusOK, contentType, processed)

}
