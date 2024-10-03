package coctrollers

import (
	"fmt"
	"forum/internal/message/comment/logics"
	"forum/internal/message/comment/requests"
	"forum/pkg/globals"
	"forum/pkg/response"
	"github.com/gin-gonic/gin"
)

// CommentMesCtrl 评论消息(前台)
func CommentMesCtrl(c *gin.Context) {

	// 获取参数
	// 从上下文中获取用户的id
	id, exists := c.Get("id")
	if !exists {
		e := response.NewAppErr(globals.StatusBadRequest, fmt.Errorf("CommentMesCtrl -> 从上下文中获取用户的id失败"), nil)
		response.Failed(c, 400, e)
		return
	}
	// 解析json格式的数据
	var req requests.CommentMesReq
	err := c.ShouldBindJSON(&req)
	if err != nil {
		e := response.NewAppErr(globals.StatusBadRequest, fmt.Errorf("CommentMesCtrl -> 解析json格式的数据失败 -> %s", err), nil)
		response.Failed(c, 400, e)
		return
	}

	// 逻辑处理
	res, err := logics.CommentMesLogic(globals.DB, &req, id)

	// 返回响应
	if err != nil {
		e := response.NewAppErr(globals.StatusInternalServerError, err, nil)
		response.Failed(c, 500, e)
		return
	}
	d := response.NewAppData(globals.StatusOK, "修改菜单成功", res)
	response.Success(c, 200, d)

}
