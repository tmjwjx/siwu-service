package controllers

import (
	"forum/internal/message/logics"
	"forum/internal/message/requests"
	"forum/pkg/globals"
	"forum/pkg/response"
	"github.com/gin-gonic/gin"
	"net/http"
)

// CommentUnreadCtrl
// @Description: 未读评论消息数量
// @param        c *gin.Context
// @Author tianjiajie 2025-02-12 22:11:18
func CommentUnreadCtrl(c *gin.Context) {
	db := globals.DB
	userId, _ := c.Get("id")
	count, err := logics.CommentUnreadCountLogic(db, userId.(uint))
	if err != nil {
		globals.Log.Errorf("Failed to get unread comments: %v", err)
		data := response.NewAppErr(globals.StatusInternalServerError, err, nil)
		response.Failed(c, http.StatusInternalServerError, data)
		return
	}
	data := response.NewAppData(globals.StatusOK, "获取未读评论消息数量成功", count)
	response.Success(c, http.StatusOK, data)
}

// CommentMesCtrl
// @Description: 评论消息
// @param        c *gin.Context
// @Author tianjiajie 2025-01-18 11:10:39
func CommentMesCtrl(c *gin.Context) {
	// 初始化需要的变量
	db := globals.DB
	var req *requests.MessageReq
	userId, ok := c.Get("id")
	if ok != true {
		data := response.NewAppErr(globals.StatusBadRequest, nil, nil)
		response.Failed(c, http.StatusBadRequest, data)
		return
	}

	// 绑定查询参数到 req 变量，如果绑定失败，返回错误信息
	if err := c.ShouldBind(&req); err != nil {
		// 日志记录错误信息
		globals.Log.Errorf("绑定req失败 err = %s", err)
		// 返回错误响应
		data := response.NewAppErr(globals.StatusBadRequest, err, nil)
		response.Failed(c, http.StatusBadRequest, data)
		return // 结束函数执行
	}

	// 进入业务层
	comment, err := logics.CommentMesLogic(db, req, userId.(uint))
	if err != nil {
		globals.Log.Errorf("获取数据失败 err = %s", err)
		data := response.NewAppErr(globals.StatusInternalServerError, err, nil)
		response.Failed(c, http.StatusInternalServerError, data)
		return
	}

	// 返回响应
	data := response.NewAppData(globals.StatusOK, "成功", comment)
	response.Success(c, http.StatusOK, data)
}
