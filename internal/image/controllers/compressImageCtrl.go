package controllers

import (
	"fmt"
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
	path := c.Param("path")
	width := c.Query("width")
	height := c.Query("height")
	level := c.Query("level")

	if path == "" {
		e := response.NewAppErr(globals.StatusInternalServerError, fmt.Errorf("上传的path不能为空"), nil)
		response.Failed(c, 500, e)
		return
	}

	if width == "" {
		width = "0"
	}

	if height == "" {
		height = "0"
	}

	if level == "" {
		level = "0"
	}

	widthInt, err := internalUtils.ChangeStringToInt(width)
	if err != nil {
		e := response.NewAppErr(globals.StatusInternalServerError, fmt.Errorf("上传的压缩目标宽度的数值不是整数"), nil)
		response.Failed(c, 500, e)
		return
	}
	if widthInt < 0 || widthInt > 8000 {
		e := response.NewAppErr(globals.StatusInternalServerError, fmt.Errorf("上传的压缩目标宽度的数值不能超出范围(0-8000)包含0和8000"), nil)
		response.Failed(c, 500, e)
		return
	}

	heightInt, err := internalUtils.ChangeStringToInt(height)
	if err != nil {
		e := response.NewAppErr(globals.StatusInternalServerError, fmt.Errorf("上传的压缩目标高度的数值不是整数"), nil)
		response.Failed(c, 500, e)
		return
	}
	if heightInt < 0 || heightInt > 8000 {
		e := response.NewAppErr(globals.StatusInternalServerError, fmt.Errorf("上传的压缩目标宽度的数值不能超出范围(0-8000)包含0和8000"), nil)
		response.Failed(c, 500, e)
		return
	}

	levelInt, err := internalUtils.ChangeStringToInt(level)
	if err != nil {
		e := response.NewAppErr(globals.StatusInternalServerError, fmt.Errorf("上传的压缩级别的数值不是整数"), nil)
		response.Failed(c, 500, e)
		return
	}
	if levelInt < 0 || levelInt > 9 {
		e := response.NewAppErr(globals.StatusInternalServerError, fmt.Errorf("上传的压缩级别的数值不能超出范围(0-9)包含0和9"), nil)
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

	// 设置响应的Accept-Ranges
	c.Header("Accept-Ranges", "bytes")

	// 直接返回压缩后的图片数据
	c.Data(http.StatusOK, contentType, processed)

}
