package controllers

import (
	"forum/internal/tag/logics"
	"forum/internal/tag/requests"
	"forum/pkg/globals"
	"forum/pkg/response"
	"github.com/gin-gonic/gin"
)

// UpdateTagUserCount 更新标签的关注人数
func UpdateTagUserCount(c *gin.Context) {
	var tag requests.Tag
	if err := c.ShouldBind(&tag); err != nil {
		e := response.NewAppErr(globals.StatusBadRequest, err, nil)
		response.Failed(c, 400, e)
		return
	}

	// 业务处理
	fansCount, err := logics.UpdateTagUserCountLogic(tag.ID)
	if err != nil {
		// 人数更新失败，返回错误响应
		e := response.NewAppErr(globals.StatusInternalServerError, err, nil)
		response.Failed(c, 500, e)
		return
	} else {
		// 人数更新成功，返回现在人数
		d := response.NewAppData(globals.StatusOK, "人数更新成功", fansCount)
		response.Success(c, 200, d)
	}
}
