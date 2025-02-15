package controllers

import (
	"forum/internal/image/logics"
	"forum/internal/image/requests"
	"forum/pkg/globals"
	"forum/pkg/response"
	"github.com/gin-gonic/gin"
	"net/http"
)

// CompressImageCtrl 压缩图片
func CompressImageCtrl(c *gin.Context) {

	var req requests.CompressImageReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
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
