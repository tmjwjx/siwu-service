package controllers

import (
	"fmt"
	"forum/internal/message/logics"
	"forum/internal/message/requests"
	"forum/pkg/globals"
	"forum/pkg/response"
	"github.com/gin-gonic/gin"
	"net/http"
)

// LikeUnreadCtrl
// @Description: 未读点赞消息数量
// @param        c *gin.Context
// @Author tianjiajie 2025-02-12 21:26:37
func LikeUnreadCtrl(c *gin.Context) {
	// 初始化需要的变量
	db := globals.DB

	// 获取用户ID
	userId, _ := c.Get("id")

	// 进入业务层
	count, err := logics.LikeUnreadCountLogic(db, userId.(uint))
	if err != nil {
		globals.Log.Errorf("Failed to get unread likes: %v", err)
		data := response.NewAppErr(globals.StatusInternalServerError, err, nil)
		response.Failed(c, http.StatusInternalServerError, data)
		return
	}

	// 返回响应
	data := response.NewAppData(globals.StatusOK, "获取未读点赞消息数量成功", count)
	response.Success(c, http.StatusOK, data)
}

// LikeMessageCtrl
// @Description: 点赞消息
// @param        c *gin.Context
// @Author tianjiajie 2024-10-04 20:59:03
func LikeMessageCtrl(c *gin.Context) {

	// 初始化需要的变量
	db := globals.DB
	var req requests.MessageReq

	// 获取用户ID
	userId, _ := c.Get("id")
	globals.Log.Infoln("userId:", userId)

	// 绑定查询参数到 req 变量，如果绑定失败，返回错误信息
	if err := c.ShouldBind(&req); err != nil {
		// 日志记录错误信息
		globals.Log.Errorf("绑定req失败 err = %s", err)
		// 返回错误响应
		data := response.NewAppErr(globals.StatusBadRequest, err, nil)
		response.Failed(c, http.StatusBadRequest, data)
		return // 结束函数执行
	}
	fmt.Printf("%v", req)

	//进入业务层
	likeList, err := logics.LikeMessageLogic(db, req, userId.(uint))
	if err != nil {
		globals.Log.Errorf("点赞消息加载失败 err = %s", err)
		data := response.NewAppErr(globals.StatusInternalServerError, err, nil)
		response.Failed(c, http.StatusInternalServerError, data)
		return
	}

	// 返回响应
	data := response.NewAppData(globals.StatusOK, "点赞消息加载成功", likeList)
	response.Success(c, http.StatusOK, data)
}
