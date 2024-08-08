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
	var tag requests.TagReq
	if err := c.ShouldBind(&tag); err != nil {
		e := response.NewAppErr(globals.StatusBadRequest, err, nil)
		response.Failed(c, 400, e)
		return
	}

	// 业务处理
	fansCount, err := logics.UpdateTagUserCountLogic(tag.ID)

	// 返回响应
	if err != nil {
		// 人数更新失败，返回错误响应
		e := response.NewAppErr(globals.StatusInternalServerError, err, nil)
		response.Failed(c, 500, e)
		return
	}
	// 人数更新成功，返回现在人数
	d := response.NewAppData(globals.StatusOK, "人数更新成功", fansCount)
	response.Success(c, 200, d)
}

// UpdateTag 更新前端的标签页
func UpdateTag(c *gin.Context) {

	// 业务处理
	tagRes, err := logics.UpdateTagArticleCountLogic() // 更新前端的标签页

	// 返回响应
	if err != nil {
		// 更新失败，返回错误响应
		e := response.NewAppErr(globals.StatusInternalServerError, err, nil)
		response.Failed(c, 500, e)
		return
	}
	// 更新成功，返回成功响应
	d := response.NewAppData(globals.StatusOK, "文章数量更新成功", tagRes)
	response.Success(c, 200, d)
}
